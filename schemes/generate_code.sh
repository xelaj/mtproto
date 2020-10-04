#!/bin/sh

go run schemes/build_tl_scheme.go < schemes/api_layer_23.json > tl_schema.go
gofmt -w tl_schema.go
