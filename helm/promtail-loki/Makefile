-include /$(PWD)/.env

SHELL := /usr/bin/env bash
NAMESPACE=""
NAME=sso-promtail-aws-loki

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

define arguments
	"${NAME}" . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" \
	--set promtail.config.clients[0].headers.Authorization="Bearer ${LOKI_AUTH_TOKEN}"
	--set promtail.config.clients[0].url="${API_GATEWAY_URL}/loki/api/v1/push"
endef

.PHONY: helm-dep
helm-dep:
	helm dependency update

.PHONY: install
install: helm-dep
install:
	@helm install $(call arguments)

.PHONY: upgrade
upgrade: helm-dep
upgrade:
	echo $NAMESPACE
	echo $API_GATEWAY_URL
	@helm upgrade --install $(call arguments)

.PHONY: lint
lint: helm-dep
lint:
	@helm upgrade --dry-run --install $(call arguments)

.PHONY: uninstall
uninstall: helm-dep
uninstall:
	@helm uninstall ${NAME} -n ${NAMESPACE}

.PHONY: template
template: helm-dep
template:
	@helm template $(call arguments) > template.yaml

.PHONY: force-install
force-install: uninstall
force-install: install
