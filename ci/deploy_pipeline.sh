#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Concourse configuration
TARGET=concourse.gemfire-ci.info
TEAM=developer
PIPELINE=cloudcache-management-cf-plugin
BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Poolsmiths configuration
POOL_NAME=us_2_6
POOLSMITHS_API_TOKEN=0d82e637-6681-4d4a-9e9f-90a71db5de0d

# Version(s) of GemFire for stand-alone testing, whitespace-separated
STANDALONE_GEMFIRE_VERSIONS="9.9 9.10 develop"

# Version(s) of PCC for testing as plugin, whitespace-separated
PCC_VERSIONS="1.10"
PIVNET_API_TOKEN=c90e06904710409eb60d55459e3b3dbd-r

cat << EOF > pipeline.yml
---
resources:
EOF
for gem in $STANDALONE_GEMFIRE_VERSIONS; do
  if [ "$gem" = "9.9" ] ; then
    cat << EOF >> pipeline.yml
- name: gemfire-$gem
  type: s3
  icon: file-cloud
  source:
    bucket: gemfire-releases
    region_name: ((aws-default-region))
    access_key_id: ((gemfire-aws-access-key-id))
    secret_access_key: ((gemfire-aws-secret-access-key))
    versioned_file: 9.9/9.9.0/pivotal-gemfire-9.9.0.tgz
EOF
  else
[ "$gem" = develop ] && br=develop || br=support/$gem
cat << EOF >> pipeline.yml
- name: gemfire-$gem
  type: gcs-resource
  icon: file-cloud
  source:
    bucket: gemfire-build-resources
    json_key: ((concourse-gcp-key))
    regexp: artifacts/gemfire/$br/pivotal-gemfire-regression-(.*).tgz
EOF
  fi
done
for pccver in $PCC_VERSIONS; do
  regex=$(echo "${pccver}." | sed 's#\.#\\\.#')
  cat << EOF >> pipeline.yml
- name: pcc-$pccver
  type: pivnet
  source:
    api_token: $PIVNET_API_TOKEN
    product_slug: p-cloudcache
    product_version: ${regex}.*
EOF
done
cat << EOF >> pipeline.yml
- name: cloudcache-management-cf-plugin
  type: git
  source:
    uri: git@github.com:gemfire/cloudcache-management-cf-plugin.git
    branch: $BRANCH
    ignore_paths:
      - "ci/deploy_pipeline.sh"
      - "**/*.md"
    private_key: ((gemfire-ci-private-key))

- name: golang-image
  type: registry-image
  source:
    repository: golang
    tag: latest

- name: gcp-env
  type: pcf-pool
  tags: [pivotal-internal-worker]
  source:
    api_token: $POOLSMITHS_API_TOKEN
    hostname: environments.toolsmiths.cf-app.com
    pool_name: $POOL_NAME

- name: cloudcache-management-cf-plugin-ci-dockerfile
  type: git
  source:
    uri: git@github.com:gemfire/cloudcache-management-cf-plugin.git
    branch: $BRANCH
    paths:
      - "ci/docker/*"
    private_key: ((!gemfire-ci-private-key))

- name: stemcell
  type: pivnet
  source:
    api_token: $PIVNET_API_TOKEN
    product_slug: stemcells-ubuntu-xenial
    product_version: 456\..*

- name: cloudcache-management-cf-plugin-ci-image
  type: docker-image
  source:
    username: "_json_key"
    password: ((!concourse-gcp-key))
    repository: gcr.io/gemfire-dev/cloudcache-management-cf-plugin-ci

- name: weekly
  type: time
  source:
    start: 3:00 AM
    stop: 11:00 PM
    days: [Wednesday]
    location: America/Los_Angeles


resource_types:
- name: pcf-pool
  type: docker-image
  source:
    repository: cftoolsmiths/toolsmiths-envs-resource

- name: pivnet
  type: docker-image
  source:
    repository: pivotalcf/pivnet-resource
    tag: latest-final

- name: gcs-resource
  type: docker-image
  source:
    repository: frodenas/gcs-resource


jobs:
- name: build-cloudcache-management-cf-plugin
  serial: true
  plan:
  - in_parallel:
    - get: cloudcache-management-cf-plugin
      trigger: true
    - get: golang-image
  - task: build-plugin
    timeout: 1h
    image: golang-image
    config:
      inputs:
        - name: cloudcache-management-cf-plugin
      platform: linux
      run:
        path: /bin/sh
        args:
          - -ec
          - |
            cd cloudcache-management-cf-plugin
            ./build.sh
- name: test-cloudcache-management-cf-plugin
  serial: true
  plan:
  - in_parallel:
    - get: cloudcache-management-cf-plugin
      trigger: true
      passed: [build-cloudcache-management-cf-plugin]
    - get: golang-image
  - task: build-plugin
    timeout: 1h
    image: golang-image
    config:
      inputs:
        - name: cloudcache-management-cf-plugin
      platform: linux
      run:
        path: /bin/sh
        args:
          - -ec
          - |
            apt-get update
            apt-get install -y jq
            cd cloudcache-management-cf-plugin
            go get github.com/onsi/ginkgo/ginkgo
            ginkgo -r

