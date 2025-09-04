import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../core/localization/app_localizations.dart';
import '../core/providers/connection_provider.dart';

class ConnectionStatus extends StatelessWidget {
  const ConnectionStatus({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final localizations = AppLocalizations.of(context);
    final connectionProvider = Provider.of<ConnectionProvider>(context);
    final status = connectionProvider.status;
    final activeConnection = connectionProvider.activeConnection;

    return Card(
      elevation: 8,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            // Connection name
            if (activeConnection != null)
              Text(
                activeConnection.name,
                style: Theme.of(context).textTheme.titleLarge,
                textAlign: TextAlign.center,
              )
            else
              Text(
                localizations.noConnections,
                style: Theme.of(context).textTheme.titleLarge,
                textAlign: TextAlign.center,
              ),

            const SizedBox(height: 16),

            // Connection protocol
            if (activeConnection != null)
              Text(
                activeConnection.protocol.toUpperCase(),
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Theme.of(context).colorScheme.secondary,
                    ),
                textAlign: TextAlign.center,
              ),

            const SizedBox(height: 24),

            // Connect/Disconnect button
            _buildConnectionButton(context, status, connectionProvider),
          ],
        ),
      ),
    );
  }

  Widget _buildConnectionButton(
    BuildContext context,
    ConnectionStatus status,
    ConnectionProvider provider,
  ) {
    final localizations = AppLocalizations.of(context);
    final theme = Theme.of(context);

    if (status.connecting) {
      // Show loading indicator while connecting
      return ElevatedButton(
        onPressed: null,
        style: ElevatedButton.styleFrom(
          backgroundColor: theme.colorScheme.primary.withOpacity(0.7),
          padding: const EdgeInsets.symmetric(vertical: 16),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        ),
        child: const SizedBox(
          height: 24,
          width: 24,
          child: CircularProgressIndicator(
            color: Colors.white,
            strokeWidth: 2,
          ),
        ),
      );
    } else if (status.connected) {
      // Show disconnect button when connected
      return ElevatedButton(
        onPressed: () async {
          await provider.disconnect();
        },
        style: ElevatedButton.styleFrom(
          backgroundColor: theme.colorScheme.error,
          padding: const EdgeInsets.symmetric(vertical: 16),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        ),
        child: Text(
          localizations.disconnect,
          style: theme.textTheme.titleMedium?.copyWith(
            color: Colors.white,
            fontWeight: FontWeight.bold,
          ),
        ),
      );
    } else {
      // Show connect button when disconnected
      final activeConnection = provider.activeConnection;
      final hasConnections = provider.connections.isNotEmpty;

      return ElevatedButton(
        onPressed: hasConnections
            ? () async {
                if (activeConnection != null) {
                  await provider.connect(activeConnection.id);
                } else if (provider.connections.isNotEmpty) {
                  await provider.connect(provider.connections.first.id);
                }
              }
            : null,
        style: ElevatedButton.styleFrom(
          backgroundColor: theme.colorScheme.primary,
          padding: const EdgeInsets.symmetric(vertical: 16),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        ),
        child: Text(
          localizations.connect,
          style: theme.textTheme.titleMedium?.copyWith(
            color: Colors.white,
            fontWeight: FontWeight.bold,
          ),
        ),
      );
    }
  }
}