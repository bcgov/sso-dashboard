SHELL := /usr/bin/env bash

.PHONY: install
install:
	go mod tidy
	pip install -r requirements.txt
	asdf reshim

.PHONY: build
build:
	go build -o build/aggregator main.go
	go build -o build/compactor compactor.go

.PHONY: format
format:
	gofmt -w -s .

.PHONY: dev
dev: build
dev:
	./build/aggregator

.PHONY: compactor
compactor: build
compactor:
	./build/compactor

define db-start
	pg_ctl start
endef

define db-setup
	bash ./setup.sh
	alembic upgrade head
endef

.PHONY: db
db:
	$(call db-start) || true
	$(call db-setup)
