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


# 创建相应证书

#kubectl  create secret tls etcd-healthcheck-client --cert=/etc/etcd/healthcheck-client.crt --key=/etc/etcd/healthcheck-client.key --namespace sdp-devops
#kubectl  create secret tls k8s-admin-client --cert=/etc/kubernetes/admin.crt --key=/etc/kubernetes/admin.key --namespace sdp-devops