# sso-loki

A Helm chart for deploying [Grafana Loki](https://github.com/grafana/loki/tree/main/production/helm/loki).

## Local deployment via Helm chart

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

### Uninstalling the Chart

```sh
make uninstall NAMESPACE=<namespace>
```

## Configuration

please find the full list of `Loki Helm values` configuration in [Loki Helm values](https://github.com/grafana/loki/blob/main/production/helm/loki/values.yaml)

## Developer Notes

- It is expected to see the Loki pods failing until MinIO pods are ready to connect.
