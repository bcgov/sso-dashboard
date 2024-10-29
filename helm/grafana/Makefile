include /$(PWD)/.env

SHELL := /usr/bin/env bash

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

define arguments
	"${NAME}" . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" \
	--set grafana."grafana\.ini"."auth\.generic_oauth".client_id=${SSO_CLIENT_ID} \
	--set grafana."grafana\.ini"."auth\.generic_oauth".client_secret=${SSO_CLIENT_SECRET} \
	--set grafana.datasources."datasources\.yaml".datasources[0].user=${DS_SSO_LOGS_USERNAME} \
	--set grafana.datasources."datasources\.yaml".datasources[0].database=${DS_SSO_LOGS_DATABASE} \
	--set grafana.datasources."datasources\.yaml".datasources[0].secureJsonData.password=${DS_SSO_LOGS_PASSWORD} \
	--set grafana.datasources."datasources\.yaml".datasources[1].user=${DS_KEYCLOAK_DEV_USERNAME} \
	--set grafana.datasources."datasources\.yaml".datasources[1].database=${DS_KEYCLOAK_DEV_DATABASE} \
	--set grafana.datasources."datasources\.yaml".datasources[1].secureJsonData.password=${DS_KEYCLOAK_DEV_PASSWORD} \
	--set grafana.datasources."datasources\.yaml".datasources[2].user=${DS_KEYCLOAK_TEST_USERNAME} \
	--set grafana.datasources."datasources\.yaml".datasources[2].database=${DS_KEYCLOAK_TEST_DATABASE} \
	--set grafana.datasources."datasources\.yaml".datasources[2].secureJsonData.password=${DS_KEYCLOAK_TEST_PASSWORD} \
	--set grafana.datasources."datasources\.yaml".datasources[3].user=${DS_KEYCLOAK_PROD_USERNAME} \
	--set grafana.datasources."datasources\.yaml".datasources[3].database=${DS_KEYCLOAK_PROD_DATABASE} \
	--set grafana.datasources."datasources\.yaml".datasources[3].secureJsonData.password=${DS_KEYCLOAK_PROD_PASSWORD} \
	--set grafana.datasources."datasources\.yaml".datasources[4].secureJsonData.httpHeaderValue1="Bearer ${LOKI_AUTH_TOKEN}" \
	--set grafana.datasources."datasources\.yaml".datasources[4].url=${API_GATEWAY_URL} \
	--set grafana.datasources."datasources\.yaml".datasources[5].user=${DS_AGGREGATOR_USERNAME} \
	--set grafana.datasources."datasources\.yaml".datasources[5].database=${DS_AGGREGATOR_DATABASE} \
	--set grafana.datasources."datasources\.yaml".datasources[5].secureJsonData.password=${DS_AGGREGATOR_PASSWORD}
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
