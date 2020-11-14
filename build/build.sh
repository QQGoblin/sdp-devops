#!/usr/bin/env bash

VERSION=v0.4.1
IMAGE=wxext-registry.101.com/sdp/sdp-devops

cd ../
go build -o build/bin/sdp-cleaner cmd/sdp_cleaner/sdp_cleaner.go
go build -o build/bin/sdp-exporter cmd/sdp_exporter/sdp_exporter.go

cd -

docker build -t ${IMAGE}:${VERSION} ./
docker push ${IMAGE}:${VERSION}