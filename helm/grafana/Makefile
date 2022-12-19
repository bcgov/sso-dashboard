SHELL := /usr/bin/env bash
NAMESPACE=""

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

.PHONY: helm-dep
helm-dep:
	helm dependency update

.PHONY: install
install: helm-dep
install:
	helm install sso-grafana . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml"

.PHONY: upgrade
upgrade: helm-dep
upgrade:
	helm upgrade --install sso-grafana . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml"

.PHONY: lint
lint: helm-dep
lint:
	helm upgrade --dry-run --install sso-grafana -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml"

.PHONY: uninstall
uninstall: helm-dep
uninstall:
	helm uninstall sso-grafana -n ${NAMESPACE}

.PHONY: force-install
force-install: uninstall
force-install: install
