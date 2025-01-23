# Generating service accounts for the CICD pipeline

The github actions need service accounts to run. The script `generate_sa.sh` will create a service acount for the prod environment of a given openshift project and give that account the roles in the dev, test, and prod environments for deploying the repositorie's resources.

## Generate the service accounts

While logged into the **Gold** instance run:

`
./generate_sa.sh <<LICENCE_PLATE>>
`

The service account, roles, and rolebindings will be created.

## Update the github action secrets

Each service account will generate a secret in the `-prod` namespace with the name `sso-dashboard-deployer-<<LICENCE_PLATE>>-token-#####`.  Copy this token into the GithHub secrets on this repos:

SANDBOX_OPENSHIFT_TOKEN
PROD_OPENSHIFT_TOKEN
