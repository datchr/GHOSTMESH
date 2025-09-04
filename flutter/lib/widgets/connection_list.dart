import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../core/localization/app_localizations.dart';
import '../core/providers/connection_provider.dart';

class ConnectionList extends StatelessWidget {
  const ConnectionList({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final localizations = AppLocalizations.of(context);
    final connectionProvider = Provider.of<ConnectionProvider>(context);
    final connections = connectionProvider.connections;
    final activeConnectionId = connectionProvider.activeConnection?.id;

    if (connections.isEmpty) {
      return Card(
        elevation: 8,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(
                  Icons.vpn_lock,
                  size: 48,
                  color: Theme.of(context).colorScheme.onSurface.withOpacity(0.5),
                ),
                const SizedBox(height: 16),
                Text(
                  localizations.noConnectionsAdded,
                  style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                        color: Theme.of(context)
                            .colorScheme
                            .onSurface
                            .withOpacity(0.7),
                      ),
                  textAlign: TextAlign.center,
                ),
                const SizedBox(height: 8),
                Text(
                  localizations.addConnectionHint,
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: Theme.of(context)
                            .colorScheme
                            .onSurface
                            .withOpacity(0.5),
                      ),
                  textAlign: TextAlign.center,
                ),
              ],
            ),
          ),
        ),
      );
    }

    return Card(
      elevation: 8,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              localizations.connections,
              style: Theme.of(context).textTheme.titleMedium,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 16),
            ListView.separated(
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              itemCount: connections.length,
              separatorBuilder: (context, index) => const Divider(),
              itemBuilder: (context, index) {
                final connection = connections[index];
                final isActive = connection.id == activeConnectionId;

                return ListTile(
                  title: Text(
                    connection.name,
                    style: TextStyle(
                      fontWeight: isActive ? FontWeight.bold : FontWeight.normal,
                      color: isActive
                          ? Theme.of(context).colorScheme.primary
                          : Theme.of(context).colorScheme.onSurface,
                    ),
                  ),
                  subtitle: Text(
                    connection.protocol.toUpperCase(),
                    style: TextStyle(
                      color: Theme.of(context)
                          .colorScheme
                          .onSurface
                          .withOpacity(0.7),
                    ),
                  ),
                  leading: _getProtocolIcon(connection.protocol, isActive),
                  trailing: IconButton(
                    icon: const Icon(Icons.delete_outline),
                    color: Theme.of(context).colorScheme.error,
                    onPressed: () {
                      _showDeleteConfirmation(context, connection, connectionProvider);
                    },
                  ),
                  onTap: () {
                    connectionProvider.setActiveConnection(connection.id);
                  },
                );
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _getProtocolIcon(String protocol, bool isActive) {
    IconData iconData;
    
    switch (protocol.toLowerCase()) {
      case 'vless':
        iconData = Icons.security;
        break;
      case 'shadowsocks':
        iconData = Icons.shield;
        break;
      case 'openvpn':
        iconData = Icons.vpn_key;
        break;
      case 'wireguard':
        iconData = Icons.wifi_tethering;
        break;
      case 'trojan':
        iconData = Icons.security;
        break;
      case 'ikev2':
        iconData = Icons.vpn_lock;
        break;
      case 'sstp':
        iconData = Icons.https;
        break;
      default:
        iconData = Icons.vpn_key;
    }

    return CircleAvatar(
      backgroundColor: isActive
          ? Colors.blue.withOpacity(0.2)
          : Colors.grey.withOpacity(0.2),
      child: Icon(
        iconData,
        color: isActive ? Colors.blue : Colors.grey,
        size: 20,
      ),
    );
  }

  void _showDeleteConfirmation(
    BuildContext context,
    Connection connection,
    ConnectionProvider provider,
  ) {
    final localizations = AppLocalizations.of(context);
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(localizations.deleteConnection),
        content: Text(
          localizations.deleteConnectionConfirmation(connection.name),
        ),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.of(context).pop();
            },
            child: Text(localizations.cancel),
          ),
          TextButton(
            onPressed: () {
              provider.removeConnection(connection.id);
              Navigator.of(context).pop();
            },
            style: TextButton.styleFrom(
              foregroundColor: Theme.of(context).colorScheme.error,
            ),
            child: Text(localizations.delete),
          ),
        ],
      ),
    );
  }
}