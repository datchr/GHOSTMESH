# GHOSTMESH VPN Client

GHOSTMESH is a multi-protocol VPN client for Windows that supports various VPN protocols and provides a user-friendly interface with a cosmic-dark theme.

## Features

- **Multi-protocol Support**: VLESS, ShadowSocks, OpenVPN, WireGuard, Trojan-GFW, IKEv2/IPsec, SSTP
- **Automatic Link Parsing**: Easily add connections by pasting VPN URLs
- **Connection Management**: Connect, disconnect, and delete VPN connections
- **Connection Metrics**: View speed, ping, jitter, and traffic statistics
- **Multiple Connection Modes**: Proxy, TUN, and TAP modes
- **Multilingual Support**: English and Russian languages
- **Cosmic-Dark Theme**: Modern and sleek user interface

## Architecture

The application is built using a combination of Go and Flutter:

- **Go Backend**: Handles VPN connections, protocol implementations, and system integration
- **Flutter Frontend**: Provides a cross-platform user interface

## Project Structure

```
├── cmd/                # Command-line entry points
│   └── main.go         # Main application entry point
├── internal/           # Internal packages
│   ├── app/            # Application logic
│   ├── config/         # Configuration management
│   ├── logger/         # Logging utilities
│   └── vpn/            # VPN protocol implementations
└── flutter/           # Flutter UI code
    ├── lib/            # Dart code
    │   ├── core/       # Core functionality
    │   ├── screens/    # UI screens
    │   └── widgets/    # Reusable UI components
    └── assets/         # Static assets
```

## Supported VPN URL Formats

- **VLESS**: `vless://uuid@hostname:port?encryption=none&security=tls&sni=example.com`
- **Shadowsocks**: `ss://base64(method:password)@hostname:port`
- **OpenVPN**: `ovpn://username:password@hostname:port`
- **WireGuard**: `wg://hostname:port?privateKey=xxx&publicKey=yyy&presharedKey=zzz`
- **Trojan**: `trojan://password@hostname:port?security=tls&sni=example.com`
- **IKEv2/IPsec**: `ikev2://username:password@hostname:port?psk=xxx`
- **SSTP**: `sstp://username:password@hostname:port`

## Building and Running

### Prerequisites

- Go 1.18 or later
- Flutter 3.0 or later
- Windows 8, 10, or 11

### Build Instructions

1. Clone the repository
2. Build the Go backend: `go build -o ghostmesh.exe ./cmd/main.go`
3. Build the Flutter frontend: `cd flutter && flutter build windows`
4. Run the application: `./ghostmesh.exe`

## License

This project is licensed under the MIT License - see the LICENSE file for details.