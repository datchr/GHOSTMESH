import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../core/localization/app_localizations.dart';
import '../core/providers/app_provider.dart';

class ConnectionModeSelector extends StatelessWidget {
  const ConnectionModeSelector({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final localizations = AppLocalizations.of(context);
    final appProvider = Provider.of<AppProvider>(context);
    final currentMode = appProvider.connectionMode;

    return Card(
      elevation: 8,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              localizations.connectionMode,
              style: Theme.of(context).textTheme.titleMedium,
              textAlign: TextAlign.center,
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _buildModeButton(
                    context,
                    localizations.proxy,
                    ConnectionMode.proxy,
                    currentMode,
                    appProvider,
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: _buildModeButton(
                    context,
                    localizations.tun,
                    ConnectionMode.tun,
                    currentMode,
                    appProvider,
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: _buildModeButton(
                    context,
                    localizations.tap,
                    ConnectionMode.tap,
                    currentMode,
                    appProvider,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildModeButton(
    BuildContext context,
    String label,
    ConnectionMode mode,
    ConnectionMode currentMode,
    AppProvider provider,
  ) {
    final theme = Theme.of(context);
    final isSelected = mode == currentMode;

    return ElevatedButton(
      onPressed: () {
        provider.setConnectionMode(mode);
      },
      style: ElevatedButton.styleFrom(
        backgroundColor: isSelected
            ? theme.colorScheme.primary
            : theme.colorScheme.surface,
        foregroundColor: isSelected
            ? theme.colorScheme.onPrimary
            : theme.colorScheme.onSurface,
        elevation: isSelected ? 4 : 0,
        padding: const EdgeInsets.symmetric(vertical: 12),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      ),
      child: Text(
        label,
        style: theme.textTheme.bodyMedium?.copyWith(
          color: isSelected
              ? theme.colorScheme.onPrimary
              : theme.colorScheme.onSurface,
          fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
        ),
      ),
    );
  }
}