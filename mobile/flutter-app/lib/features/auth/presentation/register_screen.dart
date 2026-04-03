import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/theme.dart';
import 'auth_controller.dart';

class RegisterScreen extends ConsumerStatefulWidget {
  const RegisterScreen({super.key});

  @override
  ConsumerState<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends ConsumerState<RegisterScreen> {
  final _email = TextEditingController();
  final _password = TextEditingController();

  @override
  void dispose() {
    _email.dispose();
    _password.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(authControllerProvider);
    ref.listen(authControllerProvider, (prev, next) {
      if (next.user != null && context.mounted) {
        context.go('/');
      }
    });

    return Scaffold(
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(20),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const SizedBox(height: 24),
              Text('Create account', style: Theme.of(context).textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.w800)),
              const SizedBox(height: 8),
              const Text('Join Freedom Music', style: TextStyle(color: AppTheme.textSecondary)),
              const SizedBox(height: 24),
              TextField(controller: _email, decoration: const InputDecoration(labelText: 'Email', prefixIcon: Icon(Icons.email_outlined))),
              const SizedBox(height: 12),
              TextField(
                controller: _password,
                obscureText: true,
                decoration: const InputDecoration(labelText: 'Password', prefixIcon: Icon(Icons.lock_outline)),
              ),
              const SizedBox(height: 20),
              if (state.error != null) Padding(padding: const EdgeInsets.only(bottom: 12), child: Text(state.error!, style: const TextStyle(color: Colors.red))),
              FilledButton(
                onPressed: state.loading
                    ? null
                    : () => ref.read(authControllerProvider.notifier).register(_email.text.trim(), _password.text.trim()),
                child: state.loading ? const SizedBox(height: 22, width: 22, child: CircularProgressIndicator(strokeWidth: 2)) : const Text('Sign up'),
              ),
              const SizedBox(height: 12),
              TextButton(onPressed: () => context.go('/login'), child: const Text('Already have an account? Login')),
            ],
          ),
        ),
      ),
    );
  }
}
