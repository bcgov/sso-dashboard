# sso-aggregator

A Helm chart for deploying [Keycloak event log aggregator](../../aggregator).

## Setup

This chart includes two different deployments, running the same dockerfile. The dockerfile build copies the compactor into app/compactor, and the aggregator into app/aggregator. In the compactor [deployment template](./templates/deployment-compactor.yaml) the command is used to run the compactor, whereas the [entry shell file](../../aggregator/docker-entrypoint.sh) will run the aggregator by default (last line).

## Local deployment via Helm chart

### Pre-Requisites

- Ensure below network policy exists in the namespace where loki is being deployed

```yaml
kind: NetworkPolicy
apiVersion: networking.k8s.io/v1
metadata:
  name: allow-sso-promtail
  namespace: xxxx-xxxx
spec:
  podSelector: {}
  ingress:
    - from:
        - namespaceSelector:
            matchLabels:
              name: xxxx
        - podSelector:
            matchLabels:
              app.kubernetes.io/name: promtail
  policyTypes:
    - Ingress
status: {}
```

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

- `sandbox`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/e4ca1d-prod/secrets/sso-aggregator-patroni-appusers
- `prod`: https://console.apps.gold.devops.gov.bc.ca/k8s/ns/eb75ad-prod/secrets/sso-aggregator-patroni-appusers
