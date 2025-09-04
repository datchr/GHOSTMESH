import 'dart:async';

import 'package:flutter/material.dart';

class MarqueeText extends StatefulWidget {
  final String text;
  final TextStyle? style;
  final Duration scrollDuration;
  final Duration pauseDuration;
  final double blankSpace;

  const MarqueeText({
    Key? key,
    required this.text,
    this.style,
    this.scrollDuration = const Duration(seconds: 20),
    this.pauseDuration = const Duration(seconds: 1),
    this.blankSpace = 50.0,
  }) : super(key: key);

  @override
  State<MarqueeText> createState() => _MarqueeTextState();
}

class _MarqueeTextState extends State<MarqueeText> {
  late ScrollController _scrollController;
  Timer? _timer;
  double _position = 0.0;
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
  int _currentTaglineIndex = 0;

  @override
  void initState() {
    super.initState();
    _scrollController = ScrollController();
    _startMarquee();
    _startTaglineRotation();
  }

  @override
  void dispose() {
    _timer?.cancel();
    _scrollController.dispose();
    super.dispose();
  }

  void _startMarquee() {
    _timer = Timer.periodic(const Duration(milliseconds: 50), (timer) {
      final maxScrollExtent = _scrollController.position.maxScrollExtent;
      _position += 1.0;
      if (_position > maxScrollExtent) {
        _position = 0.0;
      }
      _scrollController.jumpTo(_position);
    });
  }

  void _startTaglineRotation() {
    Timer.periodic(const Duration(seconds: 5), (timer) {
      setState(() {
        _currentTaglineIndex = (_currentTaglineIndex + 1) % _taglines.length;
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 24,
      width: double.infinity,
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: SingleChildScrollView(
        controller: _scrollController,
        scrollDirection: Axis.horizontal,
        physics: const NeverScrollableScrollPhysics(),
        child: Row(
          children: [
            Text(
              _taglines[_currentTaglineIndex],
              style: widget.style ?? const TextStyle(),
            ),
            SizedBox(width: widget.blankSpace),
            Text(
              _taglines[_currentTaglineIndex],
              style: widget.style ?? const TextStyle(),
            ),
          ],
        ),
      ),
    );
  }
}