import 'dart:async';
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class AppLocalizations {
  final Locale locale;

  AppLocalizations(this.locale);

  // Helper method to keep the code in the widgets concise
  static AppLocalizations of(BuildContext context) {
    return Localizations.of<AppLocalizations>(context, AppLocalizations)!;
  }

  // Static member to have a simple access to the delegate from the MaterialApp
  static const LocalizationsDelegate<AppLocalizations> delegate =
      _AppLocalizationsDelegate();

  late Map<String, String> _localizedStrings;

  Future<bool> load() async {
    // Load the language JSON file from the "assets/translations" folder
    String jsonString =
        await rootBundle.loadString('assets/translations/${locale.languageCode}.json');
    Map<String, dynamic> jsonMap = json.decode(jsonString);

    _localizedStrings = jsonMap.map((key, value) {
      return MapEntry(key, value.toString());
    });

    return true;
  }

  // This method will be called from every widget which needs a localized text
  String translate(String key) {
    return _localizedStrings[key] ?? key;
  }

  // Convenience method for common translations
  String get appName => translate('app_name');
  String get tagline => translate('tagline');
  String get connect => translate('connect');
  String get disconnect => translate('disconnect');
  String get addConnection => translate('add_connection');
  String get enterConnectionUrl => translate('enter_connection_url');
  String get connectionMode => translate('connection_mode');
  String get proxyMode => translate('proxy_mode');
  String get tunMode => translate('tun_mode');
  String get tapMode => translate('tap_mode');
  String get language => translate('language');
  String get english => translate('english');
  String get russian => translate('russian');
  String get speed => translate('speed');
  String get ping => translate('ping');
  String get jitter => translate('jitter');
  String get traffic => translate('traffic');
  String get received => translate('received');
  String get sent => translate('sent');
  String get noConnections => translate('no_connections');
  String get addYourFirstConnection => translate('add_your_first_connection');
  String get invalidUrl => translate('invalid_url');
  String get connectionAdded => translate('connection_added');
  String get connectionFailed => translate('connection_failed');
  String get connectionRemoved => translate('connection_removed');
  String get confirmRemoveConnection => translate('confirm_remove_connection');
  String get yes => translate('yes');
  String get no => translate('no');
  String get cancel => translate('cancel');
  String get ok => translate('ok');
  String get error => translate('error');
  String get success => translate('success');
  String get warning => translate('warning');
  String get info => translate('info');
}

class _AppLocalizationsDelegate
    extends LocalizationsDelegate<AppLocalizations> {
  const _AppLocalizationsDelegate();

  @override
  bool isSupported(Locale locale) {
    // Include all supported language codes here
    return ['en', 'ru'].contains(locale.languageCode);
  }

  @override
  Future<AppLocalizations> load(Locale locale) async {
    // AppLocalizations class is where the JSON loading actually runs
    AppLocalizations localizations = AppLocalizations(locale);
    await localizations.load();
    return localizations;
  }

  @override
  bool shouldReload(_AppLocalizationsDelegate old) => false;
}