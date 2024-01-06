package entity

type Metrics struct {
	Name              string `json:"name,omitempty"`
	TargetUtilization int    `json:"target_utilization,omitempty"`
	Type              string `json:"type,omitempty"`
}
type HPA struct {
	Name        string    `json:"name,omitempty"`
	Namespace   string    `json:"namespace,omitempty"`
	MinReplicas int       `json:"min_replica,omitempty"`
	MaxReplicas int       `json:"max_replica,omitempty"`
	Target      string    `json:"target,omitempty"`
	Metrics     []Metrics `json:"metrics,omitempty"`
}
