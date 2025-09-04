import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../core/localization/app_localizations.dart';
import '../core/providers/app_provider.dart';
import '../core/providers/connection_provider.dart';
import '../core/providers/language_provider.dart';
import '../widgets/connection_list.dart';
import '../widgets/connection_metrics.dart';
import '../widgets/connection_status.dart';
import '../widgets/language_switcher.dart';
import '../widgets/marquee_text.dart';
import '../widgets/mode_selector.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final localizations = AppLocalizations.of(context);
    
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: [Color(0xFF0A0E17), Color(0xFF141824)],
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              // App Bar with title and language switcher
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      localizations.appName,
                      style: Theme.of(context).textTheme.headlineMedium?.copyWith(
                            fontWeight: FontWeight.bold,
                            color: Colors.white,
                          ),
                    ),
                    const LanguageSwitcher(),
                  ],
                ),
              ),
              
              // Tagline with marquee effect
              SizedBox(
                height: 24,
                child: MarqueeText(
                  text: localizations.tagline,
                  style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                        color: Colors.white70,
                      ),
                ),
              ),
              
              const SizedBox(height: 24),
              
              // Connection status and metrics
              Expanded(
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 16),
                  child: Column(
                    children: [
                      // Connection status
                      const ConnectionStatus(),
                      
                      const SizedBox(height: 24),
                      
                      // Connection metrics
                      const ConnectionMetrics(),
                      
                      const SizedBox(height: 24),
                      
                      // Mode selector
                      const ModeSelector(),
                      
                      const SizedBox(height: 24),
                      
                      // Connection list
                      const Expanded(
                        child: ConnectionList(),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showAddConnectionDialog(context),
        backgroundColor: Theme.of(context).colorScheme.primary,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  void _showAddConnectionDialog(BuildContext context) {
    final localizations = AppLocalizations.of(context);
    final connectionProvider = Provider.of<ConnectionProvider>(context, listen: false);
    final TextEditingController urlController = TextEditingController();

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(localizations.addConnection),
        content: TextField(
          controller: urlController,
          decoration: InputDecoration(
            hintText: localizations.enterConnectionUrl,
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text(localizations.cancel),
          ),
          ElevatedButton(
            onPressed: () async {
              final url = urlController.text.trim();
              if (url.isNotEmpty) {
                final success = await connectionProvider.addConnection(url);
                Navigator.pop(context);
                
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      success
                          ? localizations.connectionAdded
                          : localizations.invalidUrl,
                    ),
                    backgroundColor: success
                        ? Theme.of(context).colorScheme.primary
                        : Theme.of(context).colorScheme.error,
                  ),
                );
              }
            },
            child: Text(localizations.ok),
          ),
        ],
      ),
    );
  }
}