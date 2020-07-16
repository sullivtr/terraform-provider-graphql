GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
TEST_DESTS := $(dir $(wildcard ./test/*/*test.tf))

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

fetch: ## download makefile dependencies
	@hash goreleaser 2>/dev/null || go get -u -v github.com/goreleaser/goreleaser

clean: ## cleans previously built binaries
	rm -rf ./dist

build: clean fetch ## publishes in dry run mode
	$(GOPATH)/bin/goreleaser --skip-publish --snapshot


.PHONY: test copyplugins

copyplugins: ## copy plugins to test folders
	$(eval COPY_FILES := $(wildcard ./dist/terraform-provider-graphql*/*))
	$(eval OS_ARCH := $(patsubst ./dist/terraform-provider-graphql_%/terraform-provider-graphql, %, $(COPY_FILES)))
	$(eval TEST_FOLDERS := $(foreach p,$(OS_ARCH), $(patsubst %,%terraform.d/plugins/$p,$(TEST_DESTS))))
	@sleep 1
	@mkdir -p $(TEST_FOLDERS)
	@for o in $(OS_ARCH); do \
		for f in $(TEST_DESTS); do \
			cp ./dist/terraform-provider-graphql_$$o/* $$f/terraform.d/plugins/$$o; \
		done; \
	done

test: copyplugins ## test
	@cd test && $(MAKE) test

fulltest: build test ## build and test