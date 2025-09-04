import 'dart:async';

import 'package:flutter/material.dart';

class ScrollingTagline extends StatefulWidget {
  const ScrollingTagline({Key? key}) : super(key: key);

  @override
  State<ScrollingTagline> createState() => _ScrollingTaglineState();
}

class _ScrollingTaglineState extends State<ScrollingTagline> {
  final List<String> _taglines = [
    'Communication without compromise', // English
    'Связь без компромиссов', // Russian
    'Communication sans compromis', // French
    'Comunicación sin compromiso', // Spanish
    'Kommunikation ohne Kompromisse', // German
    '妥協のないコミュニケーション', // Japanese
    '沟通无需妥协', // Chinese
    'Comunicazione senza compromessi', // Italian
    'Comunicação sem compromisso', // Portuguese
    '타협 없는 커뮤니케이션', // Korean
  ];

  int _currentIndex = 0;
  final ScrollController _scrollController = ScrollController();
  Timer? _timer;

  @override
  void initState() {
    super.initState();
    _startScrolling();
  }

  @override
  void dispose() {
    _timer?.cancel();
    _scrollController.dispose();
    super.dispose();
  }

  void _startScrolling() {
    // Change tagline every 5 seconds
    _timer = Timer.periodic(const Duration(seconds: 5), (timer) {
      setState(() {
        _currentIndex = (_currentIndex + 1) % _taglines.length;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 30,
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface.withOpacity(0.2),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Center(
        child: AnimatedSwitcher(
          duration: const Duration(milliseconds: 500),
          transitionBuilder: (Widget child, Animation<double> animation) {
            return FadeTransition(
              opacity: animation,
              child: SlideTransition(
                position: Tween<Offset>(
                  begin: const Offset(0.0, 0.5),
                  end: Offset.zero,
                ).animate(animation),
                child: child,
              ),
            );
          },
          child: Text(
            _taglines[_currentIndex],
            key: ValueKey<int>(_currentIndex),
            style: TextStyle(
              color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
              fontStyle: FontStyle.italic,
              fontSize: 14,
            ),
            textAlign: TextAlign.center,
          ),
        ),
      ),
    );
  }
}