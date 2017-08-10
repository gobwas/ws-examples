GO_PATH=$(PWD):$(PWD)/vendor:$(GOPATH)

.PHONY: vendor

all: vendor chat proxy

chat:
	GOPATH="$(GO_PATH)" go build -o ./bin/chat ./src/chat 

proxy:
	GOPATH="$(GO_PATH)" go build -o ./bin/proxy ./src/proxy 

vendor:
	if [ ! -d ./vendor ]; then \
		git submodule init; \
		git submodule update; \
	fi; \
	echo "OK"


