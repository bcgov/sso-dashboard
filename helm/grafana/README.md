# sso-grafana

A Helm chart for deploying [Grafana dashboard](https://github.com/grafana/helm-charts/tree/main/charts/grafana).

## Local deployment via Helm chart

### Pre-Requisites

#### Setup Network Policies

- Below network policy has to be added to all the namespaces, where grafana requires to access the data source

```sh
export LICENSE_PLATE=
export ENV=

# run below command after logging into each namespace (dev, test and prod)
envsubst < net-policy-sso-keycloak.yaml | oc apply -f -
```

#### Update Helm Values

- create `.env` from `.env.example` and fill the values

### Installing/Upgrading the Chart

```sh
make upgrade
```

- please find the SSO client credentials of the integration `#4492 SSO Dashboard` via [CSS app](https://bcgov.github.io/sso-requests):

### Uninstalling the Chart

```sh
make uninstall
```

## Grafana Admin credentials

once the deployment is completed, please find the Grafana admin credentials in OCP secrets below:

- `dev`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/e4ca1d-tools/secrets/sso-grafana
- `prod`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/eb75ad-tools/secrets/sso-grafana

## Configuration

please find the full list of `Grafana Helm values` configuration in [Grafana Helm values](https://github.com/grafana/helm-charts/blob/main/charts/grafana/values.yaml)
