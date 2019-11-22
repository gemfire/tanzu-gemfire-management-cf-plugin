#!/usr/bin/env bash

set -e
product_name="p-cloudcache"

while getopts "p:s:g:" opt; do
  case ${opt} in
    p) product_dir_path="${OPTARG}" ;;
    s) stemcell_path="${OPTARG}" ;;
    g) gcp_metadata_path="${OPTARG}" ;;
  esac
done

product_path="$(ls ${product_dir_path}/*.pivotal)"
product_version="$(cat ${product_dir_path}/metadata.json | jq  -r '.ProductFiles[] | select(.File | contains("p-cloudcache") ) | .FileVersion')"

env_name=$(jq -r .name < "${gcp_metadata_path}")
ops_man_user="$(jq -r .ops_manager.username < "${gcp_metadata_path}")"
ops_man_url="$(jq -r .ops_manager.url < "${gcp_metadata_path}")"
ops_man_password="$(jq -r .ops_manager.password < "${gcp_metadata_path}")"

pas_network=$(cat ${gcp_metadata_path} | jq -r '.ert_subnet')
services_network=$(cat ${gcp_metadata_path} | jq -r '.service_subnet_name')

declare -a availability_zones=($(cat ${gcp_metadata_path} | jq -r '.azs[]'))

availability_zones_for_gcp=$(printf '%s\n' "${availability_zones[@]}" | jq -R . | jq -s .)

om_exec() {
   om --target "${ops_man_url}" --username "${ops_man_user}" --password "${ops_man_password}"  --skip-ssl-validation "$@"
   return $?
}


function pid_exists() {
  ps -p "${1}" &> /dev/null
}


function upload_tile() {
  om_exec upload-product --product "${product_path}"
}


function wait_for_upload_to_finish() {
  sleep 1
  pid=$(ps -ef | grep 'om --target' | head -n 1 | awk '{print $2}')
  SECONDS=0
  while pid_exists "${pid}"; do
    echo ""
    echo "first attempt: waiting for tile to upload "${SECONDS}"s (timeout: 600s)"
    if [[ "${SECONDS}" -gt 600 ]]; then
      kill -9 "${pid}"
      om_exec upload-product --product "${product_path}"
    else
      echo -n .
      sleep 5
    fi
  done
}


function stage_tile() {
  echo ""
  echo "Staging ${product_name}"
  om_exec stage-product --product-name "${product_name}" --product-version "${product_version}"
}


function configure_tile() {
  echo ""
  echo "Configuring ${product_name}'s product properties and network settings"
  export OM_VAR_product_name=${product_name}
  export OM_VAR_availability_zones_for_gcp=${availability_zones_for_gcp}
  export OM_VAR_availability_zone_0=${availability_zones[0]}
  export OM_VAR_availability_zone_1=${availability_zones[1]}
  export OM_VAR_availability_zone_2=${availability_zones[2]}
  export OM_VAR_pas_network=${pas_network}
  export OM_VAR_services_network=${services_network}

  om_exec configure-product --config ci/configure-product-template.yml --vars-env OM_VAR
}


function upload_stemcell() {
  echo ""
  echo "Uploading the stemcell"
  om_exec upload-stemcell --stemcell "${stemcell_path}"
}


function deploy_tile() {
  pcc_deployment_guid=$(om_exec curl -p /api/v0/stemcell_assignments 2> /dev/null | jq  -r '.products[] | select( .identifier == "p-cloudcache" ) | .guid')

  echo ""
  echo "Deploying PCC product for deployment = ${pcc_deployment_guid}"
  trigger_json=$(om_exec curl -x POST --path "/api/v0/installations" -d '{"deploy_products": ["'"$pcc_deployment_guid"'"], "ignore_warnings": true }' 2> /dev/null)
  install_id=$(echo ${trigger_json} | jq .install.id)
}


function wait_for_deploy_to_finish() {
  done=false
  echo ""
  echo "Waiting for deploy to succeed for install id ${install_id}"
  echo "Expect to wait 25 minutes or more"
  while [[ ${done} != true ]]; do
    sleep 10
    status=$(om_exec curl --path "/api/v0/installations/${install_id}" 2> /dev/null | jq -r .status )
    if [[ ${status} != "running" ]]; then
      done=true
    fi
    echo -n .
  done

  if [[ ${status} = "succeeded" ]]; then
    echo "done!"
  else
    echo "deploy did not succeed: ${status}"
    exit 1
  fi
}

upload_tile
wait_for_upload_to_finish
stage_tile
configure_tile
upload_stemcell
deploy_tile
wait_for_deploy_to_finish

set +e
exit 0