import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/formatters.dart';
import '../../../core/models.dart';
import '../../catalog/data/catalog_repository.dart';
import '../../playback/presentation/playback_controller.dart';
import '../data/playlist_repository.dart';

final playlistDetailProvider = FutureProvider.family.autoDispose<PlaylistDetailVm, int>((ref, id) async {
  final repo = ref.read(playlistRepositoryProvider);
  final catalog = ref.read(catalogRepositoryProvider);
  final (pl, tracks) = await repo.details(id);
  final items = <PlaylistItem>[];
  for (final t in tracks) {
    final track = await catalog.track(t.trackId);
    if (track != null) {
      items.add(PlaylistItem(track: track, position: t.position));
    }
  }
  return PlaylistDetailVm(pl.name, items);
});

class PlaylistDetailVm {
  final String name;
  final List<PlaylistItem> items;
  PlaylistDetailVm(this.name, this.items);
}

class PlaylistItem {
  final TrackModel track;
  final int position;
  PlaylistItem({required this.track, required this.position});
}

class PlaylistDetailScreen extends ConsumerWidget {
  const PlaylistDetailScreen({super.key, required this.playlistId, required this.name});
  final int playlistId;
  final String name;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detail = ref.watch(playlistDetailProvider(playlistId));
    return Scaffold(
      appBar: AppBar(title: Text(name)),
      body: detail.when(
        data: (vm) => vm.items.isEmpty
            ? const Center(child: Text('Playlist is empty'))
            : ListView.builder(
                itemCount: vm.items.length,
                itemBuilder: (context, index) {
                  final item = vm.items[index];
                  final track = item.track;
                  return Card(
                    child: ListTile(
                      title: Text(track.title),
                      subtitle: Text('Pos ${item.position} • ${formatDuration(track.durationSec)}'),
                      trailing: IconButton(
                        icon: const Icon(Icons.play_arrow),
                        onPressed: () => ref.read(playbackControllerProvider.notifier).playTrack(track),
                      ),
                    ),
                  );
                },
              ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Failed: $e')),
      ),
    );
  }
}
