SHELL := /usr/bin/env bash
NAMESPACE=
SSO_CLIENT_ID=
SSO_CLIENT_SECRET=
AGGREGATION_DB_USER=
AGGREGATION_DB_PASS=
NAME=sso-grafana

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

define arguments
	"${NAME}" . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" \
	--set grafana."grafana\.ini"."auth\.generic_oauth".client_id=${SSO_CLIENT_ID} \
	--set grafana."grafana\.ini"."auth\.generic_oauth".client_secret=${SSO_CLIENT_SECRET} \
	--set grafana.datasources."datasources\.yaml".datasources[1].user=${AGGREGATION_DB_USER} \
	--set grafana.datasources."datasources\.yaml".datasources[1].secureJsonData.password=${AGGREGATION_DB_PASS}
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
