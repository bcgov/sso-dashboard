# sso-loki

A Helm chart for deploying [Grafana Loki](https://github.com/grafana/loki/tree/main/production/helm/loki).

## Local deployment via Helm chart

### Pre-Requisites

- It's generally a good practice to stop Alloy before restarting Grafana Loki. Alloy is responsible for scraping and sending logs to Loki, and stopping it before a restart can prevent potential issues or data loss during the restart process. Stopping Alloy temporarily ensures that it doesn’t try to send logs while Loki is restarting, preventing any potential errors due to a disrupted connection. After Loki has restarted successfully, you can start Alloy again to resume log scraping and forwarding to Loki. This sequence helps in maintaining the integrity of log data and ensures a smoother restart process for Loki.

```sh
export LICENSE_PLATE=

oc scale --replicas=0 deployment alloy-keycloak -n ${LICENSE_PLATE}-dev
oc scale --replicas=0 deployment alloy-keycloak -n ${LICENSE_PLATE}-test
oc scale --replicas=0 deployment alloy-keycloak -n ${LICENSE_PLATE}-prod
```

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
make upgrade NAMESPACE=<namespace> \
             SSO_CLIENT_ID=<sso-client-id> \
             SSO_CLIENT_SECRET=<sso-client-secret> \
             MINIO_USER=<sso-client-id> \
             MINIO_PASS=<sso-client-secret>
```

- please find the SSO client credentials of the integration `#4492 SSO Dashboard` via [CSS app](https://bcgov.github.io/sso-requests):

- please generate the secure credentials for the initial `MinIO Admin` that can be set in the MinIO deployment.

### Post Installation/Update of Loki

```sh
oc scale --replicas=1 deployment alloy-keycloak -n ${LICENSE_PLATE}-dev
oc scale --replicas=1 deployment alloy-keycloak -n ${LICENSE_PLATE}-test
oc scale --replicas=1 deployment alloy-keycloak -n ${LICENSE_PLATE}-prod
```

### Uninstalling the Chart

```sh
make uninstall NAMESPACE=<namespace>
```

## Configuration

please find the full list of `Loki Helm values` configuration in [Loki Helm values](https://github.com/grafana/loki/blob/main/production/helm/loki/values.yaml)

## Developer Notes

- It is expected to see the Loki pods failing until MinIO pods are ready to connect.
