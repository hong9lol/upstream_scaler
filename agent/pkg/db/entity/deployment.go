package entity

type Deployment struct {
	Name string         `json:"deployment_name"`
	Pods map[string]Pod `json:"pods"`
}

type Pod struct {
	Name       string               `json:"pod_name"`
	AppLabel   string               `json:"label"` // app label to matching deployment
	Containers map[string]Container `json:"container"`
}

type Container struct {
	Name       string  `json:"container_name"`
	CPURequest int64   `json:"cpu_request"`
	Usages     []Usage `json:"usages"`
}

type Usage struct {
	// reference: https://github.com/opencontainers/runc/blob/35988abe20851979d53d4a8f790b4c1b8800a10d/types/events.go#L77
	Usage     int64 `json:"usage"`
	Timestamp int64 `json:"timestamp"` // time.Now().Unix()
}
