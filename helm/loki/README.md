# Grafana Helm Chart

- Installs the log aggregation system [Loki](https://grafana.com/oss/loki)

## Installing/upgrading the Chart

To install/upgrade the chart with the release name `sso-loki`:

```console
make upgrade NAMESPACE=<namespace>
```

## Uninstalling the Chart

To uninstall/delete the my-release deployment:

```console
make uninstall NAMESPACE=<namespace>
```

- see https://github.com/grafana/loki/tree/main/production/helm/loki for more detail
