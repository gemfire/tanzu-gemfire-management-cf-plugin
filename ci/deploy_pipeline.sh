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

TARGET=concourse.gemfire-ci.info
URL=https://${TARGET}
TEAM=developer
pipelinename=cloudcache-management-cf-plugin

get_secrets_from_lastpass() {
  cat << EOF
---
poolsmiths-api-token: $(lpass show 'Shared-PCC Observability/Poolsmiths' --notes)
pivnet-api-token: $(lpass show 'Shared-PCC Observability/Pivnet' --notes)
branch: develop
EOF
}

hardcoded_secrets() {
  cat << EOF
---
poolsmiths-api-token: 0d82e637-6681-4d4a-9e9f-90a71db5de0d
pivnet-api-token: c90e06904710409eb60d55459e3b3dbd-r
branch: pcc-integration-test-pipeline
EOF
}

fly -t ${TARGET} login --team-name=${TEAM} --concourse-url=${URL}
fly -t ${TARGET} set-pipeline -p ${pipelinename} -c pipeline.yml \
  --load-vars-from <(hardcoded_secrets)