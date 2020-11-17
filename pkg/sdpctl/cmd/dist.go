package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"os"
	"path"
	"sdp-devops/pkg/sdpctl/config"
	"sdp-devops/pkg/sdpctl/sdpk8s"
	k8stools "sdp-devops/pkg/util/kubernetes"
	systools "sdp-devops/pkg/util/sys"
	"strconv"
	"strings"
)

var (
	timeOut int
)

func RunDist(cmd *cobra.Command, args []string) {

	if len(args) < 2 {
		panic("请填写源路径和目标路径")
	}
	srcPath := path.Base(args[0])
	destPath := path.Join(path.Clean(args[1]) + "/" + srcPath)
	if !path.IsAbs(destPath) {
		panic("目标文件请指定绝对路径")
	}
	destDir := path.Dir(destPath)
	cmdArr := []string{"tar", "-xmf", "-", "-C", destDir}

	kubeClientSet, kubeClientConfig := k8stools.KubeClientAndConfig(config.KubeConfigStr)
	// 返回所有需要运行运行的Node列表
	shellPodTargets := sdpk8s.GetShellPodDict(kubeClientSet)
	i := 0
	for _, v := range shellPodTargets {
		logrus.Print("------------------------------> No." + strconv.Itoa(i) + " Dist on node: " + v.Spec.NodeName + " <------------------------------")

		reader, writer := io.Pipe()
		tarExecOps := k8stools.ExecOptions{
			Command:       strings.Join(cmdArr, " "),
			ContainerName: "",
			In:            reader,
			Out:           os.Stdout,
			Err:           os.Stderr,
			Istty:         false,
			TimeOut:       timeOut,
		}

		go func() {
			defer writer.Close()
			if err := systools.MakeTar(srcPath, destPath, writer); err != nil {
				logrus.Error(err.Error())
			}
		}()

		if err := k8stools.ExecCmd(kubeClientSet, kubeClientConfig, v, tarExecOps); err != nil {
			logrus.Error(err.Error())
		}
		i += 1
	}
}

func addDistFlag(flags *pflag.FlagSet) {
	flags.IntVar(&timeOut, "timeout", 15, "单个文件传输超时时间。")
}

func NewCmdDist() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "dist [local file path] [dist folder path]",
		Short:                 "分发文件到宿主机目录，注意同名文件会被覆盖",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			RunDist(cmd, args)
		},
	}
	addDistFlag(cmd.Flags())
	return cmd
}
