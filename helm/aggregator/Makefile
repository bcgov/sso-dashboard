SHELL := /usr/bin/env bash
NAMESPACE=
IMAGE_TAG=dev
NAME=sso-aggregator

ifndef NAMESPACE
$(error NAMESPACE is not set)
endif

.PHONY: helm-dep
helm-dep:
	helm dependency update

.PHONY: install
install: helm-dep
install:
	helm install "${NAME}" . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" --set image.tag="${IMAGE_TAG}"

.PHONY: upgrade
upgrade: helm-dep
upgrade:
	helm upgrade --install "${NAME}" . -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" --set image.tag="${IMAGE_TAG}"

.PHONY: lint
lint: helm-dep
lint:
	helm upgrade --dry-run --install "${NAME}" -n "${NAMESPACE}" -f values.yaml -f "values-${NAMESPACE}.yaml" --set image.tag="${IMAGE_TAG}"

.PHONY: uninstall
uninstall: helm-dep
uninstall:
	helm uninstall "${NAME}" -n ${NAMESPACE}

.PHONY: force-install
force-install: uninstall
force-install: install
