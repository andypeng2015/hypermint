GO_BIN?=go

BUILD_DIR?=./build
HMD?=$(BUILD_DIR)/hmd
HMCLI?=$(BUILD_DIR)/hmcli
HMD_HOME?=${HOME}/.hmd
HMCLI_HOME?=${HOME}/.hmcli

MNEMONIC?="token dash time stand brisk fatal health honey frozen brown flight kitchen"
HDW_PATH?=m/44'/60'/0'/0

COMMIT_HASH:=$(shell git rev-parse --short HEAD)
VERSION:=$(shell cat version)
BUILD_FLAGS?=-ldflags "-X github.com/bluele/hypermint/pkg/consts.GitCommit=${COMMIT_HASH} -X github.com/bluele/hypermint/pkg/consts.Version=${VERSION}"

GO_BUILD_CMD=$(GO_BIN) build $(BUILD_FLAGS)
GO_TEST_FLAGS?=-v
GO_TEST_CMD=$(GO_BIN) test $(GO_TEST_FLAGS)

.PHONY: build

build: server cli

server:
	$(GO_BUILD_CMD) -o $(HMD) ./cmd/hmd

cli:
	$(GO_BUILD_CMD) -o $(HMCLI) ./cmd/hmcli

start:
	$(HMD) start --log_level="main:error" --home=$(HMD_HOME)

clean:
	@rm -rf $(HMD_HOME) $(HMCLI_HOME)

init: clean init-validator
	$(eval ADDR1 := $(shell $(HMCLI) new --password=password --silent --home=$(HMCLI_HOME) --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/1" ))
	$(eval ADDR2 := $(shell $(HMCLI) new --password=password --silent --home=$(HMCLI_HOME) --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/2" ))
	@$(HMD) init --address=$(ADDR1) --home=$(HMD_HOME)
	@echo export ADDR1='$(ADDR1)'
	@echo export ADDR2='$(ADDR2)'

init-validator:
	@$(HMD) tendermint init-validator --mnemonic=$(MNEMONIC) --hdw_path="$(HDW_PATH)/0"

test:
	$(GO_TEST_CMD) ./pkg/...

build-image:
	docker build . -t bluele/hypermint:${VERSION}