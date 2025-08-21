SHELL := /usr/bin/env bash

.PHONY: install
install:
	go mod tidy
	pip install -r requirements.txt
	sudo apt install reflex
	asdf reshim

.PHONY: build
build:
	go build -o build/aggregator ./cmd/aggregator/main.go
	go build -o build/compactor ./cmd/compactor/main.go

.PHONY: format
format:
	gofmt -w -s .

.PHONY: dev
dev: build
dev:
	./build/aggregator

.PHONY: aggregator-dev
aggregator-dev:
	reflex -r '\.go$$' -s -- sh -c 'go run ./cmd/aggregator'

.PHONY: compactor-dev
compactor-dev:
	reflex -r '\.go$$' -s -- sh -c 'go run ./cmd/comapctor'


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
