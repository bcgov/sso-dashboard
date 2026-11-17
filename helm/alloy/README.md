# Overview

This directory contains configuration for [Grafana Alloy](https://grafana.com/docs/alloy/latest/). Alloy collects and transforms Keycloak logs, then pushes them to the SSO aggregator.

## Configuration

Set the environment and aggregator URL in the Helm values.

They can be saved in a `.env` file if deploying with the local Makefile.

## Installing

After setting environment variables and logging into the correct env, set the NAMESPACE argument in the makefile and run `make install`. If adjusting values, the app will deploy when merged to dev/main.

## Legacy Migration

This Alloy instance was migrated from existing Promtail instances after they were deprecated. To keep log positions, the Promtail PVC is mounted on this instance, named `sso-promtail-aggregator-positions`. See the [legacy positions file](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.source.file/#arguments) argument for more details.
