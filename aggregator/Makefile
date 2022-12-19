SHELL := /usr/bin/env bash

.PHONY: install
install:
	go mod tidy
	pip install -r requirements.txt
	asdf reshim

.PHONY: build
build:
	go build -o build/main

.PHONY: format
format:
	gofmt -w -s .

.PHONY: start
start:
	./build/main

.PHONY: dev
dev: build
dev: start

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
