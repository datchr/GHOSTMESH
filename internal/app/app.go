package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ghostmesh/vpnclient/internal/config"
	"github.com/ghostmesh/vpnclient/internal/logger"
	"github.com/ghostmesh/vpnclient/internal/vpn"
	"github.com/go-flutter-desktop/go-flutter"
	"github.com/pkg/errors"
)

// App represents the main application
type App struct {
	config *config.Config
	logger *logger.Logger
	vpnManager *vpn.Manager
}

// NewApp creates a new application instance
func NewApp(cfg *config.Config) (*App, error) {
	return &App{
		config: cfg,
		logger: logger.GetLogger(),
		vpnManager: vpn.NewManager(),
	}, nil
}

// Run starts the application
func (a *App) Run() error {
	a.logger.Info("Starting GHOSTMESH VPN Client")

	// Initialize VPN manager
	if err := a.vpnManager.Init(); err != nil {
		return errors.Wrap(err, "failed to initialize VPN manager")
	}

	// Start Flutter engine
	return a.startFlutter()
}

// startFlutter starts the Flutter engine
func (a *App) startFlutter() error {
	// Determine the path to the Flutter assets
	execPath, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "failed to get executable path")
	}

	execDir := filepath.Dir(execPath)
	assetPath := filepath.Join(execDir, "flutter_assets")

	// Check if we're in development mode
	_, currentFile, _, ok := runtime.Caller(0)
	if ok {
		// In development mode, use the build directory
		projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentFile)))
		assetPath = filepath.Join(projectRoot, "build", "flutter_assets")
	}

	// Configure Flutter options
	options := []flutter.Option{
		flutter.WindowInitialDimensions(800, 600),
		flutter.WindowTitle("GHOSTMESH - Communication without compromise"),
		flutter.AddPlugin(&vpnPlugin{app: a}),
	}

	// Start Flutter engine
	return flutter.Run(assetPath, options...)
}

// vpnPlugin implements the Flutter plugin interface for VPN functionality
type vpnPlugin struct {
	app *App
}

// InitPlugin initializes the plugin
func (p *vpnPlugin) InitPlugin(messenger flutter.BinaryMessenger) error {
	channel := flutter.NewMethodChannel(messenger, "ghostmesh.vpn/methods")
	channel.HandleFunc("connect", p.handleConnect)
	channel.HandleFunc("disconnect", p.handleDisconnect)
	channel.HandleFunc("addConnection", p.handleAddConnection)
	channel.HandleFunc("removeConnection", p.handleRemoveConnection)
	channel.HandleFunc("getConnections", p.handleGetConnections)
	channel.HandleFunc("getConnectionStatus", p.handleGetConnectionStatus)
	channel.HandleFunc("getConnectionMetrics", p.handleGetConnectionMetrics)
	channel.HandleFunc("setConnectionMode", p.handleSetConnectionMode)
	channel.HandleFunc("setLanguage", p.handleSetLanguage)
	return nil
}

// handleConnect handles the connect method call
func (p *vpnPlugin) handleConnect(arguments interface{}) (reply interface{}, err error) {
	args, ok := arguments.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("invalid arguments")
	}

	connID, ok := args["id"].(string)
	if !ok {
		return nil, errors.New("invalid connection ID")
	}

	// Get connection from config
	conn := p.app.config.GetConnection(connID)
	if conn == nil {
		return nil, errors.New("connection not found")
	}

	// Connect to VPN
	if err := p.app.vpnManager.Connect(conn); err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}

	return map[string]interface{}{"success": true}, nil
}

// handleDisconnect handles the disconnect method call
func (p *vpnPlugin) handleDisconnect(arguments interface{}) (reply interface{}, err error) {
	// Disconnect from VPN
	if err := p.app.vpnManager.Disconnect(); err != nil {
		return nil, errors.Wrap(err, "failed to disconnect")
	}

	return map[string]interface{}{"success": true}, nil
}

// handleAddConnection handles the addConnection method call
func (p *vpnPlugin) handleAddConnection(arguments interface{}) (reply interface{}, err error) {
	args, ok := arguments.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("invalid arguments")
	}
	
	// Get URL from arguments
	urlStr, ok := args["url"].(string)
	if !ok || urlStr == "" {
		return nil, errors.New("URL is required")
	}
	
	// Parse URL using the link parser
	parser := vpn.NewParser()
	conn, err := parser.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse URL")
	}
	
	// Add connection to the VPN manager
	if err := p.app.vpnManager.AddConnection(conn); err != nil {
		return nil, errors.Wrap(err, "failed to add connection")
	}
	
	// Add connection to config
	p.app.config.AddConnection(*conn)

	// Save config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get home directory")
	}
	configPath := filepath.Join(homeDir, ".ghostmesh", "config.json")
	if err := p.app.config.Save(configPath); err != nil {
		return nil, errors.Wrap(err, "failed to save config")
	}

	return map[string]interface{}{
		"success": true,
		"id": conn.ID,
	}, nil
}

