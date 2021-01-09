module sdp-devops

go 1.15

require (
	github.com/Microsoft/go-winio v0.4.15 // indirect
	github.com/deckarep/golang-set v1.7.1
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fatih/color v1.10.0
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0
	github.com/prometheus/common v0.15.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07 // indirect
	github.com/toolkits/nux v0.0.0-20200401110743-debb3829764a
	github.com/toolkits/slice v0.0.0-20141116085117-e44a80af2484 // indirect
	github.com/toolkits/sys v0.0.0-20170615103026-1f33b217ffaf // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b // indirect
	golang.org/x/sys v0.0.0-20210104204734-6f8348627aad // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d // indirect
	google.golang.org/grpc v1.34.0
	k8s.io/api v0.16.15
	k8s.io/apimachinery v0.16.15
	k8s.io/client-go v0.16.15
	k8s.io/cri-api v0.16.15
	k8s.io/kubernetes v1.16.15
	k8s.io/metrics v0.16.15
)

replace (
	k8s.io/api => k8s.io/api v0.16.15
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.15
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.15
	k8s.io/apiserver => k8s.io/apiserver v0.16.15
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.15
	k8s.io/client-go => k8s.io/client-go v0.16.15
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.15
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.15
	k8s.io/code-generator => k8s.io/code-generator v0.16.15
	k8s.io/component-base => k8s.io/component-base v0.16.15
	k8s.io/component-helpers => k8s.io/component-helpers v0.16.15
	k8s.io/controller-manager => k8s.io/controller-manager v0.16.15
	k8s.io/cri-api => k8s.io/cri-api v0.16.15
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.15
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.15
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.15
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.15
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.15
	k8s.io/kubectl => k8s.io/kubectl v0.16.15
	k8s.io/kubelet => k8s.io/kubelet v0.16.15
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.15
	k8s.io/metrics => k8s.io/metrics v0.16.15
	k8s.io/mount-utils => k8s.io/mount-utils v0.16.15
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.15
)
