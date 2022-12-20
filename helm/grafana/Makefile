SHELL := /usr/bin/env bash
NAMESPACE=
SSO_CLIENT_ID=
SSO_CLIENT_SECRET=
SSO_AUTH_URL=
SSO_TOKEN_URL=
NAME=sso-grafana

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

define arguments
	"${NAME}" . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" \
	--set "grafana\.ini"."auth\.generic_oauth".client_id=${SSO_CLIENT_ID} \
	--set "grafana\.ini"."auth\.generic_oauth".client_secret=${SSO_CLIENT_SECRET} \
	--set "grafana\.ini"."auth\.generic_oauth".auth_url=${SSO_AUTH_URL} \
	--set "grafana\.ini"."auth\.generic_oauth".token_url=${SSO_TOKEN_URL}
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

.PHONY: force-install
force-install: uninstall
force-install: install
