package vpn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		protocol string
		valid    bool
	}{
		{
			name:     "Valid VLESS URL",
			url:      "vless://8a45f2f0-5361-4e2b-b330-1d756c448c59@example.com:443?encryption=none&security=tls&sni=example.com",
			protocol: "vless",
			valid:    true,
		},
		{
			name:     "Valid Shadowsocks URL",
			url:      "ss://YWVzLTI1Ni1nY206cGFzc3dvcmQ=@example.com:8388",
			protocol: "shadowsocks",
			valid:    true,
		},
		{
			name:     "Valid OpenVPN URL",
			url:      "ovpn://username:password@example.com:1194",
			protocol: "openvpn",
			valid:    true,
		},
		{
			name:     "Valid WireGuard URL",
			url:      "wg://example.com:51820?privateKey=abc123&publicKey=def456",
			protocol: "wireguard",
			valid:    true,
		},
		{
			name:     "Valid Trojan URL",
			url:      "trojan://password@example.com:443?security=tls&sni=example.com",
			protocol: "trojan",
			valid:    true,
		},
		{
			name:     "Valid IKEv2 URL",
			url:      "ikev2://username:password@example.com:500?psk=shared-secret",
			protocol: "ikev2",
			valid:    true,
		},
		{
			name:     "Valid SSTP URL",
			url:      "sstp://username:password@example.com:443",
			protocol: "sstp",
			valid:    true,
		},
		{
			name:     "Invalid URL",
			url:      "invalid-url",
			protocol: "",
			valid:    false,
		},
		{
			name:     "Unsupported Protocol",
			url:      "unknown://example.com:443",
			protocol: "",
			valid:    false,
		},
	}

	parser := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn, err := parser.Parse(tt.url)
			if tt.valid {
				assert.NoError(t, err)
				assert.NotNil(t, conn)
				assert.Equal(t, tt.protocol, conn.Protocol)
			} else {
				assert.Error(t, err)
			}
		})
	}
}