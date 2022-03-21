.PHONY: all test

all: build

build: protogen
	go build ./...

# {{{ Protobuf

# Protobuf definitions
PROTO_FILES := $(shell find . \( -path "./languages" -o -path "./specification" \) -prune -o -type f -name '*.proto' -print)
# Protobuf Go files
PROTO_GEN_FILES = $(patsubst %.proto, %.pb.go, $(PROTO_FILES))

# Protobuf generator
PROTO_MAKER := protoc --gogoslick_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,$\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,$\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,$\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,$\
Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,$\
plugins=grpc:.


protogen: $(PROTO_GEN_FILES)

%.pb.go: %.proto
	cd $(dir $<); $(PROTO_MAKER) --proto_path=. --proto_path=$(GOPATH)/include ./*.proto
	sed -i '' -En -e '/^package [[:alpha:]]+/,$$p' $@

# }}} Protobuf end


# {{{ Cleanup
clean: protoclean

protoclean:
	rm -rf $(PROTO_GEN_FILES)
# }}} Cleanup end

# {{{ test

PROJECT_NAME := alpha-connect
PKG := gitlab.com/alphaticks/$(PROJECT_NAME)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v /models | grep -v /legacy)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)


test:
	go test $(PKG_LIST)

race:
	go test -race -short $(PKG_LIST)

coverage:
	coverage
	go tool cover -html=coverage.cov -o coverage.html

coverhtml:
	coverage
	go tool cover -html=coverage.cov -o coverage.html

lint: ## Lint the files
	echo ${PKG_LIST}
	go fmt ${PKG_LIST}
	go vet ${PKG_LIST}
	staticcheck ${PKG_LIST}

test-short:
	go test -short $(PKG_LIST)

# }}} test
