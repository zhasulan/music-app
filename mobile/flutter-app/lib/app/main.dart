import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'router.dart';
import '../core/theme.dart';

void main() {
  runApp(const ProviderScope(child: FreedomMusicApp()));
}

class FreedomMusicApp extends ConsumerWidget {
  const FreedomMusicApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final router = ref.watch(routerProvider);
    return MaterialApp.router(
      title: 'Freedom Music',
      theme: AppTheme.dark(),
      routerConfig: router,
    );
  }
}
