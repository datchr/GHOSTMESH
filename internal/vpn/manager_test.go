package vpn

import (
	"testing"

	"github.com/ghostmesh/vpnclient/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestManager_AddConnection(t *testing.T) {
	// Create a new manager
	manager := NewManager()

	// Initialize the manager
	err := manager.Init()
	assert.NoError(t, err)

	// Create a test connection
	conn := &config.Connection{
		ID:       "test-connection",
		Name:     "Test Connection",
		Protocol: "vless",
		Config:   `{"address":"example.com","port":443,"uuid":"8a45f2f0-5361-4e2b-b330-1d756c448c59","encryption":"none","security":"tls","sni":"example.com"}`,
	}

	// Add the connection
	err = manager.AddConnection(conn)
	assert.NoError(t, err)

	// Test with unsupported protocol
	connInvalid := &config.Connection{
		ID:       "invalid-connection",
		Name:     "Invalid Connection",
		Protocol: "unsupported",
		Config:   `{}`,
	}

	// Add the connection
	err = manager.AddConnection(connInvalid)
	assert.Error(t, err)

	// Test with nil connection
	err = manager.AddConnection(nil)
	assert.Error(t, err)
}

func TestIntegration_ParserAndManager(t *testing.T) {
	// Create a new parser
	parser := NewParser()

	// Create a new manager
	manager := NewManager()

	// Initialize the manager
	err := manager.Init()
	assert.NoError(t, err)

	// Parse a VLESS URL
	url := "vless://8a45f2f0-5361-4e2b-b330-1d756c448c59@example.com:443?encryption=none&security=tls&sni=example.com"
	conn, err := parser.Parse(url)
	assert.NoError(t, err)

	// Add the connection to the manager
	err = manager.AddConnection(conn)
	assert.NoError(t, err)

	// Parse a Shadowsocks URL
	url = "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388"
	conn, err = parser.Parse(url)
	assert.NoError(t, err)

	// Add the connection to the manager
	err = manager.AddConnection(conn)
	assert.NoError(t, err)
}