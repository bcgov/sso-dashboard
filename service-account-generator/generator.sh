#!/bin/bash
set -e

usage() {
    cat <<EOF
Creates a service account for the dev test and prod environments of the project with
namespace licence plate arg.

Usages:
    $0 <project_licence_plate>

Available licence plates:
    - e4ca1d
    - eb75ad

Examples:
    $ $0 e4ca1d
EOF
}

if [ "$#" -lt 1 ]; then
    usage
    exit 1
fi

licence_plate=$1

# create service account in prod
oc -n "$licence_plate"-prod create sa sso-dashboard-deployer-"$licence_plate"

create_role_and_binding() {
  if [ "$#" -lt 2 ]; then exit 1; fi
  licence_plate=$1
  env=$2
  namespace="$licence_plate-$env"

  oc process -f ./templates/role.yaml -p NAMESPACE="$namespace" | oc -n "$namespace" apply -f -

  oc -n "$namespace" create rolebinding sso-dashboard-deployer-role-binding-"$namespace"   \
  --role=sso-dashboard-deployer-"$namespace" \
  --serviceaccount="$licence_plate"-prod:sso-dashboard-deployer-"$licence_plate"
}

# for dev, test, prod and tools create the role and role binding
create_role_and_binding "$licence_plate" "tools"

create_role_and_binding "$licence_plate" "prod"

create_role_and_binding "$licence_plate" "test"

create_role_and_binding "$licence_plate" "dev"
