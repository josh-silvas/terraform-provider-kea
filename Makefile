# -------------------------------------------------------------------------------------------
# VARIABLES: Variable declarations to be used within make to generate commands.
# -------------------------------------------------------------------------------------------
PROJECT_NAME := terraform-provider-kea
VERSION      := $(shell cat VERSION)
COMPOSE      := docker-compose --project-name $(PROJECT_NAME) --project-directory "develop" -f "develop/docker-compose.yml"

default: help

# -------------------------------------------------------------------------------------------
# RELEASE: Release management directives.
# -------------------------------------------------------------------------------------------
release: ## Run goreleaser to create a release
	@rm -rf dist/*
	@git tag -d $(VERSION) || true
	@git tag $(VERSION)
	@goreleaser --rm-dist --skip-validate --skip-announce
.PHONY: release

tag: ## Tag and push the version defined in VERSION file
	@git tag -d $(VERSION) || true
	@git tag $(VERSION)
	@git push origin $(VERSION)
.PHONY: tag

# -------------------------------------------------------------------------------------------
# CODE-QUALITY/TESTS: Linting and testing directives.
# -------------------------------------------------------------------------------------------
_lint: ## Run golangci-lint on all sub-packages
	@echo "🧪 Running golangci-lint..."
	@golangci-lint run --tests=false --exclude-use-default=false
	@echo "Completed golangci-lint."
.PHONY: _lint

_testacc: ## Run acceptance tests
	@echo "🛠 Running acceptance tests..."
	@TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m  | { grep -v 'no test files'; true; }
	@echo "Completed acceptance tests."
.PHONY: _testacc

# -------------------------------------------------------------------------------------------
# DEVELOPMENT: Development tools for use when contributing to this project.
# -------------------------------------------------------------------------------------------
develop: .env ## Build the development docker image and push to registry
	@echo "🐳 Building development docker image and pushing to registry..."
	@$(COMPOSE) build --no-cache
	@docker push jsilvas/${PROJECT_NAME}-develop:latest
.PHONY: develop

cli: .env ## Launch a bash shell inside the running container.
	@echo "🐳 Launching a bash shell 💻 inside the running container..."
	@$(COMPOSE) run --rm develop bash
.PHONY: cli

destroy: .env ## Destroy the docker-compose environment and volumes
	@$(COMPOSE) down --volumes
.PHONY: destroy

lint: .env ## Run golangci-lint on all sub-packages within docker
	@echo "🐳 Launching golangci-lint in docker..."
	@$(COMPOSE) run --rm develop make _lint
.PHONY: lint

testacc: .env ## Run acceptance tests on all sub-packages within docker
	@echo "🐳 Launching acceptance tests in docker..."
	@$(COMPOSE) run --rm develop make _testacc
.PHONY: testacc

# -------------------------------------------------------------------------------------------
# BUILD: Build the provider
# -------------------------------------------------------------------------------------------
build: ## Build the provider for local development
	@echo "Building the provider..."
	@go build -o develop/terraform-provider-kea
	@echo "Completed building the provider."
.PHONY: build

terraform-plan: ## Run terraform plan to test the provider
	@cd develop && terraform plan
.PHONY: terraform-plan

terraform-apply: ## Run terraform apply to test the provider
	@cd develop && terraform apply
.PHONY: terraform-apply

# -------------------------------------------------------------------------------------------
# DOCUMENTATION: Generate documentation
# -------------------------------------------------------------------------------------------
docs: ## Run go generate to create documentation in the docs subfolder
	@go generate ./...
	@git add docs/*
.PHONY: docs

# -------------------------------------------------------------------------------------------
# HELPERS: Internal Make Commands
# -------------------------------------------------------------------------------------------
tidy: ## Run go mod tidy and go mod vendor
	@go mod tidy && go mod vendor
.PHONY: tidy

.env:
	@if [ ! -f "develop/.env" ]; then \
	   echo "Creating environment file...\nPLEASE OVERRIDE VARIABLES IN develop/.env WITH YOUR OWN VALUES!"; \
	   cp develop/example.env develop/.env; \
	fi
.PHONY: .env

help: ## Display this help screen
	@echo "\033[1m\033[01;32m\
	$(shell echo $(PROJECT_NAME) | tr  '[:lower:]' '[:upper:]') $(VERSION) \
	\033[00m\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' \
	$(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; \
	{printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
.PHONY: help