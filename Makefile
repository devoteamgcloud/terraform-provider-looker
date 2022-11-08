ifeq ($(OS), Windows_NT)
	export VERSION=$(shell type VERSION)
else
	export VERSION=$(shell cat VERSION)
endif
export BASE_BINARY_NAME=terraform-provider-looker
export ORG=devoteamgcloud

.PHONY: build
build: ## build binary
	@go build -o build/$(ORG)/$(VERSION)/$(BASE_BINARY_NAME) .

.PHONY: format
format: ## format all the go files
	@gofmt -l -s -w .

.PHONY: docs
docs: ## generate documentation
	@go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
