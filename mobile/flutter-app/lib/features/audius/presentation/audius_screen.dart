import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:just_audio/just_audio.dart';

import '../../../core/formatters.dart';
import '../../../core/theme.dart';
import '../data/audius_repository.dart';
import '../data/audius_track.dart';

final audiusTrendingProvider = FutureProvider<List<AudiusTrack>>((ref) {
  return ref.read(audiusRepositoryProvider).trending();
});

class AudiusScreen extends ConsumerStatefulWidget {
  const AudiusScreen({super.key});

  @override
  ConsumerState<AudiusScreen> createState() => _AudiusScreenState();
}

class _AudiusScreenState extends ConsumerState<AudiusScreen> {
  final _player = AudioPlayer();
  String? _playingId;

  @override
  void dispose() {
    _player.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final tracks = ref.watch(audiusTrendingProvider);
    return Scaffold(
      appBar: AppBar(title: const Text('Audius Trending')),
      body: tracks.when(
        data: (items) => ListView.builder(
          itemCount: items.length,
          itemBuilder: (context, index) {
            final t = items[index];
            final playing = t.id == _playingId;
            return Card(
              child: ListTile(
                leading: const _ProviderBadge(),
                title: Text(t.title),
                subtitle: Text('${t.artistName} • ${formatDuration((t.duration / 1000).round())}'),
                trailing: IconButton(
                  icon: Icon(playing ? Icons.pause_circle : Icons.play_circle, color: AppTheme.secondary),
                  onPressed: t.streamUrl == null
                      ? null
                      : () async {
                          if (playing) {
                            await _player.pause();
                            setState(() => _playingId = null);
                          } else {
                            await _player.setUrl(t.streamUrl!);
                            await _player.play();
                            setState(() => _playingId = t.id);
                          }
                        },
                ),
              ),
            );
          },
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
