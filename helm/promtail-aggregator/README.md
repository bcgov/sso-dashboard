# sso-promtail (aggregator)

A Helm chart for deploying [Grafana Promtail](https://github.com/grafana/helm-charts/tree/main/charts/promtail) for our [aggregator server](../aggregator). This instance is separated from our loki promtail deployment in order to use additional labels for easier metric aggregation, which would negatively impact Loki performance.

## Local deployment via Helm chart

### Installing/Upgrading the Chart

```sh
make upgrade NAMESPACE=<namespace>
```

### Uninstalling the Chart

```sh
make uninstall NAMESPACE=<namespace>
```

## Configuration

please find the full list of `Promtail Helm values` configuration in [Promtail Helm values](https://github.com/grafana/helm-charts/blob/main/charts/promtail/values.yaml)
