#!/usr/bin/env bash

generate_secrets() {
  cat << EOF
---
poolsmiths-api-token: $(lpass show 'Shared-PCC Observability/Poolsmiths' --notes)
pivnet-api-token: $(lpass show 'Shared-PCC Observability/Pivnet' --notes)
EOF
}

pipelinename=${1:-"cloudcache-metrics-release-ci"}

#figure out what target to use or create one.
target=$(fly targets | grep concourse.gemfire-ci.info  | tail -1 | awk '{print $1}')
if [[ "${target}" == "" ]]; then
  fly login -c https://concourse.gemfire-ci.info  -t concourse-gemfire-ci-info
  target=$(fly targets | grep concourse.gemfire-ci.info  | tail -1 | awk '{print $1}')
fi


fly --target ${target} set-pipeline -p ${pipelinename}   \
  --config pipeline.yml \
  --load-vars-from <(generate_secrets)
