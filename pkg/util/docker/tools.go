package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"path"
	"sdp-devops/pkg/util/sys"
	"strings"
)

var (
	defaultAPIVerison = "1.25"
	DockerInfo        types.Info
)

// 返回Docker Client
func DockerClient(host string) *client.Client {

	var c *client.Client
	var err error
	if strings.EqualFold(host, "") {
		c, err = client.NewEnvClient()
	} else {
		c, err = client.NewClient(host, defaultAPIVerison, nil, nil)
	}
	if err != nil {
		panic(err)
	}

	return c
}

// 获取容器运行时占用的磁盘空间
func ContainerSize(containerID string, cli *client.Client) int64 {
	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		logrus.Error(err.Error())
		return 0
	}
	upperDir := containerInfo.GraphDriver.Data["UpperDir"]
	// TODO: 这里改成系统调用会不会更快？
	upperDirSize, _ := sys.CalDirSize(upperDir)
	return upperDirSize
}

// 获取容器日志的磁盘使用空间
func ContainerLogSize(containerID, dockerRootDir string, cli *client.Client) int64 {

	containerDataPath := path.Join(dockerRootDir, "containers", containerID)
	logSize, _ := sys.CalDirSize(containerDataPath)
	return logSize
}
