import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';

import '../core/localization/app_localizations.dart';
import '../core/providers/connection_provider.dart';

class AddConnectionDialog extends StatefulWidget {
  const AddConnectionDialog({Key? key}) : super(key: key);

  @override
  State<AddConnectionDialog> createState() => _AddConnectionDialogState();
}

class _AddConnectionDialogState extends State<AddConnectionDialog> {
  final _formKey = GlobalKey<FormState>();
  final _urlController = TextEditingController();
  String? _errorMessage;
  bool _isProcessing = false;

  @override
  void dispose() {
    _urlController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final localizations = AppLocalizations.of(context);
    final theme = Theme.of(context);

    return AlertDialog(
      backgroundColor: theme.colorScheme.surface.withOpacity(0.9),
      title: Text(
        localizations.addConnection,
        style: TextStyle(color: Colors.white),
      ),
      content: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            TextFormField(
              controller: _urlController,
              decoration: InputDecoration(
                labelText: localizations.connectionUrl,
                hintText: localizations.pasteConnectionUrl,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
                suffixIcon: IconButton(
                  icon: const Icon(Icons.content_paste, color: Colors.white70),
                  onPressed: _pasteFromClipboard,
                  tooltip: localizations.pasteFromClipboard,
                ),
                labelStyle: TextStyle(color: Colors.white70),
                hintStyle: TextStyle(color: Colors.white38),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(8),
                  borderSide: BorderSide(color: Colors.white24),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(8),
                  borderSide: BorderSide(color: theme.colorScheme.primary),
                ),
              ),
              style: TextStyle(color: Colors.white),
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return localizations.pleaseEnterUrl;
                }
                return null;
              },
              maxLines: 3,
              minLines: 1,
            ),
            if (_errorMessage != null) ...[  
              const SizedBox(height: 16),
              Text(
                _errorMessage!,
                style: TextStyle(color: theme.colorScheme.error),
              ),
            ],
          ],
        ),
      ),
      actions: [
        TextButton(
          onPressed: () {
            Navigator.of(context).pop();
          },
          child: Text(
            localizations.cancel,
            style: TextStyle(color: Colors.white70),
          ),
        ),
        _isProcessing
            ? const CircularProgressIndicator()
            : ElevatedButton(
                onPressed: _addConnection,
                style: ElevatedButton.styleFrom(
                  backgroundColor: theme.colorScheme.primary,
                ),
                child: Text(
                  localizations.add,
                  style: TextStyle(color: Colors.white),
                ),
              ),
      ],
    );
  }

  Future<void> _pasteFromClipboard() async {
    final data = await Clipboard.getData(Clipboard.kTextPlain);
    if (data != null && data.text != null) {
      setState(() {
        _urlController.text = data.text!;
        _errorMessage = null;
      });
    }
  }

  Future<void> _addConnection() async {
    if (_formKey.currentState!.validate()) {
      setState(() {
        _isProcessing = true;
        _errorMessage = null;
      });

      try {
        final connectionProvider =
            Provider.of<ConnectionProvider>(context, listen: false);
        final success = await connectionProvider.addConnection(_urlController.text);

        if (success) {
          if (mounted) {
            Navigator.of(context).pop();
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(localizations.connectionAdded),
                backgroundColor: Theme.of(context).colorScheme.primary,
              ),
            );
          }
        } else {
          setState(() {
            _errorMessage = AppLocalizations.of(context).invalidConnectionUrl;
            _isProcessing = false;
          });
        }
      } catch (e) {
        setState(() {
          _errorMessage = e.toString();
          _isProcessing = false;
        });
      }
    }
  }
}