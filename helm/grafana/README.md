# sso-grafana

A Helm chart for deploying [Grafana dashboard](https://github.com/grafana/helm-charts/tree/main/charts/grafana).

## Local deployment via Helm chart

### Pre-Requisites

#### Setup Network Policies

- Below network policy has to be added to all the namespaces, where grafana requires to access the data source

```yaml
# Update $LICENSE_PLATE (ex.: e4ca1d)

kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: sso-dev-sandbox-gold-grafana-access
  namespace: $LICENSE_PLATE-(dev/test/prod)
spec:
  podSelector:
    matchLabels:
      app.kubernetes.io/name: sso-patroni
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              environment: tools
              name: $LICENSE_PLATE
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: sso-grafana
  policyTypes:
    - Ingress
```

#### Update Helm Values

- Update data source username, password and database names under `values-$LICENSE_PLATE.yml` in place of `<please-replace-me>`

### Installing/Upgrading the Chart

```sh
make upgrade NAMESPACE=<namespace> \
             SSO_CLIENT_ID=<sso-client-id> \
             SSO_CLIENT_SECRET=<sso-client-secret>
```

- please find the SSO client credentials of the integration `#4492 SSO Dashboard` via [CSS app](https://bcgov.github.io/sso-requests):

### Uninstalling the Chart

```sh
make uninstall NAMESPACE=<namespace>
```

## Grafana Admin credentials

once the deployment is completed, please find the Grafana admin credentials in OCP secrets below:

- `dev`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/e4ca1d-tools/secrets/sso-grafana
- `prod`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/eb75ad-prod/secrets/sso-grafana

## Configuration

please find the full list of `Grafana Helm values` configuration in [Grafana Helm values](https://github.com/grafana/helm-charts/blob/main/charts/grafana/values.yaml)
