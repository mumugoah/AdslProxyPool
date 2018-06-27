#!/usr/bin/env bash

env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/adslProxyClient-linux