import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import '../plugin/vpn_plugin.dart';

class Connection {
  final String id;
  final String name;
  final String protocol;
  final String url;

  Connection({
    required this.id,
    required this.name,
    required this.protocol,
    required this.url,
  });

  factory Connection.fromJson(Map<String, dynamic> json) {
    return Connection(
      id: json['id'] as String,
      name: json['name'] as String,
      protocol: json['protocol'] as String,
      url: json['url'] as String,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'protocol': protocol,
      'url': url,
    };
  }
}

class ConnectionStatus {
  final bool connected;
  final bool connecting;
  final String? connectionId;

  ConnectionStatus({
    required this.connected,
    required this.connecting,
    this.connectionId,
  });

  factory ConnectionStatus.fromJson(Map<String, dynamic> json) {
    return ConnectionStatus(
      connected: json['connected'] as bool,
      connecting: json['connecting'] as bool,
      connectionId: json['connectionId'] as String?,
    );
  }
}

class ConnectionMetrics {
  final int speed; // bytes per second
  final int ping; // milliseconds
  final int jitter; // milliseconds
  final int bytesReceived;
  final int bytesSent;

  ConnectionMetrics({
    required this.speed,
    required this.ping,
    required this.jitter,
    required this.bytesReceived,
    required this.bytesSent,
  });

  factory ConnectionMetrics.fromJson(Map<String, dynamic> json) {
    return ConnectionMetrics(
      speed: json['speed'] as int,
      ping: json['ping'] as int,
      jitter: json['jitter'] as int,
      bytesReceived: json['bytesReceived'] as int,
      bytesSent: json['bytesSent'] as int,
    );
  }

  factory ConnectionMetrics.empty() {
    return ConnectionMetrics(
      speed: 0,
      ping: 0,
      jitter: 0,
      bytesReceived: 0,
      bytesSent: 0,
    );
  }
}

class ConnectionProvider extends ChangeNotifier {
  final VpnPlugin _vpnPlugin = VpnPlugin();
  
  List<Connection> _connections = [];
  ConnectionStatus _status = ConnectionStatus(connected: false, connecting: false);
  ConnectionMetrics _metrics = ConnectionMetrics.empty();
  Timer? _metricsTimer;
  StreamSubscription? _statusSubscription;
  StreamSubscription? _metricsSubscription;

  ConnectionProvider() {
    _loadConnections();
    _setupStreams();
  }

  List<Connection> get connections => _connections;
  ConnectionStatus get status => _status;
  ConnectionMetrics get metrics => _metrics;
  
  Connection? get activeConnection {
    if (_status.connectionId != null) {
      return _connections.firstWhere(
        (conn) => conn.id == _status.connectionId,
        orElse: () => throw Exception('Active connection not found'),
      );
    }
    return null;
  }

  Future<void> _loadConnections() async {
    try {
      final result = await _vpnPlugin.getConnections();
      _connections = result
          .map((json) => Connection.fromJson(json))
          .toList();
      notifyListeners();

      // Load connection status
      await _updateConnectionStatus();
    } catch (e) {
      debugPrint('Failed to load connections: $e');
    }
  }

  Future<void> _updateConnectionStatus() async {
    try {
      final result = await _vpnPlugin.getStatus();
      _status = ConnectionStatus.fromJson(result);
      notifyListeners();
    } catch (e) {
      debugPrint('Failed to update connection status: $e');
    }
  }

  Future<void> _updateConnectionMetrics() async {
    try {
      if (_status.connected) {
        final result = await _vpnPlugin.getMetrics();
        _metrics = ConnectionMetrics.fromJson(result);
        notifyListeners();
      }
    } catch (e) {
      debugPrint('Failed to update connection metrics: $e');
    }
  }

  void _setupStreams() {
    // Cancel existing subscriptions
    _statusSubscription?.cancel();
    _metricsSubscription?.cancel();
    
    // Subscribe to status updates
    _statusSubscription = _vpnPlugin.statusStream.listen((data) {
      _status = ConnectionStatus.fromJson(data);
      notifyListeners();
    });
    
    // Subscribe to metrics updates
    _metricsSubscription = _vpnPlugin.metricsStream.listen((data) {
      _metrics = ConnectionMetrics.fromJson(data);
      notifyListeners();
    });
  }

  Future<String?> addConnection(String url) async {
    try {
      final connectionId = await _vpnPlugin.addConnection(url);
      if (connectionId != null) {
        await _loadConnections(); // Reload connections
        return connectionId;
      }
      return null;
    } catch (e) {
      debugPrint('Failed to add connection: $e');
      return null;
    }
  }

  Future<bool> removeConnection(String id) async {
    try {
      final result = await _vpnPlugin.removeConnection(id);
      if (result) {
        await _loadConnections(); // Reload connections
      }
      return result;
    } catch (e) {
      debugPrint('Failed to remove connection: $e');
      return false;
    }
  }

  Future<bool> connect(String id) async {
    try {
      final result = await _vpnPlugin.connect(id);
      // Status will be updated via the status stream
      return result;
    } catch (e) {
      debugPrint('Failed to connect: $e');
      return false;
    }
  }

  Future<bool> disconnect() async {
    try {
      final result = await _vpnPlugin.disconnect();
      // Status will be updated via the status stream
      return result;
    } catch (e) {
      debugPrint('Failed to disconnect: $e');
      return false;
    }
  }
  
  Future<bool> setConnectionMode(String mode) async {
    try {
      return await _vpnPlugin.setConnectionMode(mode);
    } catch (e) {
      debugPrint('Failed to set connection mode: $e');
      return false;
    }
  }
  
  void setActiveConnection(String id) {
    // Just set the active connection without connecting
    // This is used to select a connection from the list
    notifyListeners();
  }

  @override
  void dispose() {
    _metricsTimer?.cancel();
    _statusSubscription?.cancel();
    _metricsSubscription?.cancel();
    super.dispose();
  }
}