// handleRemoveConnection handles the removeConnection method call
func (p *vpnPlugin) handleRemoveConnection(arguments interface{}) (reply interface{}, err error) {
	args, ok := arguments.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("invalid arguments")
	}

	connID, ok := args["id"].(string)
	if !ok {
		return nil, errors.New("invalid connection ID")
	}

	// Remove connection from config
	p.app.config.RemoveConnection(connID)

	// Save config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get home directory")
	}
	configPath := filepath.Join(homeDir, ".ghostmesh", "config.json")
	if err := p.app.config.Save(configPath); err != nil {
		return nil, errors.Wrap(err, "failed to save config")
	}

	return map[string]interface{}{"success": true}, nil
}

// handleGetConnections handles the getConnections method call
func (p *vpnPlugin) handleGetConnections(arguments interface{}) (reply interface{}, err error) {
	conns := make([]map[string]interface{}, 0, len(p.app.config.Connections))
	for _, conn := range p.app.config.Connections {
		conns = append(conns, map[string]interface{}{
			"id":       conn.ID,
			"name":     conn.Name,
			"protocol": conn.Protocol,
			"url":      conn.URL,
		})
	}

	return map[string]interface{}{"connections": conns}, nil
}

// handleGetConnectionStatus handles the getConnectionStatus method call
func (p *vpnPlugin) handleGetConnectionStatus(arguments interface{}) (reply interface{}, err error) {
	status := p.app.vpnManager.GetStatus()

	return map[string]interface{}{
		"connected":    status.Connected,
		"connecting":   status.Connecting,
		"connectionId": status.ConnectionID,
	}, nil
}

// handleGetConnectionMetrics handles the getConnectionMetrics method call
func (p *vpnPlugin) handleGetConnectionMetrics(arguments interface{}) (reply interface{}, err error) {
	metrics := p.app.vpnManager.GetMetrics()

	return map[string]interface{}{
		"speed":        metrics.Speed,
		"ping":         metrics.Ping,
		"jitter":       metrics.Jitter,
		"bytesReceived": metrics.BytesReceived,
		"bytesSent":    metrics.BytesSent,
	}, nil
}

// handleSetConnectionMode handles the setConnectionMode method call
func (p *vpnPlugin) handleSetConnectionMode(arguments interface{}) (reply interface{}, err error) {
	args, ok := arguments.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("invalid arguments")
	}

	mode, ok := args["mode"].(string)
	if !ok {
		return nil, errors.New("invalid mode")
	}

	// Validate mode
	if mode != "proxy" && mode != "tun" && mode != "tap" {
		return nil, errors.New("invalid mode: must be proxy, tun, or tap")
	}

	// Update config
	p.app.config.ConnectionMode = mode

	// Save config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get home directory")
	}
	configPath := filepath.Join(homeDir, ".ghostmesh", "config.json")
	if err := p.app.config.Save(configPath); err != nil {
		return nil, errors.Wrap(err, "failed to save config")
	}

	// Update VPN manager
	if err := p.app.vpnManager.SetConnectionMode(mode); err != nil {
		return nil, errors.Wrap(err, "failed to set connection mode")
	}

	return map[string]interface{}{"success": true}, nil
}

// handleSetLanguage handles the setLanguage method call
func (p *vpnPlugin) handleSetLanguage(arguments interface{}) (reply interface{}, err error) {
	args, ok := arguments.(map[interface{}]interface{})
	if !ok {
		return nil, errors.New("invalid arguments")
	}

	lang, ok := args["language"].(string)
	if !ok {
		return nil, errors.New("invalid language")
	}

	// Validate language
	if lang != "en" && lang != "ru" {
		return nil, errors.New("invalid language: must be en or ru")
	}

	// Update config
	p.app.config.Language = lang

	// Save config
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get home directory")
	}
	configPath := filepath.Join(homeDir, ".ghostmesh", "config.json")
	if err := p.app.config.Save(configPath); err != nil {
		return nil, errors.Wrap(err, "failed to save config")
	}

	return map[string]interface{}{"success": true}, nil
}