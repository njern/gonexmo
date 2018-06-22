# MAINTAINER: David López <not4rent@gmail.com>

SHELL       =  /bin/bash
PKGS        =  $(shell go list ./... | grep -v /vendor/)
APP         =  gonexmo
REPORTS_DIR ?= .reports

# Overridable by CI

COMMIT_SHORT     ?= $(shell git rev-parse --verify --short HEAD)
VERSION          ?= $(COMMIT_SHORT)
VERSION_NOPREFIX ?= $(shell echo $(VERSION) | sed -e 's/^[[v]]*//')
REPORTS_DIR      ?= .reports

#
# Common methodology based targets
#

.PHONY: prepare
prepare: setup-env

.PHONY: sanity-check
sanity-check: goimports golint vet megacheck errcheck

.PHONY: test
test:
	@echo "Running unit tests..."
	mkdir -p $(REPORTS_DIR)
	2>&1 go test -cover -v $(shell echo $(PKGS) | tr " " "\n") | tee $(REPORTS_DIR)/report-unittests.out
	cat $(REPORTS_DIR)/report-unittests.out | go-junit-report -set-exit-code > $(REPORTS_DIR)/report-unittests.xml

#
# Custom golang project related targets
#

.PHONY: goimports
goimports:
	@echo "Running goimports..."
	@test -z "`for pkg in $(PKGS); do goimports -l $(GOPATH)/src/$$pkg/*.go | tee /dev/stderr; done`"

.PHONY: golint
golint:
	@echo "Running golint..."
	@golint -set_exit_status $(PKGS)

.PHONY: vet
vet:
	@echo "Running go vet..."
	@go vet $(PKGS)

.PHONY: megacheck
megacheck:
	@echo "Running megacheck..."
	@megacheck $(PKGS)

.PHONY: errcheck
errcheck:
	@echo "Running errcheck..."
	@errcheck $(PKGS)

.PHONY: setup
setup: setup-env

.PHONY: setup-env
setup-env:
	go get -u github.com/jstemmer/go-junit-report
	go get -u github.com/golang/lint/golint
	go get -u honnef.co/go/tools/cmd/megacheck
	go get -u github.com/kisielk/errcheck
	go get -u golang.org/x/tools/cmd/goimports

#
# Debug any makefile variable
# Usage: print-<VAR>
#
print-%  : ; @echo $* = $($*)