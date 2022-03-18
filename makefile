SHELL=/bin/bash

BLACK        := $(shell tput -Txterm setaf 0)
RED          := $(shell tput -Txterm setaf 1)
GREEN        := $(shell tput -Txterm setaf 2)
YELLOW       := $(shell tput -Txterm setaf 3)
LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
PURPLE       := $(shell tput -Txterm setaf 5)
BLUE         := $(shell tput -Txterm setaf 6)
WHITE        := $(shell tput -Txterm setaf 7)
RESET := $(shell tput -Txterm sgr0)

SHA:=$(shell git describe --match 'v[0-9]*' --dirty=".d" --always)

.DEFAULT_GOAL := help

.PHONY: help
help:
	@echo "${LIGHTPURPLE}Current Commit SHA: $(SHA)${RESET}"
	@echo "${YELLOW}Targets:${RESET}"
	@echo " - ${BLUE}help:${WHITE}  shows this help message${RESET}"
	@echo " - ${BLUE}build:${WHITE} builds all go binaries into the /bin directory${RESET}"
	@echo " - ${BLUE}clean:${WHITE} removes all go binaries from the /bin directory${RESET}"

## Install dependencies
.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build: clean-tools tidy test build-tools

.PHONY: test
test:
	./runtests

.PHONY: build-tools
build-tools:
	env GOOS=darwin go build -ldflags="-X main.CommitID=${SHA} -s -w" -o "bin/tools/build_dictionary" tools/build_dictionary.go

.PHONY: clean-tools
clean-tools:
	rm -rf bin/tools

.PHONY: clean
clean:
	rm -rf bin
