package vpn

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/ghostmesh/vpnclient/internal/config"
	"github.com/ghostmesh/vpnclient/internal/logger"
	"github.com/pkg/errors"
)

// Parser is responsible for parsing VPN connection URLs
type Parser struct {
	log *logger.Logger
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{
		log: logger.GetLogger(),
	}
}
	return &Parser{
		log: logger.GetLogger(),
	}
}

// Parse parses a VPN connection URL and returns a Connection object
func (p *Parser) Parse(rawURL string) (*config.Connection, error) {
	p.log.Debug("Parsing URL: %s", rawURL)

	// Determine protocol from URL scheme
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse URL")
	}

	// Extract protocol from scheme
	protocol := parsedURL.Scheme
	if protocol == "" {
		return nil, errors.New("URL scheme (protocol) is missing")
	}

	// Parse based on protocol
	switch strings.ToLower(protocol) {
	case "vless":
		return p.parseVLESS(parsedURL)
	case "ss", "shadowsocks":
		return p.parseShadowsocks(parsedURL)
	case "ovpn", "openvpn":
		return p.parseOpenVPN(parsedURL)
	case "wg", "wireguard":
		return p.parseWireGuard(parsedURL)
	case "trojan":
		return p.parseTrojan(parsedURL)
	case "ikev2", "ipsec":
		return p.parseIKEv2(parsedURL)
	case "sstp":
		return p.parseSSTP(parsedURL)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", protocol)
	}
}

// parseVLESS parses a VLESS URL
func (p *Parser) parseVLESS(parsedURL *url.URL) (*config.Connection, error) {
	// VLESS URL format: vless://uuid@hostname:port?encryption=none&type=tcp&security=tls&sni=example.com
	
	// Extract user info (UUID)
	id := parsedURL.User.Username()
	if id == "" {
		return nil, errors.New("VLESS UUID is missing")
	}

	// Extract host and port
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) != 2 {
		return nil, errors.New("VLESS host or port is missing")
	}

	host := hostParts[0]
	port, err := strconv.Atoi(hostParts[1])
	if err != nil {
		return nil, errors.Wrap(err, "invalid VLESS port")
	}

	// Parse query parameters
	query := parsedURL.Query()
	security := query.Get("security")
	serverName := query.Get("sni")
	fingerprint := query.Get("fp")
	publicKey := query.Get("pbk")
	shortID := query.Get("sid")

	// Create VLESS config
	config := VLESSConfig{
		ID:        id,
		Address:   host,
		Port:      port,
		TLS:       security == "tls",
		Reality:   security == "reality",
		Vision:    query.Get("flow") == "xtls-rprx-vision",
		ServerName: serverName,
		Fingerprint: fingerprint,
		PublicKey: publicKey,
		ShortID:   shortID,
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal VLESS config")
	}

	// Create connection
	return &config.Connection{
		ID:       generateID(),
		Name:     fmt.Sprintf("VLESS - %s", host),
		Protocol: "vless",
		URL:      rawURL,
		Config:   string(configJSON),
	}, nil
}

