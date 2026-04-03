import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/theme.dart';
import 'playlist_controller.dart';
import 'playlist_detail_screen.dart';

class PlaylistsScreen extends ConsumerWidget {
  const PlaylistsScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(playlistsControllerProvider);
    ref.listen(playlistsControllerProvider, (prev, next) {
      if (next.error != null) {
        debugPrint('Playlists error: ${next.error}');
      }
    });
    return Scaffold(
      appBar: AppBar(
        title: const Text('Playlists'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () => ref.read(playlistsControllerProvider.notifier).load(),
          )
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => _showCreateDialog(context, ref),
        icon: const Icon(Icons.add),
        label: const Text('New playlist'),
      ),
      body: state.loading
          ? const Center(child: CircularProgressIndicator())
          : state.playlists.isEmpty
              ? const Center(child: Text('Create your first playlist'))
              : ListView.builder(
                  itemCount: state.playlists.length,
                  itemBuilder: (context, index) {
                    final p = state.playlists[index];
                    return Card(
                      child: ListTile(
                        title: Text(p.name, style: const TextStyle(fontWeight: FontWeight.w700)),
                        subtitle: Text(p.description.isEmpty ? 'No description' : p.description),
                        trailing: IconButton(
                          icon: const Icon(Icons.delete_outline, color: AppTheme.textSecondary),
                          onPressed: () async {
                            await ref.read(playlistsControllerProvider.notifier).delete(p.id);
                          },
                        ),
                        onTap: () => Navigator.of(context).push(
                          MaterialPageRoute(builder: (_) => PlaylistDetailScreen(playlistId: p.id, name: p.name)),
                        ),
                      ),
                    );
                  },
                ),
    );
  }

  void _showCreateDialog(BuildContext context, WidgetRef ref) {
    final name = TextEditingController();
    final desc = TextEditingController();
    showDialog(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('New playlist'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: name, decoration: const InputDecoration(labelText: 'Name')),
            const SizedBox(height: 8),
            TextField(controller: desc, decoration: const InputDecoration(labelText: 'Description')),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
          ElevatedButton(
            onPressed: () async {
              await ref.read(playlistsControllerProvider.notifier).create(name.text, desc.text);
              if (context.mounted) Navigator.pop(context);
            },
            child: const Text('Create'),
          ),
        ],
      ),
    );
  }
}
