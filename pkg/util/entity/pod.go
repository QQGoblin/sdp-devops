package entity

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
