import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'dart:async';

import '../features/auth/presentation/auth_controller.dart';
import '../features/auth/presentation/login_screen.dart';
import '../features/auth/presentation/register_screen.dart';
import '../features/auth/presentation/splash_screen.dart';
import 'home_screen.dart';
import 'home_shell.dart';

class RouterRefreshStream extends ChangeNotifier {
  RouterRefreshStream(Stream<dynamic> stream) {
    _sub = stream.listen((_) => notifyListeners());
  }
  late final StreamSubscription _sub;
  @override
  void dispose() {
    _sub.cancel();
    super.dispose();
  }
}

final routerProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: '/splash',
    refreshListenable: RouterRefreshStream(
      ref.watch(authControllerProvider.notifier).stream,
    ),
    routes: [
      GoRoute(path: '/splash', builder: (context, state) => const SplashScreen()),
      GoRoute(path: '/login', builder: (context, state) => const LoginScreen()),
      GoRoute(path: '/register', builder: (context, state) => const RegisterScreen()),
      GoRoute(path: '/', builder: (context, state) => const HomeScreen()),
    ],
    redirect: (context, state) {
      final authState = ref.read(authControllerProvider);
      final path = state.uri.path;
      final loggingIn = path == '/login' || path == '/register';
      final splash = path == '/splash';
      if (!authState.authenticated) {
        if (splash) return null; // splash decides
        if (!loggingIn) return '/login';
      } else {
        if (loggingIn || splash) return '/';
      }
      return null;
    },
  );
});
