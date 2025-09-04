import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../core/localization/app_localizations.dart';
import '../core/providers/connection_provider.dart';

class ConnectionMetrics extends StatelessWidget {
  const ConnectionMetrics({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final localizations = AppLocalizations.of(context);
    final connectionProvider = Provider.of<ConnectionProvider>(context);
    final metrics = connectionProvider.metrics;
    final status = connectionProvider.status;

    return Card(
      elevation: 8,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Speed and Ping
            Row(
              children: [
                Expanded(
                  child: _buildMetricItem(
                    context,
                    localizations.speed,
                    _formatSpeed(metrics.speed),
                    Icons.speed,
                    status.connected,
                  ),
                ),
                Expanded(
                  child: _buildMetricItem(
                    context,
                    localizations.ping,
                    '${metrics.ping} ms',
                    Icons.network_check,
                    status.connected,
                  ),
                ),
              ],
            ),

            const SizedBox(height: 16),

            // Jitter and Traffic
            Row(
              children: [
                Expanded(
                  child: _buildMetricItem(
                    context,
                    localizations.jitter,
                    '${metrics.jitter} ms',
                    Icons.compare_arrows,
                    status.connected,
                  ),
                ),
                Expanded(
                  child: _buildMetricItem(
                    context,
                    localizations.traffic,
                    _formatTraffic(metrics.bytesReceived, metrics.bytesSent),
                    Icons.swap_vert,
                    status.connected,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMetricItem(
    BuildContext context,
    String title,
    String value,
    IconData icon,
    bool isConnected,
  ) {
    final theme = Theme.of(context);
    final color = isConnected
        ? theme.colorScheme.primary
        : theme.colorScheme.onSurface.withOpacity(0.5);

    return Column(
      children: [
        Icon(
          icon,
          color: color,
          size: 24,
        ),
        const SizedBox(height: 8),
        Text(
          title,
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.onSurface.withOpacity(0.7),
          ),
        ),
        const SizedBox(height: 4),
        Text(
          isConnected ? value : '-',
          style: theme.textTheme.titleMedium?.copyWith(
            color: isConnected
                ? theme.colorScheme.onSurface
                : theme.colorScheme.onSurface.withOpacity(0.5),
            fontWeight: FontWeight.bold,
          ),
        ),
      ],
    );
  }

  String _formatSpeed(int bytesPerSecond) {
    if (bytesPerSecond < 1024) {
      return '$bytesPerSecond B/s';
    } else if (bytesPerSecond < 1024 * 1024) {
      final kbps = bytesPerSecond / 1024;
      return '${kbps.toStringAsFixed(1)} KB/s';
    } else {
      final mbps = bytesPerSecond / (1024 * 1024);
      return '${mbps.toStringAsFixed(1)} MB/s';
    }
  }

  String _formatTraffic(int bytesReceived, int bytesSent) {
    final received = _formatBytes(bytesReceived);
    final sent = _formatBytes(bytesSent);
    return '↓$received ↑$sent';
  }

  String _formatBytes(int bytes) {
    if (bytes < 1024) {
      return '$bytes B';
    } else if (bytes < 1024 * 1024) {
      final kb = bytes / 1024;
      return '${kb.toStringAsFixed(1)} KB';
    } else if (bytes < 1024 * 1024 * 1024) {
      final mb = bytes / (1024 * 1024);
      return '${mb.toStringAsFixed(1)} MB';
    } else {
      final gb = bytes / (1024 * 1024 * 1024);
      return '${gb.toStringAsFixed(1)} GB';
    }
  }
}