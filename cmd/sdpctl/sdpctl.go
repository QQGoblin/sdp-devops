package main

import (
	sdpLogger "sdp-devops/pkg/logger"
	"sdp-devops/pkg/sdpctl"
)

func main() {
	sdpLogger.InitLogger()
	sdpctl.Main()
}
