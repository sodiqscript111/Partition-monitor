package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"partition-monitor/internal/checker"
	"partition-monitor/internal/config"
	"partition-monitor/internal/detector"
)

var (
	previousClusterStatus detector.ClusterStatus
	runCount              int
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	// Display basic info
	fmt.Printf("Starting %s\n", cfg.Monitor.Name)
	fmt.Printf("Monitoring %d quorum groups\n", len(cfg.QuorumGroups))
	fmt.Printf("Check interval: %v\n\n", cfg.Monitor.CheckInterval)

	// Show quorum groups
	for _, group := range cfg.QuorumGroups {
		fmt.Printf("  %s: %d%% quorum, tags: %v\n",
			group.Name, group.Quorum, group.Tags)
	}
	fmt.Println()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Run initial check immediately
	runHealthChecks(cfg)

	// Setup ticker for periodic checks
	ticker := time.NewTicker(cfg.Monitor.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			runHealthChecks(cfg)
		case <-sigChan:
			fmt.Println("\nShutting down gracefully...")
			return
		}
	}
}

func runHealthChecks(cfg *config.Config) {
	runCount++

	fmt.Printf("\n╔═══════════════════════════════════╗\n")
	fmt.Printf("║ Health Check Run #%d\n", runCount)
	fmt.Printf("║ %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf("╚═══════════════════════════════════╝\n\n")

	var results []checker.HealthResult

	// Perform checks for all nodes
	for _, node := range cfg.Nodes {
		if !node.Enabled {
			continue
		}

		result := checker.HealthCheck(node)
		results = append(results, result)
		printNodeResult(result, node)
	}

	fmt.Println("\nPartition Detection:")
	detectorInstance := detector.NewPartitionDetector(cfg.QuorumGroups)
	currentStatus := detectorInstance.DetectPartitions(results, cfg.Nodes)

	for _, group := range currentStatus.Groups {
		printGroupStatus(group)
	}

	// Skip comparison on first run
	if runCount > 1 {
		checkStateChanges(currentStatus)
	}

	previousClusterStatus = currentStatus
}

func printNodeResult(result checker.HealthResult, node config.Node) {
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

func printGroupStatus(group detector.GroupStatus) {
	if group.IsPartitioned {
		fmt.Printf("   %s: PARTITIONED\n", group.GroupName)
	} else {
		fmt.Printf("  %s: HEALTHY\n", group.GroupName)
	}
	fmt.Printf("   %s\n", group.Message)
	if len(group.UnhealthyNodes) > 0 {
		fmt.Printf("   Unreachable: %v\n", group.UnhealthyNodes)
	}
	fmt.Println()
}

func checkStateChanges(current detector.ClusterStatus) {
	for i, group := range current.Groups {
		prevGroup := previousClusterStatus.Groups[i]

		if !prevGroup.IsPartitioned && group.IsPartitioned {
			sendAlert(group)
		} else if prevGroup.IsPartitioned && !group.IsPartitioned {
			sendRecoveryAlert(group)
		}
	}
}

func sendAlert(group detector.GroupStatus) {
	fmt.Printf(" ALERT: %s became PARTITIONED!\n", group.GroupName)
}

func sendRecoveryAlert(group detector.GroupStatus) {
	fmt.Printf(" RECOVERY: %s is now HEALTHY again!\n", group.GroupName)
}
