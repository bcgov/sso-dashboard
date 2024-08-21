# sso-promtail

A Helm chart for deploying [Grafana Promtail](https://github.com/grafana/helm-charts/tree/main/charts/promtail) for our loki instance.

## Labelling

For grafana loki we have found that keeping a minimal label set is ideal for performance, more information can be found in the following resources:

- [Cardinality](https://grafana.com/docs/loki/latest/get-started/labels/#cardinality)
- [Labels and Chunks](https://grafana.com/blog/2023/12/20/the-concise-guide-to-grafana-loki-everything-you-need-to-know-about-labels/)
- [Guide to label configuration](https://grafana.com/blog/2020/08/27/the-concise-guide-to-labels-in-loki/)

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
