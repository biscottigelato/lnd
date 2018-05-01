#!/bin/sh

# Generate APIs by passing the parsed protos to ./gen
protoc -I/usr/local/include -I. \
       -I$GOPATH/src \
       -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
       --plugin=protoc-gen-custom=./gen_bindings \
       --custom_out=./build \
       --proto_path=../lnrpc \
       rpc.proto



