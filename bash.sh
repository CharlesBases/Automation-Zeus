#!/bin/zsh

set -e

name=reverse

go build main.go

gopath=$GOPATH
gopath=${gopath%:*}

mv main ${gopath}/bin/${name}
