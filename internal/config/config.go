package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	// General settings
	Language string `json:"language"`
	Theme    string `json:"theme"`

	// Connection settings
	ConnectionMode string `json:"connection_mode"` // proxy, tun, tap

	// Saved connections
	Connections []Connection `json:"connections"`
}

// Connection represents a VPN connection configuration
type Connection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Protocol string `json:"protocol"` // vless, shadowsocks, openvpn, wireguard, trojan, ikev2, sstp
	URL      string `json:"url"`      // Connection URL/link
	Config   string `json:"config"`   // Protocol-specific configuration
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Language:       "en",
		Theme:          "dark",
		ConnectionMode: "proxy",
		Connections:    []Connection{},
	}
}

// Load loads the configuration from the specified file
func Load(configPath string) (*Config, error) {
	// If no config path is provided, use default location
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		configPath = filepath.Join(homeDir, ".ghostmesh", "config.json")
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		cfg := DefaultConfig()
		
		// Ensure directory exists
		dir := filepath.Dir(configPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
		
		// Save default config
		if err := cfg.Save(configPath); err != nil {
			return nil, err
		}
		
		return cfg, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse config
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Save saves the configuration to the specified file
func (c *Config) Save(configPath string) error {
	// Marshal config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(configPath, data, 0644)
}

// AddConnection adds a new connection to the configuration
func (c *Config) AddConnection(conn Connection) {
	c.Connections = append(c.Connections, conn)
}

// RemoveConnection removes a connection from the configuration
func (c *Config) RemoveConnection(id string) {
	var connections []Connection
	for _, conn := range c.Connections {
		if conn.ID != id {
			connections = append(connections, conn)
		}
	}
	c.Connections = connections
}

// GetConnection returns a connection by ID
func (c *Config) GetConnection(id string) *Connection {
	for _, conn := range c.Connections {
		if conn.ID == id {
			return &conn
		}
	}
	return nil
}