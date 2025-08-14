NAME      := postfix-log-filter
BUILDNAME := plf
VERSION   := v0.1.0
REVISION  := $(shell git rev-parse --short HEAD)

SRCS    := $(shell find . -type f -name '*.go')
DSTDIR  := /usr/local/bin
USER    := root
GROUP   := root
LDFLAGS := -ldflags="-s -w -X \"main.PackageName=$(NAME)\" -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\" -extldflags \"-static\""

GOVER     := $(shell go version | awk '{ print substr($$3, 3) }' | tr "." " ")
VER_JUDGE := $(shell if [ $(word 1,$(GOVER)) -eq 1 ] && [ $(word 2,$(GOVER)) -le 10 ]; then echo 0; else echo 1; fi)

.PHONY: run
run:
	@go run *.go

.PHONY: init
init:
ifeq ($(VER_JUDGE),1)
	@go mod init $(NAME) && go mod tidy
else
	@echo "Packageの取得は手動で行ってください"
endif

.PHONY: modup
modup:
	@go mod tidy

.PHONY: build
build: $(SRCS)
	@go build -a -tags netgo -installsuffix netgo $(LDFLAGS) -o bin/$(BUILDNAME)

.PHONY: install
install:
	@command cp -r bin/$(BUILDNAME) $(DSTDIR)/
	@chown $(USER):$(GROUP) $(DSTDIR)/$(BUILDNAME)

.PHONY: uninstall
	@rm -f $(DSTDIR)/$(BUILDNAME)

.PHONY: clean
clean:
	@rm -rf bin/*
	@rm -rf vendor/*

.PHONY: test
test:
	@go test
