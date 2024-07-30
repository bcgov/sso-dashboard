SHELL := /usr/bin/env bash
NAMESPACE=
MINIO_USER=
MINIO_PASS=
SSO_CLIENT_ID=
SSO_CLIENT_SECRET=
NAME=sso-loki

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

define arguments
	"${NAME}" . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" \
	--set root-minio.rootUser=${MINIO_USER} \
	--set root-minio.rootPassword=${MINIO_PASS} \
	--set root-minio.oidc.clientId=${SSO_CLIENT_ID} \
	--set root-minio.oidc.clientSecret=${SSO_CLIENT_SECRET} \
	--set loki.loki.storage.s3.accessKeyId=${MINIO_USER} \
	--set loki.loki.storage.s3.secretAccessKey=${MINIO_PASS}
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
