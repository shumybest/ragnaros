.PHONY: default build docker test clean bench

BINARY = {{ .App.ProjectName }}
Image = {{ .App.ProjectName }}

PWD            := $(shell pwd)
PKG            := {{ .App.ProjectName }}
TRAVIS_COMMIT ?= `git rev-parse HEAD`
GOCMD          = go
BUILD_TIME     = `date +'%Y-%m-%d %H:%M:%S'`
VERSION        = 1.0.0
GOFLAGS       ?= $(GOFLAGS:)
LDFLAGS       := "-w -X '$(PKG)/config.GitCommit=$(TRAVIS_COMMIT)' \
                  -X '$(PKG)/config.Version=$(VERSION)' \
		          -X '$(PKG)/config.BuildTime=$(BUILD_TIME)'"
GOARCH = amd64
CGO_ENABLED = 0
TAGS_OPT = -tags linux
GOOS = linux

ifeq ($(TARGET), local)
  TAGS_OPT = -tags darwin
  GOOS = darwin
endif

ifeq ($(CGO), 1)
  CGO_ENABLED = 1
endif

ifeq ($(ARCH), arm)
  GOARCH = arm64
  ifeq ($(CGO), 1)
    CC = aarch64-linux-gnu-gcc
    CXX = aarch64-linux-gnu-g++
  endif
endif

default: build test

build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOCMD) build ${GOFLAGS} -ldflags ${LDFLAGS} ${TAGS_OPT} -o ${BINARY}

docker: build
	@docker build -t ${Image} --build-arg BINARY=${BINARY} -f Dockerfile .

test:
	$(GOCMD) test -race -v $(shell go list ./... | grep -v '/vendor/')

bench:
	$(GOCMD) test -race -run none -v --bench=. ./...

coverage:
	go test -run none -v -cover ./...

clean:
	rm ${BINARY}