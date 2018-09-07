#!/bin/bash

set -e
timetag=`date +%Y%m%d%H%M%S`

export GOPATH
cd servers/
export GOPATH="/data/jenkins/go"
export GOROOT="/usr/local/go"
go env
GOOS=linux GOARCH=amd64 go build