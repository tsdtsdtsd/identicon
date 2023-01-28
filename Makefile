GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test
COVERFILE=coverage.out

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: test test/cover

all: help

## Test:
test: ## Run all tests
	$(GOTEST) -v -race ./...
	
cover: ## Run tests and open coverage in browser
	$(GOTEST) -v -coverpkg=./... -covermode=atomic -coverprofile=$(COVERFILE) ./...
	$(GOCOVER) -func=$(COVERFILE)
	$(GOCOVER) -html=$(COVERFILE)
	rm $(COVERFILE)

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)