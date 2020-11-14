#!/usr/bin/env bash

VERSION=v0.4.0
IMAGE=wxext-registry.101.com/sdp/sdp-devops

cd ../
go build -o bin/sdp-cleaner cmd/sdp-cleaner/sdp-cleaner.go
go build -o bin/sdp-exporter cmd/sdp-exporter/sdp-exporter.go

cd -

docker build -t ${IMAGE}:${VERSION} ./
docker push ${IMAGE}:${VERSION}