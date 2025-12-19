SHELL := /bin/sh

.PHONY: tidy deps swag run build

ROOT := /Users/idris.s/dev/fandom/notifications
SWAG := $(shell go env GOPATH)/bin/swag

tidy:
	cd $(ROOT) && go mod tidy

deps:
	cd $(ROOT) && go get github.com/gin-gonic/gin && go get github.com/swaggo/files && go get github.com/swaggo/gin-swagger && go get github.com/swaggo/swag

swag:
	@command -v $(SWAG) >/dev/null 2>&1 || (echo "Installing swag CLI..." && go install github.com/swaggo/swag/cmd/swag@latest)
	cd $(ROOT) && $(SWAG) init -g cmd/server/main.go -o docs

run:
	cd $(ROOT) && go run ./cmd/server

build:
	cd $(ROOT) && go build -o bin/server ./cmd/server

generate-key:
	cd $(ROOT) && go run ./cmd/generate-key -name "API Key"


