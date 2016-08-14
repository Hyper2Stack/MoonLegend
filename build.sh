#!/bin/bash

export GOPATH=$(cd `dirname $0`;pwd)
go build -o ${GOPATH}/bin/moonlegend src/controller/*.go
go build -o ${GOPATH}/bin/ml src/ml/*.go
