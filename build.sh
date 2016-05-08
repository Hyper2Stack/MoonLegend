#!/bin/bash

export GOPATH=$(readlink -f $(dirname $0))
go build src/controller/moonlegend.go src/controller/parsejson.go
