import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../plugin/vpn_plugin.dart';

class AppProvider extends ChangeNotifier {
  String _connectionMode = 'proxy'; // proxy, tun, tap
  final String _connectionModeKey = 'connection_mode';
  final VpnPlugin _vpnPlugin = VpnPlugin();

  AppProvider() {
    _loadSettings();
  }

  String get connectionMode => _connectionMode;

  Future<void> _loadSettings() async {
    final prefs = await SharedPreferences.getInstance();
    final savedMode = prefs.getString(_connectionModeKey);
    if (savedMode != null) {
      _connectionMode = savedMode;
      notifyListeners();
    }
  }

  Future<void> setConnectionMode(String mode) async {
    if (mode != _connectionMode) {
      // Save to preferences
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString(_connectionModeKey, mode);
      
      // Update connection mode via plugin
      final success = await _vpnPlugin.setConnectionMode(mode);
      
      if (success) {
        _connectionMode = mode;
        notifyListeners();
      }
    }
  }
}