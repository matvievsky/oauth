#!/bin/bash
set -e
source .env
go run ./cmd/tokenizer/main.go get \
    --hydra-client-id $TOKENIZER_HYDRA_CLIENT_ID \
    --user-login $TOKENIZER_USER_LOGIN \
    --user-password $TOKENIZER_USER_PASSWORD
