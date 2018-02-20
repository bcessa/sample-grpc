.PHONY: proto
default: build
LD_FLAGS=""
BINARY_NAME=sample-grpc
KNOWN_TYPES="github.com/gogo/protobuf/types"
PROTO_TYPES="Mgoogle/protobuf/any.proto=$(KNOWN_TYPES),Mgoogle/protobuf/empty.proto=$(KNOWN_TYPES),Mgoogle/protobuf/duration.proto=$(KNOWN_TYPES),Mgoogle/protobuf/field_mask.proto=$(KNOWN_TYPES),Mgoogle/protobuf/struct.proto=$(KNOWN_TYPES),Mgoogle/protobuf/timestamp.proto=$(KNOWN_TYPES),Mgoogle/protobuf/wrappers.proto=$(KNOWN_TYPES)"
PROTO_OUT = "Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto,$(PROTO_TYPES)"

proto: ## Compile protocol files
	@protoc \
	--proto_path=./vendor/github.com/gogo/protobuf \
	--proto_path=./vendor/github.com/gogo/protobuf/protobuf \
	--proto_path=./vendor/github.com/grpc-ecosystem/grpc-gateway \
	--proto_path=./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
	--proto_path=./vendor \
	--proto_path=. \
	--gogofaster_out=$(PROTO_OUT),plugins=grpc:. \
	--grpc-gateway_out=$(PROTO_TYPES),logtostderr=true:. \
	--swagger_out=logtostderr=true:. \
	--descriptor_set_out=proto/service.desc \
	proto/*.proto

build: ## Build for local system
	@go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)

linux: ## Build for AMD64 linux systems
	@env GOOS=linux GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-linux-amd64

install: ## Install to local machine
	@go build -v -ldflags $(LD_FLAGS) -i -o ${GOPATH}/bin/$(BINARY_NAME)

deps-list: ## List installed vendor dependencies
	dep status

dependencies: ## Install required vendor dependencies
	dep ensure -v

help: ## Display available make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-16s\033[0m %s\n", $$1, $$2}'