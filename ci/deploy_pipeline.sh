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
TEAM=main
PIPELINE=tanzu-gemfire-management-cf-plugin
BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Poolsmiths configuration
POOL_NAME=us_2_9
POOLSMITHS_API_TOKEN=0d82e637-6681-4d4a-9e9f-90a71db5de0d

# To get PCC and stemcell snapshots and releases
PIVNET_API_TOKEN=c90e06904710409eb60d55459e3b3dbd-r

# To get access to latest stemcells
PIVNET_API_TOKEN_LATEST=4b4b8fadaac9460dbdc0f6dd08c2939a-r

# Blessed PCC tile from https://pcc1.ci.cf-app.com/teams/main/pipelines/cloudcache-1.10.x
BLESSED_KEY="AKIAJPMW5"G'JQJPVJVR6Q'
BLESSED_SECRET='n41r3zdlwKCtln'"w22plhTJUhE"/9iBWQmh4p26fPY

# Version(s) of GemFire for stand-alone testing, whitespace-separated
# 3-digit indicates a released version, anything else is latest nightly build"
STANDALONE_GEMFIRE_VERSIONS="9.9.5 9.10.6 9.10 1.13 1.14 develop"

# Version(s) of PCC+stemcell for testing as plugin, whitespace-separated.
# Last one in the list will be taken from blessed bucket, the rest from pivnet
# Special string latest gives latest stemcell, for example 1.9+456 1.10+latest
# See https://docs.google.com/spreadsheets/d/1iYp71cfXVXCeJF5mm9Wh6KoVCjjm64eAICprjwSK1Zk/edit and https://bosh.io/stemcells/
PCC_VERSIONS="1.11+621 1.12+621 1.13+latest"

cat << EOF > pipeline.yml
---
resources:
EOF
for gem in $STANDALONE_GEMFIRE_VERSIONS; do
  if [ $(echo $gem|tr . '\n' | wc -l) -eq 3 ] ; then
    mm=${gem%.*}
    cat << EOF >> pipeline.yml
- name: gemfire-$gem
  type: s3
  icon: file-cloud
  source:
    bucket: gemfire-releases
    region_name: ((aws-default-region))
    access_key_id: ((gemfire-aws-access-key-id))
    secret_access_key: ((gemfire-aws-secret-access-key))
    versioned_file: ${mm}/${gem}/pivotal-gemfire-${gem}.tgz
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
    initial_path: foo
EOF
  fi
done
for pccstemvers in $PCC_VERSIONS; do
  latestpccver="${pccstemvers%+*}"
done
for pccstemvers in $PCC_VERSIONS; do
  pccver="${pccstemvers%+*}"
  stemcellver="${pccstemvers#*+}"
  regex=$(echo "${pccver}." | sed 's#\.#\\\.#')
  [ "$pccver" = "$latestpccver" ] && getfrompivnet=false || getfrompivnet=true
  cat << EOF >> pipeline.yml
- name: pcc-$pccver
EOF
  if [ "$getfrompivnet" = "false" ] ; then
    cat << EOF >> pipeline.yml
  type: s3
  source:
    bucket: cloudcache-eng
    access_key_id: $BLESSED_KEY
    secret_access_key: $BLESSED_SECRET
    private: true
    regexp: blessed-p-cloudcache-tile/p-cloudcache-(${regex}.*).pivotal

EOF
  else
    cat << EOF >> pipeline.yml
  type: pivnet
  source:
    api_token: $PIVNET_API_TOKEN
    product_slug: p-cloudcache
    product_version: ${regex}.*
EOF
  fi
  cat << EOF >> pipeline.yml
- name: pcc-env-$pccver
  type: pcf-pool
  tags: [nimbus]
  source:
    api_token: $POOLSMITHS_API_TOKEN
    hostname: environments.toolsmiths.cf-app.com
    pool_name: $POOL_NAME
EOF
done
cat << EOF >> pipeline.yml
- name: tanzu-gemfire-management-cf-plugin
  type: git
  source:
    uri: git@github.com:gemfire/tanzu-gemfire-management-cf-plugin.git
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

- name: tanzu-gemfire-management-cf-plugin-ci-dockerfile
  type: git
  source:
    uri: git@github.com:gemfire/tanzu-gemfire-management-cf-plugin.git
    branch: $BRANCH
    paths:
      - "ci/docker/*"
    private_key: ((!gemfire-ci-private-key))

EOF
for stemcellver in $(for pccstemvers in $PCC_VERSIONS; do echo "${pccstemvers#*+}"; done | sort -u); do
  cat << EOF >> pipeline.yml
- name: stemcell-$stemcellver
  type: pivnet
  source:
    product_slug: stemcells-ubuntu-xenial
EOF
  if [ "$stemcellver" = "latest" ] ; then
    cat << EOF >> pipeline.yml
    api_token: $PIVNET_API_TOKEN
    sort_by: semver

EOF
  else
    cat << EOF >> pipeline.yml
    api_token: $PIVNET_API_TOKEN_LATEST
    product_version: ${stemcellver}\..*

EOF
  fi
done
cat << EOF >> pipeline.yml
- name: tanzu-gemfire-management-cf-plugin-ci-image
  type: docker-image
  source:
    username: "_json_key"
    password: ((!concourse-gcp-key))
    repository: gcr.io/gemfire-dev/tanzu-gemfire-management-cf-plugin-ci

