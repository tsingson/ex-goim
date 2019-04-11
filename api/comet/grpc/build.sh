#!/usr/bin/env bash

protoc  -I=/Users/qinshen/go/src -I=/usr/local/include  -I=./ --gofast_out=plugins=grpc:.  ./*.proto

