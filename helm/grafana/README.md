# sso-grafana

A Helm chart for deploying [Grafana dashboard](https://github.com/grafana/helm-charts/tree/main/charts/grafana).

## Local deployment via Helm chart

### Installing/Upgrading the Chart

```sh
make upgrade NAMESPACE=<namespace> \
             SSO_CLIENT_ID=<sso-client-id> \
             SSO_CLIENT_SECRET=<sso-client-secret> \
             AGGREGATION_DB_USER=<aggregation-db-user> \
             AGGREGATION_DB_PASS=<aggregation-db-user>
```

- please find the SSO client credentials of the integration `#4492 SSO Dashboard` via [CSS app](https://bcgov.github.io/sso-requests):

- please find the DB admin credentials in OCP secrets:

  - `dev`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/c6af30-prod/secrets/sso-aggregator-patroni-appusers
  - `prod`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/eb75ad-prod/secrets/sso-aggregator-patroni-appusers

### Uninstalling the Chart

```sh
make uninstall NAMESPACE=<namespace>
```

## Grafana Admin credentials

once the deployment is completed, please find the Grafana admin credentials in OCP secrets below:

- `dev`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/c6af30-prod/secrets/sso-grafana
- `prod`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/eb75ad-prod/secrets/sso-grafana

## Configuration

please find the full list of `Grafana Helm values` configuration in [Grafana Helm values](https://github.com/grafana/helm-charts/blob/main/charts/grafana/values.yaml)
