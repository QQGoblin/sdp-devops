package kubernetes

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

const (
	DefaultCRIEndpoint = "unix:///var/run/crio/crio.sock"
)

// 创建gRPC连接
func CRIRuntimeClient(endPoint string) (pb.RuntimeServiceClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(endPoint, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		errMsg := errors.Wrapf(err, "connect endpoint '%s', make sure you are running as root and the endpoint has been started", endPoint)
		logrus.Error(errMsg)
		return nil, conn
	} else {
		logrus.Debugf("connected successfully using endpoint: %s", endPoint)
	}
	runtimeClient := pb.NewRuntimeServiceClient(conn)
	return runtimeClient, conn
}

// 关闭gRPC连接
func CloseCRIConn(conn *grpc.ClientConn) {
	if conn != nil {
		conn.Close()
	}
}
