import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class LanguageProvider extends ChangeNotifier {
  Locale _locale = const Locale('en', '');
  final String _prefsKey = 'language_code';

  LanguageProvider() {
    _loadSavedLanguage();
  }

  Locale get locale => _locale;

  Future<void> _loadSavedLanguage() async {
    final prefs = await SharedPreferences.getInstance();
    final savedLanguage = prefs.getString(_prefsKey);
    if (savedLanguage != null) {
      _locale = Locale(savedLanguage, '');
      notifyListeners();
    }
  }

  Future<void> setLanguage(String languageCode) async {
    if (languageCode != _locale.languageCode) {
      _locale = Locale(languageCode, '');
      final prefs = await SharedPreferences.getInstance();
      await prefs.setString(_prefsKey, languageCode);
      notifyListeners();
    }
  }

  bool isCurrentLanguage(String languageCode) {
    return _locale.languageCode == languageCode;
  }
}