- name: weekly
  type: time
  source:
    start: 3:00 AM
    stop: 10:00 AM
    days: [Monday]
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
- name: build-tanzu-gemfire-management-cf-plugin
  serial: true
  plan:
  - in_parallel:
    - get: tanzu-gemfire-management-cf-plugin
      trigger: true
    - get: golang-image
  - task: ginkgo
    timeout: 1h
    image: golang-image
    config:
      inputs:
        - name: tanzu-gemfire-management-cf-plugin
      platform: linux
      run:
        path: /bin/sh
        args:
          - -ec
          - |
            cd tanzu-gemfire-management-cf-plugin
            go get github.com/onsi/ginkgo/ginkgo
            ginkgo -r

- name: build-docker-image
  plan:
  - get: weekly
    trigger: true
  - get: tanzu-gemfire-management-cf-plugin-ci-dockerfile
    trigger: true
  - put: tanzu-gemfire-management-cf-plugin-ci-image
    params:
      build: tanzu-gemfire-management-cf-plugin-ci-dockerfile/ci/docker
      tag_as_latest: true
EOF
for gem in $STANDALONE_GEMFIRE_VERSIONS; do
  cat << EOF >> pipeline.yml
- name: test-cloudcache-management-cf-standalone-$gem
  serial: true
  plan:
  - in_parallel:
    - get: tanzu-gemfire-management-cf-plugin
      trigger: true
      passed: [build-tanzu-gemfire-management-cf-plugin]
    - get: golang-image
    - get: gemfire-$gem
EOF
if [ $(echo $gem|tr . '\n' | wc -l) -ne 3 ] ; then
  cat << EOF >> pipeline.yml
      trigger: true
EOF
fi
function buildPlugin {
  cat << EOF >> pipeline.yml
  - task: build-plugin
    timeout: 1h
    image: golang-image
    config:
      inputs:
        - name: tanzu-gemfire-management-cf-plugin
      outputs:
        - name: pcc-plugin
      platform: linux
      run:
        path: /bin/sh
        args:
          - -ec
          - |
            cd tanzu-gemfire-management-cf-plugin
            ./build.sh
            cp gemfire ../pcc-plugin/
EOF
}
buildPlugin
cat << EOF >> pipeline.yml
  - task: standalone-test
    timeout: 1h
    config:
      image_resource:
        type: docker-image
        source:
          repository: openjdk
          tag: 8
      inputs:
        - name: tanzu-gemfire-management-cf-plugin
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
            cd tanzu-gemfire-management-cf-plugin
            gemfire="../pcc-plugin/gemfire"
            \$gemfire --help
            cf="\$gemfire http://localhost:7070" ci/smoke-test.bash
EOF
done
for gem in $STANDALONE_GEMFIRE_VERSIONS ; do
  standalone_targets="$standalone_targets,test-cloudcache-management-cf-standalone-$gem"
done
standalone_targets=$(echo "$standalone_targets"|cut -c2-)
for pccstemvers in $PCC_VERSIONS; do
  pccver="${pccstemvers%+*}"
  stemcellver="${pccstemvers#*+}"
  cat << EOF >> pipeline.yml
- name: test-cloudcache-management-cf-pcc-$pccver
  plan:
  - aggregate:
    - get: tanzu-gemfire-management-cf-plugin
      trigger: true
      passed: [$standalone_targets]
    - get: pcc-$pccver
    - get: stemcell-$stemcellver
      params:
        preserve_filename: true
        globs: ["*google*"]
    - get: tanzu-gemfire-management-cf-plugin-ci-image
    - get: golang-image
    - put: pcc-env-$pccver
      tags: [nimbus]
      params:
        action: claim
EOF
buildPlugin
cat << EOF >> pipeline.yml
  - task: install-pcc
    image: tanzu-gemfire-management-cf-plugin-ci-image
    config:
      platform: linux
      inputs:
      - name: pcc-env-$pccver
      - name: pcc-$pccver
      - name: tanzu-gemfire-management-cf-plugin
      - name: stemcell-$stemcellver
      run:
        path: bash
        args:
        - -exc
        - |
          cd tanzu-gemfire-management-cf-plugin
          ci/install.sh -p ../pcc-$pccver -s ../stemcell-$stemcellver/*.tgz -g ../pcc-env-$pccver/metadata
          ci/create-service.bash -g ../pcc-env-$pccver/metadata
  - task: plugin-test
    image: tanzu-gemfire-management-cf-plugin-ci-image
    config:
      platform: linux
      inputs:
      - name: pcc-env-$pccver
      - name: tanzu-gemfire-management-cf-plugin
      - name: pcc-plugin
      run:
        path: bash
        args:
        - -exc
        - |
          cd tanzu-gemfire-management-cf-plugin
          ci/login.bash ../pcc-env-$pccver/metadata
          cf="cf gemfire test" ci/smoke-test.bash
    ensure:
      aggregate:
      - put: pcc-env-$pccver
        tags: [nimbus]
        params:
          action: unclaim
          env_file: pcc-env-$pccver/metadata
EOF
done

fly -t ${TARGET}-${TEAM} login --team-name=${TEAM} --concourse-url=https://${TARGET}
fly -t ${TARGET}-${TEAM} set-pipeline -p ${PIPELINE} -c pipeline.yml
rm pipeline.yml
