# sso-dashboard

SSO Keycloak dashboard services provide the ability to monitor real-time statistical data and event logs.

## Local Development Environment

- Install asdf
- Run `make local-setup` to install necessary tooling

## Benefits

1. De-coupling the auditing service from the authentication service (Keycloak) and reducing the amount of Keycloak SQL transactions and DB data storage; gives better maintainability of the Keycloak instances.

   - see [Impact of enabling Audit Events for Keycloak](https://keycloak.discourse.group/t/impact-of-enabling-audit-events-for-keycloak/13552/2)

1. Full control of the log ingestion and data store process that gives better performance displaying the dashboard data and log events in a separate business intelligent tool rather than in Keycloak UI.

## Technical Considerations

1. an access to the Keycloak logs without significant impacts on Keycloak operational performance.

1. a functional log consumer that can be used to filter the logs and extract metadata before the data stored.

1. a solution to store the aggregated historical data and logs for a longer term.

1. a dashboard tool to display the aggregated data and option to search log events.

1. a dashboard that has authorization integration to support multi-tenant workspaces.

## Architectural Design

1. `Promtail` & `Loki`: collect, transform and load raw log data for the designated time period.

1. `Loki` & `S3`: provide the Amazon S3 compatible Object Storage to store/read compacted event data by Loki.

1. `Promtail` & `Custom Go server`: collect, and upsert the aggreated event historial data in DB.

1. `Grafana`: connect Loki and the aggregation DB to visualize the logs and stats.

   ![SSO Dashboard Architecture Diagram](assets/sso-dashboard.drawio.svg)

1. Loki in AWS breakdown:

   ![SSO Loki on AWS Diagram](assets/sso-dashboard-aws.drawio.svg)

### Loki in AWS ECS Cluster

Loki has a helm chart for deploying in kubernetes. For the deployment in an ECS cluster there are a few changes to note:

- Service discovery can be used in ECS to replace services in k8s. Since we cannot use this in the BCGov AWS, it has been replaced with a network load balancer. This is necessary to allow read and write tasks to communicate on port 7946. If not working, you will see "empty ring" errors.
- ECS does not support config maps. To replace this a custom image was built with custom configuration files. Configurations that will be changed at runtime can set their values with the syntax ${ENV_VAR:-default}, and environment variables can be used to configure them. Values consistent across environments can be hardcoded.
- The helm chart includes a deployment "gateway". This is an nginx reverse proxy which provides path-based routing to the read and write services. It has been replaced with listener rules on the application load balancer.

<!-- ![image](https://user-images.githubusercontent.com/36021827/211399712-5bbeaa67-2994-460f-a12b-368b13187cdd.png) -->

## Deployment

The helm charts for the promtail instances and grafana dashboard can be installed with make commands. These automate adding environment variables from .env files in their directories. See the directory readmes for more information. They will deploy on merge to dev for sandbox, and main for production.

The Loki setup is deployed with terraform into AWS. It deploys automatically on merge to dev/main.

GitHub CD pipeline scripts are triggered based on the directory that has changed; When deploying for the first time you should deploy promtail last, as it will give not found errors until the receiving resources (loki and aggregator) are up and running.

The terraform account for deployment is restricted to the required resource types for this repository. If adding new resources not currently required, you will get a permission denied error. Expand the permissions on the `sso-dashboard-boundary` as needed.

## Service accounts

Service accounts are already generated and added to github secrets, see below for the related OC secret to see the token value. If needing to recreate the service account, see the [service-account-generator directory](/service-account-generator/README.md) for how to do so.

## GitHub secrets

The following secrets are set in the GitHub secrets of the repository and can be found in [OCP secret](https://console.apps.silver.devops.gov.bc.ca/k8s/ns/6d70e7-tools/secrets/sso-team-sso-dashboard-github-secrets)

### Sandbox

- `SANDBOX_OPENSHIFT_SERVER`: the OpenShift online server URL.
- `SANDBOX_OPENSHIFT_TOKEN`: The OpenShift session token. The token can be found in the sso-dashboard-deployer-e4ca1d-token secret in the prod namespace.
- `GRAFANA_SANDBOX_ENV`: Contains all secrets necessary to deploy grafana as an env file, see [the example env file](/helm/grafana/.env.example) for the list. The values are saved in the openshift secret sso-grafana-env in the tools namespace for reference.

### Production

- `PROD_OPENSHIFT_SERVER`: the OpenShift online server URL.
- `PROD_OPENSHIFT_TOKEN`: The OpenShift session token. The token can be found in the sso-dashboard-deployer-eb75ad-token secret in the prod namespace.
- `GRAFANA_PROD_ENV`: Contains all secrets necessary to deploy grafana as an env file, see [the example env file](/helm/grafana/.env.example) for the list. The values are saved in the openshift secret sso-grafana-env in the tools namespace for reference.