- name: build-docker-image
  plan:
  - get: weekly
    trigger: true
  - get: cloudcache-management-cf-plugin-ci-dockerfile
    trigger: true
  - put: cloudcache-management-cf-plugin-ci-image
    params:
      build: cloudcache-management-cf-plugin-ci-dockerfile/ci/docker
      tag_as_latest: true
EOF
for gem in $STANDALONE_GEMFIRE_VERSIONS; do
  cat << EOF >> pipeline.yml
- name: test-cloudcache-management-cf-standalone-$gem
  serial: true
  plan:
  - in_parallel:
    - get: cloudcache-management-cf-plugin
      trigger: true
      passed: [test-cloudcache-management-cf-plugin]
    - get: golang-image
    - get: gemfire-$gem
  - task: build-plugin
    timeout: 1h
    image: golang-image
    config:
      inputs:
        - name: cloudcache-management-cf-plugin
      outputs:
        - name: pcc-plugin
      platform: linux
      run:
        path: /bin/sh
        args:
          - -ec
          - |
            cd cloudcache-management-cf-plugin
            ./build.sh
            cp pcc ../pcc-plugin/
  - task: standalone-test
    timeout: 1h
    config:
      image_resource:
        type: docker-image
        source:
          repository: openjdk
          tag: 8
      inputs:
        - name: cloudcache-management-cf-plugin
        - name: pcc-plugin
        - name: gemfire-$gem
          path: gemfire
      platform: linux
      run:
        path: /bin/sh
        args:
          - -ecx
          - |
            tar xzf gemfire/pivotal-gemfire-*.tgz
            [ -x bin/gfsh ] && gfsh=bin/gfsh || gfsh=*gemfire*/bin/gfsh
            \$gfsh -e "version --full" -e "start locator"
            cd cloudcache-management-cf-plugin
            pcc="../pcc-plugin/pcc"
            \$pcc --help
            cf="\$pcc http://localhost:7070" ci/smoke-test.bash
EOF
done
for gem in $STANDALONE_GEMFIRE_VERSIONS ; do
  standalone_targets="$standalone_targets,test-cloudcache-management-cf-standalone-$gem"
done
standalone_targets=$(echo "$standalone_targets"|cut -c2-)
for pccver in $PCC_VERSIONS; do
  cat << EOF >> pipeline.yml
- name: test-cloudcache-management-cf-pcc-$pccver
  plan:
  - aggregate:
    - get: cloudcache-management-cf-plugin
      trigger: true
      passed: [$standalone_targets]
    - get: pcc-$pccver
    - get: stemcell
      params:
        preserve_filename: true
        globs: ["*google*"]
    - get: cloudcache-management-cf-plugin-ci-image
    - get: golang-image
    - put: gcp-env
      tags: [pivotal-internal-worker]
      params:
        action: claim
  - task: install-pcc
    image: cloudcache-management-cf-plugin-ci-image
    config:
      platform: linux
      inputs:
      - name: gcp-env
      - name: pcc-$pccver
      - name: cloudcache-management-cf-plugin
      - name: stemcell
      run:
        path: bash
        args:
        - -exc
        - |
          cd cloudcache-management-cf-plugin
          ci/install.sh -p "../pcc-$pccver" -s "\$(ls ../stemcell/*.tgz)" -g "../gcp-env/metadata"
          ci/create-service.bash -g "../gcp-env/metadata"
  - task: build-plugin
    timeout: 1h
    image: golang-image
    config:
      inputs:
        - name: cloudcache-management-cf-plugin
      outputs:
        - name: pcc-plugin
      platform: linux
      run:
        path: /bin/sh
        args:
          - -ec
          - |
            cd cloudcache-management-cf-plugin
            ./build.sh
            cp pcc ../pcc-plugin/
  - task: smoke-test
    image: cloudcache-management-cf-plugin-ci-image
    config:
      platform: linux
      inputs:
      - name: gcp-env
      - name: cloudcache-management-cf-plugin
      - name: pcc-plugin
      run:
        path: bash
        args:
        - -exc
        - |
          cd cloudcache-management-cf-plugin
          ci/login.bash "../gcp-env/metadata"
          cf="cf pcc test" ci/smoke-test.bash
    ensure:
      aggregate:
      - put: gcp-env
        tags: [pivotal-internal-worker]
        params:
          action: unclaim
          env_file: gcp-env/metadata
EOF
done

fly -t ${TARGET} login --team-name=${TEAM} --concourse-url=https://${TARGET}
fly -t ${TARGET} set-pipeline -p ${PIPELINE} -c pipeline.yml
rm pipeline.yml
