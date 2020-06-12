GIT_BRANCH   := $(shell git rev-parse --abbrev-ref HEAD || echo 'n/a')
GIT_REVISION :=	$(shell git describe --always --tags --dirty || echo 'n/a')

GO_BUILD_DATE  ?=$(shell date -u +%FT%T)
GO_VERSION_PKG ?= github.com/squarescale/hssh

GO_LD_FLAGS ?= -ldflags "-X $(GO_VERSION_PKG)/version.GitCommit=$(GIT_REVISION) \
                         -X $(GO_VERSION_PKG)/version.GitBranch=$(GIT_BRANCH)   \
                         -X $(GO_VERSION_PKG)/version.BuildDate=$(GO_BUILD_DATE)"

all:
	export GOPROXY=https://gocenter.io && \
	export GO111MODULE=on && \
	go build $(GO_LD_FLAGS)
