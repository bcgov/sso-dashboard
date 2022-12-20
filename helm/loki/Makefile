SHELL := /usr/bin/env bash
NAMESPACE=
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY=
NAME=sso-loki

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

define arguments
	"${NAME}" . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" \
	--set loki.storage.s3.access_key_id=${S3_ACCESS_KEY_ID} \
	--set loki.storage.s3.secret_access_key=${S3_SECRET_ACCESS_KEY}
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
