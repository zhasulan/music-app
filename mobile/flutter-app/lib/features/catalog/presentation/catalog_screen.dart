import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/formatters.dart';
import '../../../core/theme.dart';
import '../../library/presentation/library_controller.dart';
import '../../playback/presentation/playback_controller.dart';
import '../../playlists/presentation/playlist_picker.dart';
import '../data/catalog_repository.dart';

final tracksProvider = FutureProvider((ref) => ref.read(catalogRepositoryProvider).tracks());

class CatalogScreen extends ConsumerWidget {
  const CatalogScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final tracksAsync = ref.watch(tracksProvider);
    final library = ref.watch(libraryControllerProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Catalog'),
        centerTitle: false,
      ),
      body: tracksAsync.when(
        data: (tracks) => ListView.builder(
          padding: const EdgeInsets.symmetric(vertical: 8),
          itemCount: tracks.length,
          itemBuilder: (context, index) {
            final t = tracks[index];
            final liked = library.likedIds.contains(t.id);
            return Card(
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                child: Row(
                  children: [
                    CircleAvatar(
                      backgroundColor: AppTheme.primary.withOpacity(0.15),
                      child: const Icon(Icons.music_note, color: AppTheme.secondary),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(t.title, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 16)),
                          const SizedBox(height: 4),
                          Text(formatDuration(t.durationSec), style: const TextStyle(color: AppTheme.textSecondary)),
                        ],
                      ),
                    ),
                    IconButton(
                      icon: Icon(Icons.favorite, color: liked ? AppTheme.secondary : AppTheme.textSecondary),
                      onPressed: () => ref.read(libraryControllerProvider.notifier).toggleLike(t.id),
                    ),
                    IconButton(
                      icon: const Icon(Icons.play_arrow, color: AppTheme.secondary),
                      onPressed: () {
                        debugPrint('UI play tap ${t.id}');
                        ref.read(playbackControllerProvider.notifier).playTrack(t);
                      },
                    ),
                    IconButton(
                      icon: const Icon(Icons.add, color: AppTheme.textSecondary),
                      onPressed: () => showPlaylistPicker(context, ref, trackId: t.id),
                    ),
                  ],
                ),
              ),
            );
          },
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Failed to load tracks: $e')),
      ),
    );
  }
}
