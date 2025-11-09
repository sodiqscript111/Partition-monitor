package main

import (
	"fmt"
	"partition-monitor/internal/checker"
	"partition-monitor/internal/config"
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
	fmt.Printf("Quorum: %d%%\n\n", cfg.Monitor.Quorum)

	fmt.Println("Checking Nodes...\n")

	for _, node := range cfg.Nodes {
		if !node.Enabled {
			continue // Skip disabled nodes
		}

		fmt.Printf("%s (%s): %s\n", node.Name, node.Type, node.Address)
		fmt.Printf("   Tags: %v\n", node.Tags)

		// Perform health check
		result := checker.HealthCheck(node)

		// Display result with proper formatting
		statusIcon := "✅"
		if result.Status == "unhealthy" {
			statusIcon = "❌"
		}

		fmt.Printf("   %s Status: %s\n", statusIcon, result.Status)
		fmt.Printf("    Latency: %v\n", result.Latency)

		if result.Error != nil {
			fmt.Printf("    Error: %v\n", result.Error)
		}

		fmt.Println() // Blank line between nodes
	}
}
