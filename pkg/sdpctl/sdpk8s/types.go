package sdpk8s

// 打印Pod的状态信息
type PodBriefInfo struct {
	Name        string
	UID         string
	NameSpace   string
	Status      string
	Node        string
	NetworkMode string
	Restart     int
	PID         int
}

type NodeBriefInfo struct {
	Name          string
	Role          string
	UnSche        string
	Env           string
	Type          string
	Label         string
	CPU           string
	Memory        string
	MemoryUsage   string
	MemoryRequest string
	MemoryLimits  string
	Pod           string
}
