package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"path"
	"sdp-devops/pkg/util/entity"
	"sdp-devops/pkg/util/sys"
	"strconv"
	"strings"
)

var (
	defaultAPIVerison = "1.25"
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

// 根据 K8S 容器的命名规则从 docker 容器名称解析对应的Pod信息
// 映射规则如下：
// k8s_POD_${POD_FULL_NAME}_${NAMESPACES}_${POD_ID}_${Restart_Count}
// k8s_${Container_NAME}_${POD_FULL_NAME}_${NAMESPACES}_${POD_ID}_${Restart_Count}
//
func GetPodInfo(cli *client.Client) []entity.PodBriefInfo {

	podBriefInfoList := make([]entity.PodBriefInfo, 0)
	if containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{}); err != nil {
		logrus.Error(err)
	} else {
		for _, container := range containers {
			if strings.EqualFold(container.State, "running") && strings.HasPrefix(container.Names[0], "/k8s_POD_") {
				infoStr := strings.Replace(container.Names[0], "/k8s_POD_", "", -1)
				info := strings.Split(infoStr, "_")
				podBriefInfo := entity.PodBriefInfo{
					Name:        info[0],
					UID:         info[2],
					NameSpace:   info[1],
					Status:      container.State,
					Node:        "local",
					NetworkMode: container.HostConfig.NetworkMode,
					Restart:     0,
					PID:         1,
				}
				if restartCount, err := strconv.Atoi(info[3]); err != nil {
					logrus.Error(err)
					continue
				} else {
					podBriefInfo.Restart = restartCount
				}
				if containerInspect, err := cli.ContainerInspect(context.Background(), container.ID); err == nil {
					podBriefInfo.PID = containerInspect.State.Pid
				} else {
					logrus.Error(err)
					continue
				}
				podBriefInfoList = append(podBriefInfoList, podBriefInfo)
			}
		}
	}
	return podBriefInfoList
}
