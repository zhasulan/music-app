import 'package:flutter/material.dart' hide SearchController;
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/formatters.dart';
import '../../../core/theme.dart';
import '../../playback/presentation/playback_controller.dart';
import '../../playlists/presentation/playlist_picker.dart';
import '../../library/presentation/library_controller.dart';
import '../data/search_repository.dart';
import 'search_controller.dart';

class SearchScreen extends ConsumerWidget {
  const SearchScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(searchControllerProvider);
    final controller = ref.read(searchControllerProvider.notifier);

    return Scaffold(
      appBar: AppBar(title: const Text('Search')),
      body: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(12),
            child: TextField(
              decoration: const InputDecoration(
                hintText: 'Search tracks, albums, artists',
                prefixIcon: Icon(Icons.search),
              ),
              onChanged: controller.updateQuery,
            ),
          ),
          if (state.loading) const LinearProgressIndicator(minHeight: 2),
          if (state.error != null)
            Padding(
              padding: const EdgeInsets.all(12),
              child: Text(state.error!, style: const TextStyle(color: Colors.redAccent)),
            ),
          Expanded(
            child: state.query.isEmpty
                ? const _Placeholder()
                : ListView(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                    children: [
                      if (state.suggestions.isNotEmpty) _suggestions(state.suggestions, controller),
                      _resultsSection('Tracks', state.results.where((r) => r.type == 'track').toList(), ref),
                      _resultsSection('Albums', state.results.where((r) => r.type == 'album').toList(), ref),
                      _resultsSection('Artists', state.results.where((r) => r.type == 'artist').toList(), ref),
                    ],
                  ),
          ),
        ],
      ),
    );
  }

  Widget _suggestions(List<String> suggestions, SearchController controller) {
    return Wrap(
      spacing: 8,
      runSpacing: 4,
      children: suggestions
          .map((s) => ActionChip(
                label: Text(s),
                onPressed: () => controller.updateQuery(s),
              ))
          .toList(),
    );
  }

  Widget _resultsSection(String title, List<SearchResult> results, WidgetRef ref) {
    if (results.isEmpty) return const SizedBox.shrink();
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.only(top: 16, bottom: 8),
          child: Text(title, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w700)),
        ),
        ...results.map((r) => _ResultTile(result: r, ref: ref)),
      ],
    );
  }
}

class _Placeholder extends StatelessWidget {
  const _Placeholder();

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Text(
          'Search for tracks, albums or artists.',
          style: TextStyle(color: AppTheme.textSecondary.withOpacity(0.8)),
          textAlign: TextAlign.center,
        ),
      ),
    );
  }
}

class _ResultTile extends StatelessWidget {
  const _ResultTile({required this.result, required this.ref});
  final SearchResult result;
  final WidgetRef ref;

  @override
  Widget build(BuildContext context) {
    final track = result.track;
    final liked = track != null ? ref.watch(libraryControllerProvider.select((s) => s.likedIds.contains(track.id))) : false;

    return Card(
      child: ListTile(
        leading: Icon(
          result.type == 'track'
              ? Icons.music_note
              : result.type == 'album'
                  ? Icons.album
                  : Icons.person,
          color: AppTheme.secondary,
        ),
        title: Text(result.title),
        subtitle: track != null
            ? Text('${formatDuration(track.durationSec)} • ${result.provider}')
            : Text(result.type),
        trailing: track == null
            ? null
            : Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  IconButton(
                    icon: Icon(Icons.favorite, color: liked ? AppTheme.secondary : AppTheme.textSecondary),
                    onPressed: () => ref.read(libraryControllerProvider.notifier).toggleLike(track.id),
                  ),
                  IconButton(
                    icon: const Icon(Icons.play_arrow, color: AppTheme.secondary),
                    onPressed: () => ref.read(playbackControllerProvider.notifier).playTrack(track),
                  ),
                  IconButton(
                    icon: const Icon(Icons.add, color: AppTheme.textSecondary),
                    onPressed: () => showPlaylistPicker(context, ref, trackId: track.id),
                  ),
                ],
              ),
      ),
    );
  }
}
