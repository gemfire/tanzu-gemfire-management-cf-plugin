#!/usr/bin/env bash

set -e
product_name="p-cloudcache"

while getopts "p:s:g:" opt; do
  case ${opt} in
    g) gcp_metadata_path="${OPTARG}" ;;
  esac
done

env_name=$(jq -r .name < "${gcp_metadata_path}")
ops_man_user="$(jq -r .ops_manager.username < "${gcp_metadata_path}")"
ops_man_url="$(jq -r .ops_manager.url < "${gcp_metadata_path}")"
ops_man_password="$(jq -r .ops_manager.password < "${gcp_metadata_path}")"

cf_org="system"
cf_space="test_space"
service_instance_name="test"


om_exec() {
   om --target "${ops_man_url}" --username "${ops_man_user}" --password "${ops_man_password}"  --skip-ssl-validation "$@"
   return $?
}


function login_to_cf_cli() {
   cf_guid="$(om_exec curl -s -x GET --path /api/v0/deployed/products -s 2> /dev/null | jq -r 'map(select(.type=="cf")) | .[].guid')"
   uaa_admin_password="$(om_exec  curl -s -x GET --path /api/v0/deployed/products/${cf_guid}/credentials/.uaa.admin_credentials -s 2> /dev/null| jq -r .credential.value.password)"
   cf login -a https://api.sys.${env_name}.cf-app.com -u "admin" -p "${uaa_admin_password}" --skip-ssl-validation
}


function create_space() {

  echo "Creating space ${cf_space} in org ${cf_org}."
  cf create-space ${cf_space} -o ${cf_org}
  cf target -o ${cf_org} -s ${cf_space}
}


function create_service_instance() {
  local plan

  echo "Creating service"
  #expecting plan to be smallest (non-dev) plan. Currently extra-small
  plan="$(cf m | grep p-cloudcache | awk '{print $2}' | cut -d ',' -f 1)"

  cf create-service p-cloudcache "$plan" "${service_instance_name}"
}


function wait_for_create_service() {
  counter=0
  status="$(cf service "${service_instance_name}" | grep "status:")"
  echo ""
  echo "waiting for create service to complete. It should take about 12 minutes"

  while [[ "${status}" == *"status:"* ]] && [[ "${counter}" -lt 240 ]]; do
    if [[ "${status}" == *"status:"*"create failed"* ]]; then
      error "Service instance failed to create"
    elif [[ "${status}" == *"status:"*"create succeeded"* ]]; then
      echo ""
      echo "finished!"
      break
    fi
    sleep 5
    echo -n .
    status="$(cf service "${service_instance_name}" | grep "status:" )"
    let counter=${counter}+1
  done

}

login_to_cf_cli
create_space
create_service_instance
wait_for_create_service
cf create-service-key "${service_instance_name}" myKey

set +e
exit 0