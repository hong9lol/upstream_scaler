package common

type Deployment struct {
	Name string `json:"deployment_name"`
	Pods []Pod  `json:"pods"`
}

type Pod struct {
	Name   string  `json:"pod_name"`
	Usages []Usage `json:"usages"`
}

type Usage struct {
	// reference: https://github.com/opencontainers/runc/blob/35988abe20851979d53d4a8f790b4c1b8800a10d/types/events.go#L77
	Usage     uint64 `json:"usage"`
	Timestamp uint64 `json:"timestamp"` // time.Now().Unix()
}
