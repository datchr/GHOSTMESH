package vpn

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/ghostmesh/vpnclient/internal/config"
	"github.com/ghostmesh/vpnclient/internal/logger"
	"github.com/pkg/errors"
)

// VLESSProvider implements the Provider interface for VLESS protocol
type VLESSProvider struct {
	log        *logger.Logger
	conn       *config.Connection
	mode       string
	cmd        *exec.Cmd
	config     *VLESSConfig
	lastRx     int64
	lastTx     int64
	lastUpdate time.Time
}

// VLESSConfig represents the VLESS configuration
type VLESSConfig struct {
	ID        string `json:"id"`
	Address   string `json:"address"`
	Port      int    `json:"port"`
	TLS       bool   `json:"tls"`
	Reality   bool   `json:"reality"`
	Vision    bool   `json:"vision"`
	ServerName string `json:"server_name"`
	Fingerprint string `json:"fingerprint"`
	PublicKey string `json:"public_key"`
	ShortID   string `json:"short_id"`
}

// Init initializes the VLESS provider
func (p *VLESSProvider) Init(conn *config.Connection, mode string) error {
	p.log = logger.GetLogger()
	p.conn = conn
	p.mode = mode
	p.lastUpdate = time.Now()

	// Parse config
	var config VLESSConfig
	if err := json.Unmarshal([]byte(conn.Config), &config); err != nil {
		return errors.Wrap(err, "failed to parse VLESS config")
	}

	p.config = &config

	return nil
}

// Connect connects to the VLESS server
func (p *VLESSProvider) Connect() error {
	// Prepare xray-core command
	configPath, err := p.writeConfig()
	if err != nil {
		return errors.Wrap(err, "failed to write config")
	}

	// Start xray-core
	p.cmd = exec.Command("xray", "-c", configPath)
	if err := p.cmd.Start(); err != nil {
		return errors.Wrap(err, "failed to start xray-core")
	}

	// Wait for connection to establish
	time.Sleep(1 * time.Second)

	return nil
}

// Disconnect disconnects from the VLESS server
func (p *VLESSProvider) Disconnect() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}

	// Kill xray-core process
	if err := p.cmd.Process.Kill(); err != nil {
		return errors.Wrap(err, "failed to kill xray-core process")
	}

	return nil
}

// GetMetrics returns the current connection metrics
func (p *VLESSProvider) GetMetrics() (Metrics, error) {
	// In a real implementation, we would get metrics from xray-core
	// For now, we'll return dummy metrics
	now := time.Now()
	elapsed := now.Sub(p.lastUpdate).Seconds()

	// Simulate traffic
	rx := p.lastRx + int64(1000000*elapsed) // 1 MB/s
	tx := p.lastTx + int64(500000*elapsed)  // 500 KB/s

	// Calculate speed
	rxSpeed := int64(float64(rx-p.lastRx) / elapsed)
	txSpeed := int64(float64(tx-p.lastTx) / elapsed)

	// Update last values
	p.lastRx = rx
	p.lastTx = tx
	p.lastUpdate = now

	return Metrics{
		Speed:         rxSpeed + txSpeed,
		Ping:          50,  // 50ms
		Jitter:        5,   // 5ms
		BytesReceived: rx,
		BytesSent:     tx,
	}, nil
}

// ParseURL parses a VLESS URL and returns a connection config
func (p *VLESSProvider) ParseURL(urlStr string) (*config.Connection, error) {
	// Check if URL is a VLESS URL
	if !strings.HasPrefix(urlStr, "vless://") {
		return nil, errors.New("not a VLESS URL")
	}

	// Parse URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse URL")
	}

	// Extract ID from username
	id := u.User.Username()

	// Extract host and port
	hostParts := strings.Split(u.Host, ":")
	if len(hostParts) != 2 {
		return nil, errors.New("invalid host format")
	}

	host := hostParts[0]
	port := 0
	fmt.Sscanf(hostParts[1], "%d", &port)

	// Parse query parameters
	query := u.Query()

	// Create config
	config := VLESSConfig{
		ID:        id,
		Address:   host,
		Port:      port,
		TLS:       query.Get("security") == "tls",
		Reality:   query.Get("reality") == "1",
		Vision:    query.Get("flow") == "xtls-rprx-vision",
		ServerName: query.Get("sni"),
		Fingerprint: query.Get("fp"),
		PublicKey: query.Get("pbk"),
		ShortID:   query.Get("sid"),
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal config")
	}

	// Generate name from URL
	name := fmt.Sprintf("VLESS - %s:%d", host, port)

	// Create connection
	conn := &config.Connection{
		ID:       generateID(),
		Name:     name,
		Protocol: "vless",
		URL:      urlStr,
		Config:   string(configJSON),
	}

	return conn, nil
}

// writeConfig writes the VLESS configuration to a temporary file
func (p *VLESSProvider) writeConfig() (string, error) {
	// In a real implementation, we would write the config to a file
	// For now, we'll just return a dummy path
	return "/tmp/vless.json", nil
}

// generateID generates a random ID for a connection
func generateID() string {
	return fmt.Sprintf("conn-%d", time.Now().UnixNano())
}