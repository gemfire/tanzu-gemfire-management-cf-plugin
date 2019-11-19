#!/usr/bin/env bash

set -e

service_name=service-metrics

function stop_service_synchronously() {
  echo "stopping cloudcache-metrics"

  /var/vcap/bosh/bin/monit stop ${service_name}

  echo "waiting for ${service_name} to stop"

  until /var/vcap/bosh/bin/monit summary | grep ${service_name} | grep "not monitored$" > /dev/null 2>&1;
  do
    sleep 1
    echo -n .
  done
}

function overwrite_service() {
  echo "overwriting cloudcache-metrics"
  cp /tmp/cloudcache-metrics.jar /var/vcap/packages/cloudcache-metrics/cloudcache-metrics.jar
}

function restart_service() {
  /var/vcap/bosh/bin/monit start ${service_name}

  echo "waiting for ${service_name} to start"

  until /var/vcap/bosh/bin/monit summary | grep ${service_name} | grep "running$" > /dev/null 2>&1;
  do
    sleep 1
    echo -n .
  done
}

stop_service_synchronously
overwrite_service
restart_service

set +e