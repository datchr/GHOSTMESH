import 'dart:async';
import 'package:flutter/services.dart';

class VpnPlugin {
  static const MethodChannel _channel = MethodChannel('com.ghostmesh.vpnclient/vpn');
  static final VpnPlugin _instance = VpnPlugin._internal();

  // Singleton instance
  factory VpnPlugin() {
    return _instance;
  }

  VpnPlugin._internal();

  // Connection management
  Future<bool> connect(String connectionId) async {
    try {
      final result = await _channel.invokeMethod('connect', {'connectionId': connectionId});
      return result == true;
    } catch (e) {
      print('Error connecting to VPN: $e');
      return false;
    }
  }

  Future<bool> disconnect() async {
    try {
      final result = await _channel.invokeMethod('disconnect');
      return result == true;
    } catch (e) {
      print('Error disconnecting from VPN: $e');
      return false;
    }
  }

  // Connection management
  Future<String?> addConnection(String url) async {
    try {
      final result = await _channel.invokeMethod('addConnection', {'url': url});
      if (result is Map && result['success'] == true) {
        return result['id'] as String;
      }
      return null;
    } catch (e) {
      print('Error adding connection: $e');
      return null;
    }
  }

  Future<bool> removeConnection(String connectionId) async {
    try {
      final result = await _channel.invokeMethod('removeConnection', {'connectionId': connectionId});
      return result == true;
    } catch (e) {
      print('Error removing connection: $e');
      return false;
    }
  }

  Future<List<Map<String, dynamic>>> getConnections() async {
    try {
      final result = await _channel.invokeMethod('getConnections');
      return List<Map<String, dynamic>>.from(result);
    } catch (e) {
      print('Error getting connections: $e');
      return [];
    }
  }

  // Status and metrics
  Future<Map<String, dynamic>> getStatus() async {
    try {
      final result = await _channel.invokeMethod('getStatus');
      return Map<String, dynamic>.from(result);
    } catch (e) {
      print('Error getting status: $e');
      return {
        'connected': false,
        'connecting': false,
        'connectionId': null,
      };
    }
  }

  Future<Map<String, dynamic>> getMetrics() async {
    try {
      final result = await _channel.invokeMethod('getMetrics');
      return Map<String, dynamic>.from(result);
    } catch (e) {
      print('Error getting metrics: $e');
      return {
        'ping': 0,
        'jitter': 0,
        'speed': 0,
        'bytesReceived': 0,
        'bytesSent': 0,
      };
    }
  }

  // Settings
  Future<bool> setConnectionMode(String mode) async {
    try {
      final result = await _channel.invokeMethod('setConnectionMode', {'mode': mode});
      return result == true;
    } catch (e) {
      print('Error setting connection mode: $e');
      return false;
    }
  }

  // Stream for status updates
  Stream<Map<String, dynamic>> get statusStream {
    return const EventChannel('com.ghostmesh.vpnclient/status')
        .receiveBroadcastStream()
        .map((event) => Map<String, dynamic>.from(event));
  }

  // Stream for metrics updates
  Stream<Map<String, dynamic>> get metricsStream {
    return const EventChannel('com.ghostmesh.vpnclient/metrics')
        .receiveBroadcastStream()
        .map((event) => Map<String, dynamic>.from(event));
  }
}