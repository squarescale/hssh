GIT_BRANCH   := $(shell git rev-parse --abbrev-ref HEAD || echo 'n/a')
GIT_REVISION :=	$(shell git describe --always --tags --dirty || echo 'n/a')

GO_BUILD_DATE  ?=$(shell date -u +%FT%T)
GO_VERSION_PKG ?= github.com/squarescale/hssh

GO_LD_FLAGS ?= -ldflags "-X main.GitCommit=$(GIT_REVISION) \
                         -X main.GitBranch=$(GIT_BRANCH)   \
                         -X main.BuildDate=$(GO_BUILD_DATE)"

all:
	export GOPROXY=https://proxy.golang.org && \
	export GO111MODULE=on && \
	go build $(GO_LD_FLAGS)
