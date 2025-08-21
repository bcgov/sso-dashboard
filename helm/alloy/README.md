# Overview

This repository contains configuration for [grafana alloy](https://grafana.com/docs/alloy/latest/), which is used to send our keycloak logs to the sso-aggregator and loki instance.

## Configuration

The following environment variables are required:
- API_GATEWAY_URL: Base url to the loki instance.
- LOKI_AUTH_TOKEN: Bearer token for loki.

They can be saved in a `.env` file if deploying with the local Makefile.

## Installing

After setting environment variables and logging into the correct env, set the NAMESPACE argument in the makefile and run `make install`. If adjusting values, the app will deploy when merged to dev/main.

## Legacy Migration

This alloy instance was migrated from existing promtail instances after they were deprecated. To keep the log positions, the promtail PVC is mounted on this instance, named `sso-promtail-aggregator-positions`. See the [legacy positions file](https://grafana.com/docs/alloy/latest/reference/components/loki/loki.source.file/#arguments) argument for more details.
