package node

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// NodeManager handles connections to multiple Ethereum nodes
type NodeManager struct {
	primaryClient *ethclient.Client
	backupClients []*ethclient.Client
	wsClients     []*ethclient.Client
	healthChecks  map[string]bool
	mutex         sync.RWMutex
	endpoints     []string
	wsEndpoints   []string
}

// NewNodeManager creates a new node manager
func NewNodeManager(rpcEndpoints, wsEndpoints []string) (*NodeManager, error) {
	if len(rpcEndpoints) == 0 {
		return nil, fmt.Errorf("at least one RPC endpoint is required")
	}

	manager := &NodeManager{
		healthChecks: make(map[string]bool),
		endpoints:    rpcEndpoints,
		wsEndpoints:  wsEndpoints,
	}

	// Connect to primary client
	var err error
	manager.primaryClient, err = ethclient.Dial(rpcEndpoints[0])
	if err != nil {
		return nil, fmt.Errorf("connecting to primary endpoint: %w", err)
	}

	// Connect to backup clients
	for i := 1; i < len(rpcEndpoints); i++ {
		client, err := ethclient.Dial(rpcEndpoints[i])
		if err != nil {
			fmt.Printf("Warning: Failed to connect to backup endpoint %s: %v\n", rpcEndpoints[i], err)
			continue
		}
		manager.backupClients = append(manager.backupClients, client)
	}

	// Connect to WebSocket clients
	for _, endpoint := range wsEndpoints {
		client, err := ethclient.Dial(endpoint)
		if err != nil {
			fmt.Printf("Warning: Failed to connect to WebSocket endpoint %s: %v\n", endpoint, err)
			continue
		}
		manager.wsClients = append(manager.wsClients, client)
	}

	// Initialize health checks
	for _, endpoint := range rpcEndpoints {
		manager.healthChecks[endpoint] = true
	}

	return manager, nil
}

// GetClient returns a healthy client
func (m *NodeManager) GetClient() *ethclient.Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check if primary client is healthy
	if m.healthChecks[m.endpoints[0]] {
		return m.primaryClient
	}

	// Select a random backup client
	if len(m.backupClients) > 0 {
		idx := rand.Intn(len(m.backupClients))
		return m.backupClients[idx]
	}

	// Fallback to primary even if it's unhealthy
	return m.primaryClient
}

// GetWSClient returns a WebSocket client for subscriptions
func (m *NodeManager) GetWSClient() *ethclient.Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if len(m.wsClients) == 0 {
		return nil
	}

	// Select a random WebSocket client
	idx := rand.Intn(len(m.wsClients))
	return m.wsClients[idx]
}

// StartHealthCheck begins periodic health checks of nodes
func (m *NodeManager) StartHealthCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.checkHealth(ctx)
		}
	}
}

// checkHealth verifies connectivity of all nodes
func (m *NodeManager) checkHealth(ctx context.Context) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check primary
	if m.primaryClient != nil {
		_, err := m.primaryClient.BlockNumber(ctx)
		m.healthChecks[m.endpoints[0]] = err == nil
	}

	// Check backups
	for i, client := range m.backupClients {
		if i+1 < len(m.endpoints) {
			_, err := client.BlockNumber(ctx)
			m.healthChecks[m.endpoints[i+1]] = err == nil
		}
	}
}

// Close closes all connections
func (m *NodeManager) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.primaryClient != nil {
		m.primaryClient.Close()
	}

	for _, client := range m.backupClients {
		client.Close()
	}

	for _, client := range m.wsClients {
		client.Close()
	}
}
