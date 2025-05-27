-include /$(PWD)/.env

SHELL := /usr/bin/env bash
NAMESPACE=""
NAME=alloy

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

define arguments
	"${NAME}" -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" grafana/alloy --version 1.0.2 \
	--set extraConfiguration.loki.bearerToken="${LOKI_AUTH_TOKEN}" \
	--set extraConfiguration.loki.url="${API_GATEWAY_URL}/loki/api/v1/push"
endef


.PHONY: install
install:
	@helm install $(call arguments)

.PHONY: upgrade
upgrade:
	@helm upgrade --install $(call arguments)

.PHONY: uninstall
uninstall:
	@helm uninstall ${NAME} -n ${NAMESPACE}
