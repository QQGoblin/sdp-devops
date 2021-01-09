#!/usr/bin/env bash

VERSION=v0.4.4-beta
IMAGE=wxext-registry.101.com/sdp/sdp-devops

cd ../
go build -o build/bin/sdp-cleaner cmd/sdp_cleaner/sdp_cleaner.go
go build -o build/bin/sdp-exporter cmd/sdp_exporter/sdp_exporter.go
go build -o build/bin/sdp-alert cmd/sdp_alert/sdp_alert.go

cd -

docker build -t ${IMAGE}:${VERSION} ./
docker push ${IMAGE}:${VERSION}