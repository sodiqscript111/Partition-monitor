package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	
)

type Config struct {
	Monitor MonitorConfig `yaml:"monitor"`
	QuorumGroups []QuorumGroup `yaml:"quorum_groups"`
	Nodes   []Node        `yaml:"nodes"`
}

type MonitorConfig struct {
	Name          string        `yaml:"name"`
	CheckInterval time.Duration `yaml:"check_interval"` // ← Fixed: matches YAML + correct type
	CheckTimeout  time.Duration `yaml:"check_timeout"`  // ← Fixed: matches YAML + correct type
	Quorum        int           `yaml:"quorum"`
}

type Node struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Address string   `yaml:"address"`
	Enabled bool     `yaml:"enabled"`
	Tags    []string `yaml:"tags"`
}

type QuorumGroup struct {
	Name    string   `yaml:"name"`
    Quorum int      `yaml:"quorum"`
	Tags    []string `yaml:"tags"`
}

func LoadConfig(path string) (*Config, error) {  // ← Fixed typo: "LoadConfing" → "LoadConfig"
func LoadConfig(path string) (*Config, error) { // ← Fixed typo: "LoadConfing" → "LoadConfig"
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)  // Better error message
		return nil, fmt.Errorf("failed to read config file: %w", err) // Better error message


	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)  // Better error message
		return nil, fmt.Errorf("failed to parse YAML: %w", err) // Better error message


}
