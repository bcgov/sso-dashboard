SHELL := /usr/bin/env bash
NAMESPACE=""
NAME=sso-promtail

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

define arguments
	"${NAME}" . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml"
endef

.PHONY: helm-dep
helm-dep:
	helm dependency update

.PHONY: install
install: helm-dep
install:
	helm install $(call arguments)

.PHONY: upgrade
upgrade: helm-dep
upgrade:
	helm upgrade --install $(call arguments)

.PHONY: lint
lint: helm-dep
lint:
	helm upgrade --dry-run --install $(call arguments)

.PHONY: uninstall
uninstall: helm-dep
uninstall:
	helm uninstall ${NAME} -n ${NAMESPACE}

.PHONY: template
template: helm-dep
template:
	helm template $(call arguments) > template.yaml

.PHONY: force-install
force-install: uninstall
force-install: install
