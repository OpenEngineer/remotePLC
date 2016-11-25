#!/bin/bash
export CGO_ENABLED=0

go build -ldflags "-linkmode external -extldflags -static"
