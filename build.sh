#!/bin/bash

export GOPATH=$(readlink -f $(dirname $0))
go build -o ${GOPATH}/bin/moonlegend src/controller/*.go
