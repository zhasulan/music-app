import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import 'auth_controller.dart';

class SplashScreen extends ConsumerWidget {
  const SplashScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.listen(authControllerProvider, (prev, next) {
      if (next.user != null && context.mounted) {
        context.go('/');
      }
    });

    ref.read(authControllerProvider.notifier).loadFromStorage().then((_) {
      final auth = ref.read(authControllerProvider);
      if (context.mounted) {
        if (auth.authenticated) {
          context.go('/');
        } else {
          context.go('/login');
        }
      }
    });

    return const Scaffold(
      body: Center(child: CircularProgressIndicator()),
    );
  }
}
