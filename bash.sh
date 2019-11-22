#!/bin/zsh

set -e

name=db

go build main.go

gopath=$GOPATH
gopath=${gopath%:*}

mv main ${gopath}/bin/${name}
