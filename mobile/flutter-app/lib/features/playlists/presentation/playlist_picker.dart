import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/theme.dart';
import '../data/playlist_repository.dart';
import 'playlist_controller.dart';

Future<void> showPlaylistPicker(BuildContext context, WidgetRef ref, {required String trackId}) async {
  // ensure playlists are loaded
  await ref.read(playlistsControllerProvider.notifier).load();
  final playlists = ref.read(playlistsControllerProvider).playlists;
  showModalBottomSheet(
    context: context,
    backgroundColor: AppTheme.surfaceAlt,
    builder: (ctx) {
      if (playlists.isEmpty) {
        return const Padding(
          padding: EdgeInsets.all(24),
          child: Text('No playlists yet. Create one first.', style: TextStyle(color: Colors.white)),
        );
      }
      return ListView.separated(
        padding: const EdgeInsets.symmetric(vertical: 12),
        itemCount: playlists.length,
        separatorBuilder: (_, __) => const Divider(height: 1, color: Colors.white10),
        itemBuilder: (_, index) {
          final p = playlists[index];
          return ListTile(
            title: Text(p.name, style: const TextStyle(fontWeight: FontWeight.w700)),
            subtitle: Text(p.description.isEmpty ? 'No description' : p.description),
            trailing: const Icon(Icons.add, color: Colors.white70),
            onTap: () async {
              debugPrint('Add track $trackId to playlist ${p.id}');
              await ref.read(playlistRepositoryProvider).addTrack(p.id, trackId);
              if (context.mounted) Navigator.pop(context);
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('Added to ${p.name}')),
              );
            },
          );
        },
      );
    },
  );
}
