#!/usr/bin/env bash
set -ex -o pipefail

service_instance_name="test"

function expect {
  set +x
  tee result.json
  for s; do
    if ! grep -q "$s" result.json ; then
      echo "Test Failed: expected '$s' but it was not found" 1>&2
      exit 1
    fi
  done
  set -x
}

$cf commands
$cf commands | grep '^list *[^ ]*s\b' | sed 's/.--.*//' | while read cmd; do
  $cf $cmd | expect '"statusCode": "OK"'
done
