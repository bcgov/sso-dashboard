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

1. `Alloy` & `Loki`: collect, transform and load raw log data for the designated time period.

1. `Loki` & `S3`: provide the Amazon S3 compatible Object Storage to store/read compacted event data by Loki.

1. `Alloy` & `Custom Go server`: collect, and upsert the aggreated event historial data in DB.

1. `Grafana`: connect Loki and the aggregation DB to visualize the logs and stats.

   ![SSO Dashboard Architecture Diagram](assets/sso-dashboard.drawio.svg)

1. Loki in AWS breakdown:

   ![SSO Loki on AWS Diagram](assets/sso-dashboard-aws.drawio.svg)

### Loki in AWS ECS Cluster

Loki has a helm chart for deploying in kubernetes. For the deployment in an ECS cluster there are a few changes to note:

- Service discovery can be used in ECS to replace services in k8s. Since we cannot use this in the BCGov AWS, it has been replaced with a network load balancer. This is necessary to allow read and write tasks to communicate on port 7946. If not working, you will see "empty ring" errors.
- ECS does not support config maps. To replace this a custom image was built with custom configuration files. Configurations that will be changed at runtime can set their values with the syntax ${ENV_VAR:-default}, and environment variables can be used to configure them. Values consistent across environments can be hardcoded.
- The helm chart includes a deployment "gateway". This is an nginx reverse proxy which provides path-based routing to the read and write services. It has been replaced with listener rules on the application load balancer.
- When deploying locally, you will need to use the values from the [terraform workflow file](/.github/workflows/terraform.yaml#97) to populate a var file, refer to the dev or prod block depending on environment. The secret value loki_auth_token can be found in the tools namespace secret loki-auth-token.

<!-- ![image](https://user-images.githubusercontent.com/36021827/211399712-5bbeaa67-2994-460f-a12b-368b13187cdd.png) -->

## Deployment

The helm charts for the alloy instance and grafana dashboard can be installed with make commands. These automate adding environment variables from .env files in their directories. See the directory readmes for more information. They will deploy on merge to dev for sandbox, and main for production.

The Loki setup is deployed with terraform into AWS. It deploys automatically on merge to dev/main.

GitHub CD pipeline scripts are triggered based on the directory that has changed; When deploying for the first time you should deploy alloy last, as it will give not found errors until the receiving resources (loki and aggregator) are up and running.

The terraform account for deployment is restricted to the required resource types for this repository. If adding new resources not currently required, you will get a permission denied error. Expand the permissions on the `sso-dashboard-boundary` as needed.

When doing an initial webhook setup to integrate with [AWS SNS](https://aws.amazon.com/sns) you need to confirm the url you gave is correct. AWS will send a link to the provided URL to confirm. You can find it in the `content_raw.SubscribeURL` parameter to confirm. e.g for rocket chat the script:

``` javascript
class Script {
    process_incoming_request({ request }) {
      return {
        content:{
         text: `@here ${JSON.parse(request.content_raw).SubscribeURL}`
         }
      };
    }
  }
```

Would output the url to follow.

### Grafana Git Sync

Grafana's native Git Sync feature connects a Grafana instance to a GitHub repository so that dashboards are version-controlled. Dashboard changes must go through a pull request and be approved before they are applied to Grafana.

**How it works**

- Grafana polls the configured repository branch at a set interval (currently `3600` seconds).
- A GitHub App installed on the target repository (`bcgov/sso-grafana-dashboards`) authenticates the connection. The DevOps team provides the App ID, Installation ID, and private key.
- The Terraform configuration in [`grafana/git-sync/terraform/`](/grafana/git-sync/terraform/) manages the Git Sync settings via the Grafana API (repository URL, branch, sync path, and sync interval).

**Repository layout**

| Environment | Grafana URL | Repo path synced |
|---|---|---|
| Sandbox | `https://sso-grafana-sandbox.apps.gold.devops.gov.bc.ca` | `apps/sso/environments/sandbox` |
| Production | `https://sso-grafana.apps.gold.devops.gov.bc.ca` | `apps/sso/environments/production` |

**Initial setup steps**

1. **Obtain GitHub App credentials** from the DevOps team:
   - App ID
   - Installation ID
   - Private key (PEM file)

2. **Store the credentials as GitHub Actions secrets** in this repository:
   | Secret name | Value |
   |---|---|
   | `SANDBOX_SSO_GRAFANA_SERVICE_ACCOUNT_TOKEN` | Grafana service account token for sandbox |
   | `PROD_SSO_GRAFANA_SERVICE_ACCOUNT_TOKEN` | Grafana service account token for production |
   | `TERRAFORM_DEPLOY_ROLE_ARN_DEV` | AWS IAM role ARN for sandbox Terraform state |
   | `TERRAFORM_DEPLOY_ROLE_ARN_PROD` | AWS IAM role ARN for production Terraform state |

3. **Run the GitHub Action** at [`.github/workflows/terraform-git-sync.yaml`](/.github/workflows/terraform-git-sync.yaml):
   - Navigate to **Actions → Terraform For Git Sync → Run workflow**.
   - Select the target environment (`sandbox` or `production`).
   - The workflow authenticates to AWS, initialises Terraform against the S3 backend, and applies the Git Sync configuration to the chosen Grafana instance.

4. **Verify** by opening the Grafana instance, navigating to **Administration → Provisioning**, and confirming the repository is connected and the last sync succeeded.

**Making dashboard changes**

1. Edit or add dashboard JSON files in `bcgov/sso-grafana-dashboards` under the appropriate environment path.
2. Open a pull request targeting the `main` branch.
3. Once approved and merged, Grafana will pick up the changes on the next sync cycle (within 1 hour).

## Service accounts

Service accounts are already generated and added to github secrets. If needing to recreate the service account, see the [service-account-generator directory](/service-account-generator/README.md) for how to do so
