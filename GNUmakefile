# Pins version of tools to match that of what we use in CI
# assumes you are using tfenv to manage your terraform versions
# and have a go version >= 1.21 where toolchains were introduced. 
export GOTOOLCHAIN := go1.17
export TFENV_TERRAFORM_VERSION := 1.0.5
GORELEASER_VERSION := v1.5.0

GOPATH := $(shell go env GOPATH)
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
TEST_DESTS := $(dir $(wildcard ./e2e/*/*test.tf))

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fetch: ## download makefile dependencies
	go install github.com/goreleaser/goreleaser@$(GORELEASER_VERSION)

clean: ## cleans previously built binaries and test folders
	@for f in $(TEST_DESTS); do \
	  rm -rf $$f/.terraform; \
	  rm -rf $$f/terraform.d; \
	done
	@rm -rf ./dist;

build: clean fetch ## publishes in dry run mode
	$(GOPATH)/bin/goreleaser --skip-publish --snapshot --skip-sign 


.PHONY: unittest e2etest copyplugins

copyplugins: ## copy plugins to test folders
	$(eval COPY_FILES := $(filter %/, $(wildcard ./dist/terraform-provider-graphql*/)))
	$(eval OS_ARCH := $(patsubst ./dist/terraform-provider-graphql_%/, %, $(COPY_FILES)))
	$(eval TEST_FOLDERS := $(foreach p,$(OS_ARCH), $(patsubst %,%terraform.d/plugins/gqlprovidertf.com/examplecorp/graphql/2.0.0/$p,$(TEST_DESTS))))
	@sleep 1
	@mkdir -p $(TEST_FOLDERS);
	@for o in $(OS_ARCH); do \
		for f in $(TEST_DESTS); do \
	    	cp ./dist/terraform-provider-graphql_$$o/* $$f/terraform.d/plugins/gqlprovidertf.com/examplecorp/graphql/2.0.0/$$o; \
		done; \
	done

unittest: ## run the graphql package unit & acceptance tests
	go test -v -coverprofile=cover.out ./... | awk '!/\[no test files\]/ {print}'

e2etest: ## run the E2E tests (requires a build and copied plugins)
	@cd e2e && $(MAKE) test

fulltest: build copyplugins unittest e2etest ## build, copy plugins, then run unit and E2E tests