#!/bin/bash
#Don't forget to put $GOPATH/bin into PATH for gogoslick_out to work with protoc
protoc --gogoslick_out=plugins=grpc:. --proto_path=. --proto_path="$GOPATH/src" ./*.proto