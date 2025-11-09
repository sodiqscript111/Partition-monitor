package checker

import (
	"fmt"
	"net"
	"net/http"
	"partition-monitor/internal/config"
	"time"
)

type HealthResult struct {
	NodeName  string
	Status    string
	Latency   time.Duration
	Error     error
	Timestamp time.Time
}

type HealthChecker struct {
	timeout    time.Duration
	httpClient *http.Client
}

func HttpHealthChecker(nodeName, nodeUrl string) HealthResult {
	result := HealthResult{
		NodeName:  nodeName,
		Timestamp: time.Now(),
	}

	start := time.Now()

	resp, err := http.Get(nodeUrl)

	result.Latency = time.Since(start)

	if err != nil {
		result.Status = "unhealthy"
		result.Error = err
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = "healthy"
	} else {
		result.Status = "unhealthy"
		result.Error = fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	return result
}

func TcpChecker(nodeName, nodeUrl string) HealthResult {
	result := HealthResult{
		NodeName:  nodeName,
		Timestamp: time.Now(),
	}

	start := time.Now()
	conn, err := net.Dial("tcp", nodeUrl)
	if err != nil {
		fmt.Println("There was an error:", err)
		result.Status = "unhealthy"
		result.Error = err
		return result
	}
	defer conn.Close()

	result.Latency = time.Since(start)
	result.Status = "healthy"
	return result
}

func HealthCheck(node config.Node) HealthResult {
	if node.Type == "http" {
		return HttpHealthChecker(node.Name, node.Address)
	}

	if node.Type == "tcp" {
		return TcpChecker(node.Name, node.Address)
	}

	// Handle unknown type
	return HealthResult{
		NodeName:  node.Name,
		Status:    "unhealthy",
		Error:     fmt.Errorf("unknown node type: %s", node.Type),
		Timestamp: time.Now(),
	}
}
