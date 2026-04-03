import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/formatters.dart';
import '../../../core/theme.dart';
import '../../catalog/presentation/catalog_screen.dart';
import '../../catalog/data/catalog_repository.dart';
import '../../playback/presentation/playback_controller.dart';
import '../../playlists/presentation/playlist_picker.dart';
import '../../audius/data/audius_repository.dart';
import '../../audius/data/audius_track.dart';
import '../../../core/models.dart';
import 'library_controller.dart';

class LibraryScreen extends ConsumerWidget {
  const LibraryScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final library = ref.watch(libraryControllerProvider);
    final tracksAsync = ref.watch(likedTrackDetailsProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Your Library')),
      body: library.loading
          ? const Center(child: CircularProgressIndicator())
          : tracksAsync.when(
        data: (tracks) {
          final likedTracks = tracks.where((t) => library.likedIds.contains(t.id)).toList();
          if (likedTracks.isEmpty) {
            return const Center(child: Text('Nothing liked yet'));
          }
          return ListView.builder(
            padding: const EdgeInsets.symmetric(vertical: 8),
            itemCount: likedTracks.length,
            itemBuilder: (context, index) {
              final t = likedTracks[index];
              return Card(
                child: Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                  child: Row(
                    children: [
                      const Icon(Icons.favorite, color: AppTheme.secondary),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(t.title, style: const TextStyle(fontWeight: FontWeight.w700)),
                            Text(formatDuration(t.durationSec), style: const TextStyle(color: AppTheme.textSecondary)),
                          ],
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.play_arrow, color: AppTheme.secondary),
                        onPressed: () => ref.read(playbackControllerProvider.notifier).playTrack(t),
                      ),
                      IconButton(
                        icon: const Icon(Icons.remove_circle_outline, color: AppTheme.textSecondary),
                        onPressed: () => ref.read(libraryControllerProvider.notifier).toggleLike(t.id),
                      ),
                      IconButton(
                        icon: const Icon(Icons.playlist_add, color: AppTheme.textSecondary),
                        onPressed: () => showPlaylistPicker(context, ref, trackId: t.id),
                      ),
                    ],
                  ),
                ),
              );
            },
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Failed to load library: $e')),
      ),
    );
  }
}

final likedTrackDetailsProvider = FutureProvider<List<TrackModel>>((ref) async {
  final likedIds = ref.watch(libraryControllerProvider).likedIds;
  final catalog = ref.read(catalogRepositoryProvider);
  final audius = ref.read(audiusRepositoryProvider);
  final List<TrackModel> result = [];
  for (final id in likedIds) {
    if (id.startsWith('audius:')) {
      final a = await audius.track(id);
      if (a != null) result.add(a.toTrackModel());
    } else {
      final t = await catalog.track(id);
      if (t != null) result.add(t);
    }
  }
  return result;
});
