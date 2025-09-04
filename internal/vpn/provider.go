package vpn

import (
	"github.com/ghostmesh/vpnclient/internal/config"
)

// Provider is the interface for VPN providers
type Provider interface {
	// Init initializes the provider with the given connection and mode
	Init(conn *config.Connection, mode string) error

	// Connect connects to the VPN
	Connect() error

	// Disconnect disconnects from the VPN
	Disconnect() error

	// GetMetrics returns the current metrics
	GetMetrics() (Metrics, error)
}