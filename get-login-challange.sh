#!/bin/bash
set -e
source .env
go run ./cmd/tokenizer/main.go get-login-challenge