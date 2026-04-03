import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/theme.dart';
import '../../auth/presentation/auth_controller.dart';
import '../data/profile_repository.dart';

final profileFutureProvider = FutureProvider((ref) async {
  return ref.read(profileRepositoryProvider).me();
});

class ProfileScreen extends ConsumerStatefulWidget {
  const ProfileScreen({super.key});

  @override
  ConsumerState<ProfileScreen> createState() => _ProfileScreenState();
}

class _ProfileScreenState extends ConsumerState<ProfileScreen> {
  final _name = TextEditingController();
  final _country = TextEditingController();

  @override
  void dispose() {
    _name.dispose();
    _country.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final authUser = ref.watch(authControllerProvider).user;
    final profileAsync = ref.watch(profileFutureProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Profile')),
      body: profileAsync.when(
        data: (profile) {
          _name.text = profile.name;
          _country.text = profile.country;
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(profile.email, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w700)),
                      const SizedBox(height: 6),
                      Text('User ID: ${authUser?.id ?? '-'}', style: const TextStyle(color: AppTheme.textSecondary)),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 12),
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    children: [
                      TextField(controller: _name, decoration: const InputDecoration(labelText: 'Name')),
                      const SizedBox(height: 8),
                      TextField(controller: _country, decoration: const InputDecoration(labelText: 'Country')),
                      const SizedBox(height: 12),
                      FilledButton(
                        onPressed: () async {
                          await ref.read(profileRepositoryProvider).update(name: _name.text, country: _country.text);
                          ref.refresh(profileFutureProvider);
                          if (mounted) {
                            ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Profile updated')));
                          }
                        },
                        child: const Text('Save changes'),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Failed to load profile: $e')),
      ),
    );
  }
}
