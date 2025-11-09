package detector

import (
	"partition-monitor/internal/config"
)

type GroupStatus struct {
	GroupName      string
	IsPartitioned  bool
	HealthyCount   int
	UnhealthyCount int
	TotalNodes     int
	QuorumRequired int
	UnhealthyNodes []string
	Message        string
}
type ClusterStatus struct {
	Groups         []GroupStatus
	AnyPartitioned bool
	AllHealthy     bool
}

type PartitionDetector struct {
	quorumGroups []config.QuorumGroup
}

func NewPartitionDetector(quorumGroups []config.QuorumGroup) *PartitionDetector {
	return &PartitionDetector{
		quorumGroups: quorumGroups,
	}
}
