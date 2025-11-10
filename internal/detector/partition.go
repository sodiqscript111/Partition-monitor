package detector

import (
	"fmt"
	"partition-monitor/internal/checker"
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
	fmt.Println("🔧 Creating PartitionDetector with", len(quorumGroups), "groups")
	for _, group := range quorumGroups {
		fmt.Printf("   - %s: %d%% quorum, tags: %v\n", group.Name, group.Quorum, group.Tags)
	}
	fmt.Println()

	return &PartitionDetector{
		quorumGroups: quorumGroups,
	}
}

func nodeInGroup(node config.Node, groupTags []string) bool {
	fmt.Printf("  Checking if node '%s' (tags: %v) matches group tags: %v\n",
		node.Name, node.Tags, groupTags)

	for _, nodeTag := range node.Tags {
		for _, groupTag := range groupTags {
			if nodeTag == groupTag {
				fmt.Printf("       MATCH! Found '%s'\n", nodeTag)
				return true
			}
		}
	}

	fmt.Println("       No match")
	return false
}

func (pd *PartitionDetector) DetectPartitions(
	allResults []checker.HealthResult,
	allNodes []config.Node,
) ClusterStatus {

	fmt.Println(" DETECTING PARTITIONS...")
	fmt.Printf("   Total results: %d, Total nodes: %d\n\n", len(allResults), len(allNodes))

	var groupStatuses []GroupStatus

	// Loop through each group
	for _, group := range pd.quorumGroups {
		fmt.Printf(" Checking group: %s\n", group.Name)

		// Mini-Step B: Find nodes in this group
		var nodesInThisGroup []config.Node

		fmt.Println("   Finding nodes for this group...")
		for _, node := range allNodes {
			if !node.Enabled {
				fmt.Printf("   Skipping disabled node: %s\n", node.Name)
				continue
			}

			if nodeInGroup(node, group.Tags) {
				nodesInThisGroup = append(nodesInThisGroup, node)
			}
		}

		fmt.Printf("    Nodes in group: %d\n", len(nodesInThisGroup))
		for _, n := range nodesInThisGroup {
			fmt.Printf("      - %s\n", n.Name)
		}
		fmt.Println()

		// Mini-Step C: Find results for those nodes
		var resultsForThisGroup []checker.HealthResult

		fmt.Println("   Finding health results for these nodes...")
		for _, result := range allResults {
			for _, node := range nodesInThisGroup {
				if result.NodeName == node.Name {
					resultsForThisGroup = append(resultsForThisGroup, result)
					fmt.Printf("       Found result for: %s (status: %s)\n",
						result.NodeName, result.Status)
					break
				}
			}
		}
		fmt.Println()

		// Mini-Step D: Count healthy vs unhealthy
		healthyCount := 0
		var unhealthyNodes []string

		fmt.Println("   Counting healthy vs unhealthy...")
		for _, result := range resultsForThisGroup {
			if result.Status == "healthy" {
				healthyCount++
				fmt.Printf("  %s: healthy\n", result.NodeName)
			} else {
				unhealthyNodes = append(unhealthyNodes, result.NodeName)
				fmt.Printf("      ❌ %s: unhealthy\n", result.NodeName)
			}
		}

		totalNodes := len(resultsForThisGroup)
		fmt.Printf("    Summary: %d healthy, %d unhealthy, %d total\n\n",
			healthyCount, len(unhealthyNodes), totalNodes)

		// Mini-Step E: Calculate quorum
		quorumRequired := (totalNodes*group.Quorum + 99) / 100
		isPartitioned := healthyCount < quorumRequired

		fmt.Printf("    Quorum calculation:\n")
		fmt.Printf("      Total nodes: %d\n", totalNodes)
		fmt.Printf("      Quorum percentage: %d%%\n", group.Quorum)
		fmt.Printf("      Required healthy: %d\n", quorumRequired)
		fmt.Printf("      Actual healthy: %d\n", healthyCount)
		fmt.Printf("      Result: %v\n\n", map[bool]string{true: "🚨 PARTITIONED", false: "✅ HEALTHY"}[isPartitioned])

		// Mini-Step F: Create message
		var message string
		if isPartitioned {
			message = fmt.Sprintf("PARTITIONED: %d/%d nodes healthy (need %d)",
				healthyCount, totalNodes, quorumRequired)
		} else {
			message = fmt.Sprintf("Healthy: %d/%d nodes reachable (quorum: %d)",
				healthyCount, totalNodes, quorumRequired)
		}

		status := GroupStatus{
			GroupName:      group.Name,
			IsPartitioned:  isPartitioned,
			HealthyCount:   healthyCount,
			UnhealthyCount: totalNodes - healthyCount,
			TotalNodes:     totalNodes,
			QuorumRequired: quorumRequired,
			UnhealthyNodes: unhealthyNodes,
			Message:        message,
		}

		groupStatuses = append(groupStatuses, status)
		fmt.Println("   Group status created")
		fmt.Println("═══════════════════════════════════\n")
	}

	// Mini-Step H: Calculate overall status
	fmt.Println(" Calculating overall cluster status...")

	anyPartitioned := false
	allHealthy := true

	for _, status := range groupStatuses {
		if status.IsPartitioned {
			anyPartitioned = true
			allHealthy = false
			fmt.Printf("    %s is partitioned\n", status.GroupName)
		} else {
			fmt.Printf(" %s is healthy\n", status.GroupName)
		}
	}

	fmt.Printf("\n   Final: AnyPartitioned=%v, AllHealthy=%v\n\n", anyPartitioned, allHealthy)

	return ClusterStatus{
		Groups:         groupStatuses,
		AnyPartitioned: anyPartitioned,
		AllHealthy:     allHealthy,
	}
}
