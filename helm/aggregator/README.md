# sso-aggregator

A Helm chart for deploying [Keycloak event log aggregator](../../aggregator).

## Local deployment via Helm chart

### Installing/Upgrading the Chart

```sh
make upgrade NAMESPACE=<namespace> IMAGE_TAG=<image-tag>
```

please find the published image tags in [SSO Aggregator package](https://github.com/bcgov/sso-dashboard/pkgs/container/sso-aggregator).

### Uninstalling the Chart

```sh
make uninstall NAMESPACE=<namespace>
```

## Database Admin credentials

once the deployment is completed with the patroni database created, please find the DB admin credentials in OCP secrets below to be used for Grafana datasource configuration:

- `dev`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/c6af30-prod/secrets/sso-aggregator-patroni-appusers
- `prod`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/eb75ad-prod/secrets/sso-aggregator-patroni-appusers
