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
	helm install sso-loki . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml"

.PHONY: upgrade
upgrade: helm-dep
upgrade:
	helm upgrade --install sso-loki . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml"

.PHONY: lint
lint: helm-dep
lint:
	helm upgrade --dry-run --install sso-loki -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml"

.PHONY: uninstall
uninstall: helm-dep
uninstall:
	helm uninstall sso-loki -n ${NAMESPACE}

.PHONY: force-install
force-install: uninstall
force-install: install
