GO_PATH=$(PWD):$(PWD)/vendor:$(GOPATH)

.PHONY: vendor cli

all: vendor chat proxy

cli:
	GOPATH="$(GO_PATH)" go build -o ./bin/cli ./src/chat/client

chat:
	GOPATH="$(GO_PATH)" go build -o ./bin/chat ./src/chat 

proxy:
	GOPATH="$(GO_PATH)" go build -o ./bin/proxy ./src/proxy 

vendor:
	git submodule init; \
	git submodule update; \

