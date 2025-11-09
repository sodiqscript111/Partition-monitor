# Partition Monitor

Partition Monitor is a lightweight Go-based tool for continuously monitoring distributed nodes and detecting network partitions across quorum groups. It performs periodic health checks and reports when groups become partitioned or recover.

---

##  Features

- Periodic node health checks (HTTP/TCP)
- Partition detection across quorum-based groups
- State change alerts (partitioned → healthy or vice versa)
- Graceful shutdown with signal handling
- Simple YAML-based configuration

---

## ⚙️ Configuration

All settings are defined in `config.yaml`.  
Example structure:

```yaml
monitor:
  name: "Cluster Monitor"
  check_interval: 10s
  check_timeout: 3s

nodes:
  - name: "Node A"
    address: "http://localhost:8080/health"
    type: "http"
    enabled: true
    tags: ["api", "region1"]

  - name: "Node B"
    address: "10.0.0.2:9000"
    type: "tcp"
    enabled: true
    tags: ["db", "region1"]

quorum_groups:
  - name: "Region1"
    quorum: 60
    tags: ["region1"]

```

Usage
Build
```
go build -o partition-monitor
```

Run
```
./partition-monitor
```


-When started, the monitor:

-Loads the configuration file

-Runs an initial health check immediately

-Continues running checks at the configured interval

-Displays node and group health in the console