// parseShadowsocks parses a Shadowsocks URL
func (p *Parser) parseShadowsocks(parsedURL *url.URL) (*config.Connection, error) {
	// Shadowsocks URL format: ss://base64(method:password)@hostname:port
	
	// Extract user info (method:password)
	userInfo := parsedURL.User.String()
	if userInfo == "" {
		return nil, errors.New("Shadowsocks user info is missing")
	}

	// Decode base64 if needed
	if !strings.Contains(userInfo, ":") {
		decodedBytes, err := base64.StdEncoding.DecodeString(userInfo)
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode Shadowsocks user info")
		}
		userInfo = string(decodedBytes)
	}

	// Split method and password
	parts := strings.SplitN(userInfo, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid Shadowsocks user info format")
	}

	method := parts[0]
	password := parts[1]

	// Extract host and port
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) != 2 {
		return nil, errors.New("Shadowsocks host or port is missing")
	}

	host := hostParts[0]
	port, err := strconv.Atoi(hostParts[1])
	if err != nil {
		return nil, errors.Wrap(err, "invalid Shadowsocks port")
	}

	// Create Shadowsocks config
	config := struct {
		Method   string `json:"method"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Port     int    `json:"port"`
	}{
		Method:   method,
		Password: password,
		Address:  host,
		Port:     port,
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal Shadowsocks config")
	}

	// Create connection
	return &config.Connection{
		ID:       generateID(),
		Name:     fmt.Sprintf("Shadowsocks - %s", host),
		Protocol: "shadowsocks",
		URL:      parsedURL.String(),
		Config:   string(configJSON),
	}, nil
}

// parseOpenVPN parses an OpenVPN URL
func (p *Parser) parseOpenVPN(parsedURL *url.URL) (*config.Connection, error) {
	// OpenVPN URL format: ovpn://username:password@hostname:port
	
	// Extract user info
	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()

	// Extract host and port
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) != 2 {
		return nil, errors.New("OpenVPN host or port is missing")
	}

	host := hostParts[0]
	port, err := strconv.Atoi(hostParts[1])
	if err != nil {
		return nil, errors.Wrap(err, "invalid OpenVPN port")
	}

	// Create OpenVPN config
	config := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Port     int    `json:"port"`
	}{
		Username: username,
		Password: password,
		Address:  host,
		Port:     port,
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal OpenVPN config")
	}

	// Create connection
	return &config.Connection{
		ID:       generateID(),
		Name:     fmt.Sprintf("OpenVPN - %s", host),
		Protocol: "openvpn",
		URL:      parsedURL.String(),
		Config:   string(configJSON),
	}, nil
}

// parseWireGuard parses a WireGuard URL
func (p *Parser) parseWireGuard(parsedURL *url.URL) (*config.Connection, error) {
	// WireGuard URL format: wg://hostname:port?privateKey=xxx&publicKey=yyy&presharedKey=zzz
	
	// Extract host and port
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) != 2 {
		return nil, errors.New("WireGuard host or port is missing")
	}

	host := hostParts[0]
	port, err := strconv.Atoi(hostParts[1])
	if err != nil {
		return nil, errors.Wrap(err, "invalid WireGuard port")
	}

	// Parse query parameters
	query := parsedURL.Query()
	privateKey := query.Get("privateKey")
	publicKey := query.Get("publicKey")
	presharedKey := query.Get("presharedKey")

	// Create WireGuard config
	config := struct {
		PrivateKey   string `json:"private_key"`
		PublicKey    string `json:"public_key"`
		PresharedKey string `json:"preshared_key"`
		Address      string `json:"address"`
		Port         int    `json:"port"`
	}{
		PrivateKey:   privateKey,
		PublicKey:    publicKey,
		PresharedKey: presharedKey,
		Address:      host,
		Port:         port,
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal WireGuard config")
	}

	// Create connection
	return &config.Connection{
		ID:       generateID(),
		Name:     fmt.Sprintf("WireGuard - %s", host),
		Protocol: "wireguard",
		URL:      parsedURL.String(),
		Config:   string(configJSON),
	}, nil
}

// parseTrojan parses a Trojan URL
func (p *Parser) parseTrojan(parsedURL *url.URL) (*config.Connection, error) {
	// Trojan URL format: trojan://password@hostname:port?security=tls&sni=example.com
	
	// Extract password
	password := parsedURL.User.Username()
	if password == "" {
		return nil, errors.New("Trojan password is missing")
	}

	// Extract host and port
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) != 2 {
		return nil, errors.New("Trojan host or port is missing")
	}

	host := hostParts[0]
	port, err := strconv.Atoi(hostParts[1])
	if err != nil {
		return nil, errors.Wrap(err, "invalid Trojan port")
	}

	// Parse query parameters
	query := parsedURL.Query()
	security := query.Get("security")
	serverName := query.Get("sni")

	// Create Trojan config
	config := struct {
		Password   string `json:"password"`
		Address    string `json:"address"`
		Port       int    `json:"port"`
		TLS        bool   `json:"tls"`
		ServerName string `json:"server_name"`
	}{
		Password:   password,
		Address:    host,
		Port:       port,
		TLS:        security == "tls" || security == "", // Default to TLS
		ServerName: serverName,
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal Trojan config")
	}

	// Create connection
	return &config.Connection{
		ID:       generateID(),
		Name:     fmt.Sprintf("Trojan - %s", host),
		Protocol: "trojan",
		URL:      parsedURL.String(),
		Config:   string(configJSON),
	}, nil
}

// parseIKEv2 parses an IKEv2/IPsec URL
func (p *Parser) parseIKEv2(parsedURL *url.URL) (*config.Connection, error) {
	// IKEv2 URL format: ikev2://username:password@hostname:port?psk=xxx
	
	// Extract user info
	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()

	// Extract host and port
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) < 1 {
		return nil, errors.New("IKEv2 host is missing")
	}

	host := hostParts[0]
	port := 500 // Default IKEv2 port
	if len(hostParts) > 1 {
		var err error
		port, err = strconv.Atoi(hostParts[1])
		if err != nil {
			return nil, errors.Wrap(err, "invalid IKEv2 port")
		}
	}

	// Parse query parameters
	query := parsedURL.Query()
	psk := query.Get("psk")

	// Create IKEv2 config
	config := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Port     int    `json:"port"`
		PSK      string `json:"psk"`
	}{
		Username: username,
		Password: password,
		Address:  host,
		Port:     port,
		PSK:      psk,
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal IKEv2 config")
	}

	// Create connection
	return &config.Connection{
		ID:       generateID(),
		Name:     fmt.Sprintf("IKEv2 - %s", host),
		Protocol: "ikev2",
		URL:      parsedURL.String(),
		Config:   string(configJSON),
	}, nil
}

// parseSSTP parses an SSTP URL
func (p *Parser) parseSSTP(parsedURL *url.URL) (*config.Connection, error) {
	// SSTP URL format: sstp://username:password@hostname:port
	
	// Extract user info
	username := parsedURL.User.Username()
	password, _ := parsedURL.User.Password()

	// Extract host and port
	hostParts := strings.Split(parsedURL.Host, ":")
	if len(hostParts) < 1 {
		return nil, errors.New("SSTP host is missing")
	}

	host := hostParts[0]
	port := 443 // Default SSTP port
	if len(hostParts) > 1 {
		var err error
		port, err = strconv.Atoi(hostParts[1])
		if err != nil {
			return nil, errors.Wrap(err, "invalid SSTP port")
		}
	}

	// Create SSTP config
	config := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Address  string `json:"address"`
		Port     int    `json:"port"`
	}{
		Username: username,
		Password: password,
		Address:  host,
		Port:     port,
	}

	// Convert config to JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal SSTP config")
	}

	// Create connection
	return &config.Connection{
		ID:       generateID(),
		Name:     fmt.Sprintf("SSTP - %s", host),
		Protocol: "sstp",
		URL:      parsedURL.String(),
		Config:   string(configJSON),
	}, nil
}

// generateID generates a unique ID for a connection
func generateID() string {
	// Generate a random UUID-like string
	return fmt.Sprintf("%d", time.Now().UnixNano())
}