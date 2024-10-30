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

The helm charts for the promtail instances and grafana dashboard can be installed with make commands. These automate adding environment variables from .env files in their directories. See the directory readmes for more information.

The Loki setup is deployed with terraform into AWS. It deploys automatically on merge to dev/main.

GitHub CD pipeline scripts are triggered based on the directory that has changed; there is a recommended deployment order when deploying the resources for the very first time:

1. `Loki`: deploys the `MinIO` and `Loki` resources, `read`, `write`, and `gateway`.
1. `Aggregator`: deploys the `Aggregator` and `Compactor` with the `Postgres DB`.
1. `Grafana`: deploys the `Grafana` dashboard with the two `datasources` configured above.
1. `Promtail`: deploys the `Promtail` in multiple namespaces to collect the Keycloak disk logs.

## GitHub secrets

The following secrets are set in the GitHub secrets of the repository and can be found in [OCP secret](https://console.apps.silver.devops.gov.bc.ca/k8s/ns/6d70e7-tools/secrets/sso-team-sso-dashboard-github-secrets)

### Sandbox

- `SANDBOX_OPENSHIFT_SERVER`: the OpenShift online server URL.
- `SANDBOX_OPENSHIFT_TOKEN`: : the OpenShift session token.
  - please the find the secret in [Sandbox Deployer Secret](https://console.apps.gold.devops.gov.bc.ca/k8s/ns/c6af30-tools/secrets/oc-deployer-token-9tgwm)
- `SANDBOX_OPENSHIFT_NAMESPACE`: the namespace name to deploy `Grafana`, `Loki`, and `Aggregator`.
- `SANDBOX_SSO_CLIENT_ID`: the SSO integration credentials, `client id`, to set in `Grafana` and `MinIO` dashboard UI.
- `SANDBOX_SSO_CLIENT_SECRET`: the SSO integration credentials, `client secret`, to set in `Grafana` and `MinIO` dashboard UI.
  - please find the integration `#4492 SSO Dashboard` via [CSS app](https://bcgov.github.io/sso-requests)

### Production

- `PROD_OPENSHIFT_SERVER`: the OpenShift online server URL.
- `PROD_OPENSHIFT_TOKEN`: : the OpenShift session token.
  - please the find the secret in [Sandbox Deployer Secret](https://console.apps.gold.devops.gov.bc.ca/k8s/ns/eb75ad-tools/secrets/oc-deployer-token-b99cz)
- `PROD_OPENSHIFT_NAMESPACE`: the namespace name to deploy `Grafana`, `Loki`, and `Aggregator`.
- `PROD_SSO_CLIENT_ID`: the SSO integration credentials, `client id`, to set in `Grafana` and `MinIO` dashboard UI.
- `PROD_SSO_CLIENT_SECRET`: the SSO integration credentials, `client secret`, to set in `Grafana` and `MinIO` dashboard UI.
  - please find the integration `#4492 SSO Dashboard` via [CSS app](https://bcgov.github.io/sso-requests)
