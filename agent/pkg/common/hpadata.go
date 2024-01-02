package common

type MetricsObject struct {
	Name              string `json:"ame,omitempty"`
	TargetUtilization int    `json:"target_utilization,omitempty"`
	Type              string `json:"type,omitempty"`
}
type HPAObject struct {
	Name        string          `json:"name,omitempty"`
	Namespace   string          `json:"namespace,omitempty"`
	MinReplicas int             `json:"min_replica,omitempty"`
	MaxReplicas int             `json:"max_replica,omitempty"`
	Target      string          `json:"target,omitempty"`
	Metrics     []MetricsObject `json:"metrics,omitempty"`
}
