-include /$(PWD)/.env

SHELL := /usr/bin/env bash
NAMESPACE="e4ca1d-prod"
NAME=alloy

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

define arguments
	"${NAME}" -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" grafana/alloy --version 1.8.2
endef

.PHONY: deps
deps:
	@helm repo add grafana https://grafana.github.io/helm-charts

.PHONY: install
install:
	@helm install $(call arguments)

.PHONY: upgrade
upgrade:
	@helm upgrade --install $(call arguments)

.PHONY: uninstall
uninstall:
	@helm uninstall ${NAME} -n ${NAMESPACE}
