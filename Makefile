VERSION := $(shell cat ./VERSION)

all: build

build: vendor
	CGO_ENABLED=0 go build -v

test: vendor
	go test -v ./...

fmt:
	go fmt -x ./...

release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)

vendor: bootstrap
	dep ensure


HAS_DEP := $(shell command -v dep;)
HAS_LINT := $(shell command -v golint;)

bootstrap:
ifndef HAS_DEP
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: build test fmt vendor release
