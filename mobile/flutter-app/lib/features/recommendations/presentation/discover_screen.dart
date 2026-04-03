import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/formatters.dart';
import '../../../core/models.dart';
import '../../../core/theme.dart';
import '../../library/presentation/library_controller.dart';
import '../../playback/presentation/playback_controller.dart';
import '../../playlists/presentation/playlist_picker.dart';
import '../data/recommendation_repository.dart';

final trendingProvider = FutureProvider<List<TrackModel>>((ref) {
  return ref.read(recommendationRepositoryProvider).trending();
});

final recentProvider = FutureProvider<List<TrackModel>>((ref) {
  return ref.read(recommendationRepositoryProvider).recentlyPlayed();
});

final forYouProvider = FutureProvider<List<TrackModel>>((ref) {
  return ref.read(recommendationRepositoryProvider).forYou();
});

class DiscoverScreen extends ConsumerWidget {
  const DiscoverScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final trending = ref.watch(trendingProvider);
    final recent = ref.watch(recentProvider);
    final forYou = ref.watch(forYouProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Discover')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(vertical: 12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _section(context, 'Trending', trending, ref),
            _section(context, 'Recently played', recent, ref),
            _section(context, 'For you', forYou, ref),
          ],
        ),
      ),
    );
  }

  Widget _section(BuildContext context, String title, AsyncValue<List<TrackModel>> asyncTracks, WidgetRef ref) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
            child: Text(title, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w700)),
          ),
          SizedBox(
            height: 180,
            child: asyncTracks.when(
              data: (tracks) => tracks.isEmpty
                  ? const _EmptyMessage()
                  : ListView.separated(
                      scrollDirection: Axis.horizontal,
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                      itemCount: tracks.length,
                      separatorBuilder: (_, __) => const SizedBox(width: 10),
                      itemBuilder: (context, index) => _TrackCard(track: tracks[index], ref: ref),
                    ),
              loading: () => const Center(child: Padding(
                padding: EdgeInsets.all(16),
                child: CircularProgressIndicator(),
              )),
              error: (e, _) => Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12),
                child: Text('Error: $e', style: const TextStyle(color: Colors.redAccent)),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class _EmptyMessage extends StatelessWidget {
  const _EmptyMessage();

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      child: Container(
        padding: const EdgeInsets.all(16),
        width: 200,
        child: const Text('Nothing here yet', style: TextStyle(color: AppTheme.textSecondary)),
      ),
    );
  }
}

class _TrackCard extends StatelessWidget {
  const _TrackCard({required this.track, required this.ref});
  final TrackModel track;
  final WidgetRef ref;

  @override
  Widget build(BuildContext context) {
    final liked = ref.watch(libraryControllerProvider.select((s) => s.likedIds.contains(track.id)));
    return Card(
      child: Container(
        width: 200,
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Align(
              alignment: Alignment.center,
              child: Icon(Icons.album, size: 42, color: AppTheme.secondary),
            ),
            const SizedBox(height: 8),
            Text(track.title, maxLines: 2, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w700)),
            const SizedBox(height: 4),
            Text(formatDuration(track.durationSec), style: const TextStyle(color: AppTheme.textSecondary)),
            const Spacer(),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                IconButton(
                  icon: Icon(Icons.favorite, color: liked ? AppTheme.secondary : AppTheme.textSecondary),
                  onPressed: () => ref.read(libraryControllerProvider.notifier).toggleLike(track.id),
                ),
                IconButton(
                  icon: const Icon(Icons.play_circle, color: AppTheme.secondary),
                  onPressed: () => ref.read(playbackControllerProvider.notifier).playTrack(track),
                ),
                IconButton(
                  icon: const Icon(Icons.add, color: AppTheme.textSecondary),
                  onPressed: () => showPlaylistPicker(context, ref, trackId: track.id),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
