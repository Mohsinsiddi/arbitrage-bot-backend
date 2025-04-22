package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config holds the application configuration
type Config struct {
	RPCEndpoints    []string          `yaml:"rpcEndpoints"`
	WSEndpoints     []string          `yaml:"wsEndpoints"`
	PrivateKeys     []string          `yaml:"privateKeys"`
	ContractAddress string            `yaml:"contractAddress"`
	TokenAddresses  map[string]string `yaml:"tokenAddresses"`
	DEXRouters      map[string]string `yaml:"dexRouters"`
	MinProfit       string            `yaml:"minProfit"`
	GasSettings     GasSettings       `yaml:"gasSettings"`
}

// GasSettings contains gas-related configuration
type GasSettings struct {
	MaxGasPrice      string `yaml:"maxGasPrice"`
	PriorityFee      string `yaml:"priorityFee"`
	GasLimitMultiplier float64 `yaml:"gasLimitMultiplier"`
	SpeedUpThreshold  int     `yaml:"speedUpThreshold"`
}

// LoadConfig reads and parses the configuration file
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	// Validate config
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validateConfig checks if the configuration is valid
func validateConfig(cfg *Config) error {
	if len(cfg.RPCEndpoints) == 0 {
		return fmt.Errorf("no RPC endpoints provided")
	}
	if len(cfg.PrivateKeys) == 0 {
		return fmt.Errorf("no private keys provided")
	}
	if cfg.ContractAddress == "" {
		return fmt.Errorf("contract address is required")
	}
	return nil
}
