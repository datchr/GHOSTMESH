package vpn

import (
	"fmt"
	"sync"
	"time"

	"github.com/ghostmesh/vpnclient/internal/config"
	"github.com/ghostmesh/vpnclient/internal/logger"
	"github.com/pkg/errors"
)

// Status represents the current VPN connection status
type Status struct {
	Connected    bool
	Connecting   bool
	ConnectionID string
}

// Metrics represents VPN connection metrics
type Metrics struct {
	Speed         int64 // bytes per second
	Ping          int   // milliseconds
	Jitter        int   // milliseconds
	BytesReceived int64
	BytesSent     int64
}

// Manager manages VPN connections
type Manager struct {
	log           *logger.Logger
	connectionMode string
	currentConn   *config.Connection
	provider      Provider
	status        Status
	metrics       Metrics
	mu            sync.Mutex
	metricsTimer  *time.Ticker
	providers     map[string]Provider
}

// NewManager creates a new VPN manager
func NewManager() *Manager {
	return &Manager{
		log:           logger.GetLogger(),
		connectionMode: "proxy",
		status:        Status{},
		metrics:       Metrics{},
		providers:     make(map[string]Provider),
	}
}

// Init initializes the VPN manager
func (m *Manager) Init() error {
	// Register VPN providers
	m.registerProviders()

	// Start metrics collection
	m.startMetricsCollection()

	return nil
}

// registerProviders registers all supported VPN providers
func (m *Manager) registerProviders() {
	// Register VLESS provider
	m.providers["vless"] = &VLESSProvider{}

	// Register ShadowSocks provider
	m.providers["shadowsocks"] = &ShadowsocksProvider{}

	// Register OpenVPN provider
	m.providers["openvpn"] = &OpenVPNProvider{}

	// Register WireGuard provider
	m.providers["wireguard"] = &WireGuardProvider{}

	// Register Trojan-GFW provider
	m.providers["trojan"] = &TrojanProvider{}

	// Register IKEv2/IPsec provider
	m.providers["ikev2"] = &IKEv2Provider{}

	// Register SSTP provider
	m.providers["sstp"] = &SSTPProvider{}
}

// startMetricsCollection starts collecting metrics
func (m *Manager) startMetricsCollection() {
	m.metricsTimer = time.NewTicker(1 * time.Second)
	go func() {
		for range m.metricsTimer.C {
			m.updateMetrics()
		}
	}()
}

// updateMetrics updates the connection metrics
func (m *Manager) updateMetrics() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.status.Connected || m.provider == nil {
		return
	}

	// Get metrics from provider
	metrics, err := m.provider.GetMetrics()
	if err != nil {
		m.log.Error("Failed to get metrics: %v", err)
		return
	}

	m.metrics = metrics
}

// Connect connects to a VPN
func (m *Manager) Connect(conn *config.Connection) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if already connected
	if m.status.Connected {
		return errors.New("already connected")
	}

	// Check if connecting
	if m.status.Connecting {
		return errors.New("connection in progress")
	}

	// Get provider for protocol
	provider, ok := m.providers[conn.Protocol]
	if !ok {
		return errors.Errorf("unsupported protocol: %s", conn.Protocol)
	}

	// Update status
	m.status.Connecting = true
	m.status.ConnectionID = conn.ID

	// Initialize provider
	if err := provider.Init(conn, m.connectionMode); err != nil {
		m.status.Connecting = false
		m.status.ConnectionID = ""
		return errors.Wrap(err, "failed to initialize provider")
	}

	// Connect
	if err := provider.Connect(); err != nil {
		m.status.Connecting = false
		m.status.ConnectionID = ""
		return errors.Wrap(err, "failed to connect")
	}

	// Update status
	m.status.Connected = true
	m.status.Connecting = false
	m.currentConn = conn
	m.provider = provider

	return nil
}

// Disconnect disconnects from the VPN
func (m *Manager) Disconnect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if connected
	if !m.status.Connected {
		return errors.New("not connected")
	}

	// Disconnect
	if err := m.provider.Disconnect(); err != nil {
		return errors.Wrap(err, "failed to disconnect")
	}

	// Update status
	m.status.Connected = false
	m.status.ConnectionID = ""
	m.currentConn = nil
	m.provider = nil

	// Reset metrics
	m.metrics = Metrics{}

	return nil
}

// GetStatus returns the current VPN status
func (m *Manager) GetStatus() Status {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.status
}

// GetMetrics returns the current VPN metrics
func (m *Manager) GetMetrics() Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.metrics
}

// SetConnectionMode sets the connection mode
func (m *Manager) SetConnectionMode(mode string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate mode
	if mode != "proxy" && mode != "tun" && mode != "tap" {
		return errors.New("invalid mode: must be proxy, tun, or tap")
	}

	// Check if connected
	if m.status.Connected {
		return errors.New("cannot change mode while connected")
	}

	m.connectionMode = mode

	return nil
}

// AddConnection adds a new connection to the manager
func (m *Manager) AddConnection(conn *config.Connection) error {
	// Validate connection
	if conn == nil {
		return errors.New("connection is nil")
	}

	// Check if provider exists for protocol
	_, ok := m.providers[conn.Protocol]
	if !ok {
		return errors.Errorf("unsupported protocol: %s", conn.Protocol)
	}

	return nil
}