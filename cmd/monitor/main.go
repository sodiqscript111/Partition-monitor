package main

import (
	"fmt"
	"partition-monitor/internal/checker"
	"partition-monitor/internal/config"
	"partition-monitor/internal/detector" // ← ADD THIS IMPORT!
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	fmt.Printf("Monitor Name: %s\n", cfg.Monitor.Name)
	fmt.Printf("Check Interval: %v\n", cfg.Monitor.CheckInterval)
	fmt.Printf("Check Timeout: %v\n", cfg.Monitor.CheckTimeout)
	fmt.Printf("Quorum Groups: %d\n\n", len(cfg.QuorumGroups)) // ← Changed from Quorum

	// Display quorum groups
	for _, group := range cfg.QuorumGroups {
		fmt.Printf("  %s: %d%% quorum, tags: %v\n",
			group.Name, group.Quorum, group.Tags)
	}

	fmt.Println("\nChecking Nodes...\n")

	// Collect all results - IMPORTANT: Store them in a slice!
	var results []checker.HealthResult // ← ADD THIS

	for _, node := range cfg.Nodes {
		if !node.Enabled {
			continue
		}

		// Perform health check
		result := checker.HealthCheck(node)
		results = append(results, result) // ← ADD THIS - Save the result!

		// Display result
		statusIcon := "✅"
		if result.Status == "unhealthy" {
			statusIcon = "❌"
		}

		fmt.Printf(" %s (%s): %s\n", result.NodeName, node.Type, node.Address)
		fmt.Printf("   Tags: %v\n", node.Tags)
		fmt.Printf("   %s Status: %s\n", statusIcon, result.Status)
		fmt.Printf("     Latency: %v\n", result.Latency)

		if result.Error != nil {
			fmt.Printf("     Error: %v\n", result.Error)
		}

		fmt.Println()
	}

	// ═══════════════════════════════════════════════════════
	// THIS IS THE NEW PART - PARTITION DETECTION!
	// ═══════════════════════════════════════════════════════

	fmt.Println("═══════════════════════════════════")
	fmt.Println("🔍 PARTITION DETECTION BY GROUP")
	fmt.Println("═══════════════════════════════════\n")

	// Create partition detector with your quorum groups
	partitionDetector := detector.NewPartitionDetector(cfg.QuorumGroups)

	// Detect partitions across all groups
	clusterStatus := partitionDetector.DetectPartitions(results, cfg.Nodes)

	// Display each group's status
	for _, groupStatus := range clusterStatus.Groups {
		if groupStatus.IsPartitioned {
			fmt.Printf(" %s: PARTITIONED\n", groupStatus.GroupName)
		} else {
			fmt.Printf("%s: HEALTHY\n", groupStatus.GroupName)
		}

		fmt.Printf("   %s\n", groupStatus.Message)

		if len(groupStatus.UnhealthyNodes) > 0 {
			fmt.Printf("   Unreachable: %v\n", groupStatus.UnhealthyNodes)
		}
		fmt.Println()
	}

	// Overall cluster status
	fmt.Println("───────────────────────────────────")
	if clusterStatus.AllHealthy {
		fmt.Println("ALL GROUPS HEALTHY - No action needed")
	} else if clusterStatus.AnyPartitioned {
		fmt.Println("PARTITION DETECTED - ALERT TRIGGERED!")
		fmt.Println("   Some services may be unavailable")
	}
}
