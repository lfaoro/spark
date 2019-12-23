#!/usr/bin/env bash

release="$0"

npm i -g grpc-web
curl -o /tmp/grpcweb -OL https://github.com/grpc/grpc-web/releases/download/1.0.7/protoc-gen-grpc-web-"${release}"-darwin-x86_64
chmod +x /tmp/grpcweb
sudo mv /tmp/grpcweb /usr/local/bin/protoc-gen-grpc-web
