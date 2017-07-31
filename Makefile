GO_PATH=$(PWD):$(PWD)/vendor:$(GOPATH)

all: chat proxy

chat:
	GOPATH="$(GO_PATH)" go build -o ./bin/chat ./src/chat 

proxy:
	GOPATH="$(GO_PATH)" go build -o ./bin/proxy ./src/proxy 


