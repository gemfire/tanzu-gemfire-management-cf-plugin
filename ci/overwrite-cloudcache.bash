#!/usr/bin/env bash
set -e
gcp_metadata_path=${1}

env_name=$(jq -r .name < "${gcp_metadata_path}")
ops_man_user="$(jq -r .ops_manager.username < "${gcp_metadata_path}")"
ops_man_url="$(jq -r .ops_manager.url < "${gcp_metadata_path}")"
ops_man_password="$(jq -r .ops_manager.password < "${gcp_metadata_path}")"

om_exec() {
   om --target "${ops_man_url}" --username "${ops_man_user}" --password "${ops_man_password}"  --skip-ssl-validation "$@"
   return $?
}


function login_and_target_cf_space() {
  echo "getting cf guid"
  cf_guid="$(om_exec curl -s -x GET --path /api/v0/deployed/products -s | jq -r 'map(select(.type=="cf")) | .[].guid')"
  if [[ $? -ne 0 ]] || [[ "$cf_guid" == "" ]]; then
    error "failed to get the cf guid"
    exit 1
  fi

  echo "getting auth info"
  uaa_admin_password="$(om_exec  curl -s -x GET --path /api/v0/deployed/products/${cf_guid}/credentials/.uaa.admin_credentials -s | jq -r .credential.value.password)"
  if [[ $? -ne 0 ]] || [[ "$uaa_admin_password" == "" ]]; then
    error "failed to get the uaa_admin_password"
    exit 1
  fi

  echo "cf_guid = ${cf_guid}, password=${uaa_admin_password}"

  cf login -a https://api.sys.${env_name}.cf-app.com -u "admin" -p "${uaa_admin_password}" -o "system" -s "test_space" --skip-ssl-validation
}


function overwrite() {
  echo "Overwrite cloudcache-metrics-release on VMs"
  export GUID=$(cf service test --guid)
  for i in `bosh -d service-instance_${GUID} vms | cut -f1`;
  do

    echo "push cloudcache-metrics-release jar file to vm ${i}"
    bosh -d service-instance_${GUID} scp cloudcache-metrics-jar/cloudcache-metrics.jar ${i}:/tmp/
    bosh -d service-instance_${GUID} scp cloudcache-metrics-source/ci/remote-upgradecloudcachemetrics.bash ${i}:/tmp/

    echo "Start actual overwrite of cloudcache-metrics.jar"
    bosh -d service-instance_${GUID} ssh ${i} "sudo bash /tmp/remote-upgradecloudcachemetrics.bash"
  done
}

login_and_target_cf_space
overwrite

set +e