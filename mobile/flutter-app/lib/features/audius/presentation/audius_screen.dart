import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/formatters.dart';
import '../../../core/theme.dart';
import '../data/audius_repository.dart';
import '../data/audius_track.dart';
import '../../playback/presentation/playback_controller.dart';
import '../../library/presentation/library_controller.dart';
import '../../playlists/presentation/playlist_picker.dart';

typedef AudiusPage = (int offset, int limit);

final audiusTrendingProviderFamily = FutureProvider.autoDispose.family<List<AudiusTrack>, AudiusPage>((ref, page) {
  final (offset, limit) = page;
  return ref.read(audiusRepositoryProvider).trending(limit: limit, offset: offset);
});

final audiusTrendingProvider = FutureProvider<List<AudiusTrack>>((ref) {
  return ref.read(audiusRepositoryProvider).trending();
});

class AudiusScreen extends ConsumerStatefulWidget {
  const AudiusScreen({super.key});

  @override
  ConsumerState<AudiusScreen> createState() => _AudiusScreenState();
}

class _AudiusScreenState extends ConsumerState<AudiusScreen> {
  int _offset = 0;
  final int _pageSize = 20;

  @override
  Widget build(BuildContext context) {
    final tracks = ref.watch(audiusTrendingProviderFamily((_offset, _pageSize)));
    final playbackState = ref.watch(playbackControllerProvider);
    return Scaffold(
      appBar: AppBar(title: const Text('Audius Trending')),
      body: tracks.when(
        data: (items) => RefreshIndicator(
          onRefresh: () async {
            setState(() {
              _offset = 0;
            });
            ref.refresh(audiusTrendingProviderFamily((_offset, _pageSize)));
          },
          child: ListView.builder(
            itemCount: items.length + 1,
            itemBuilder: (context, index) {
              if (index == items.length) {
                return TextButton(
                  onPressed: () {
                    setState(() => _offset += _pageSize);
                    ref.refresh(audiusTrendingProviderFamily((_offset, _pageSize)));
                  },
                  child: const Text('Load more'),
                );
              }
              final t = items[index];
              final playing = playbackState.playing && playbackState.track?.id == t.id;
              return Card(
                child: ListTile(
                  leading: const _ProviderBadge(),
                  title: Text(t.title),
                  subtitle: Text('${t.artistName} • ${formatDuration(t.duration)}'),
                  trailing: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      IconButton(
                        icon: Icon(playing ? Icons.pause_circle : Icons.play_circle, color: AppTheme.secondary),
                        onPressed: t.streamUrl == null
                            ? null
                            : () async {
                                if (playing) {
                                  await ref.read(playbackControllerProvider.notifier).pause();
                                } else {
                                  final playback = ref.read(playbackControllerProvider.notifier);
                                  await playback.playExternal(
                                    id: t.id,
                                    title: t.title,
                                    artist: t.artistName,
                                    url: t.streamUrl!,
                                    durationSec: t.duration,
                                  );
                                }
                              },
                      ),
                      IconButton(
                        icon: const Icon(Icons.favorite_border, color: AppTheme.textSecondary),
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
          ),
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Audius unavailable', style: TextStyle(color: Colors.red[200]))),
      ),
    );
  }
}

class _ProviderBadge extends StatelessWidget {
  const _ProviderBadge();

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: Colors.purple.withOpacity(0.15),
        borderRadius: BorderRadius.circular(8),
      ),
      child: const Text('Audius', style: TextStyle(color: Colors.purpleAccent, fontWeight: FontWeight.w700)),
    );
  }
}
