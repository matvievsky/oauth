# Makefile
GO_CMD=go
SOURCE_ENV=source .env
MAIN_GO=./cmd/tokenizer/main.go

.PHONY: get-token get-login-challenge

get-token:
	@$(SOURCE_ENV) && $(GO_CMD) run $(MAIN_GO) get

get-login-challenge:
	@$(SOURCE_ENV) && $(GO_CMD) run $(MAIN_GO) get-login-challenge

default: get-token
