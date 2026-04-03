import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/formatters.dart';
import '../../../core/theme.dart';
import 'playback_controller.dart';

class PlayerBar extends ConsumerWidget {
  const PlayerBar({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(playbackControllerProvider);
    final controller = ref.read(playbackControllerProvider.notifier);
    final track = state.track;
    if (track == null) return const SizedBox.shrink();
    final duration = state.durationMs > 0 ? state.durationMs : track.durationSec * 1000;
    final position = state.positionMs.clamp(0, duration);

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: const BoxDecoration(
        color: AppTheme.surfaceAlt,
        border: Border(top: BorderSide(color: Colors.white12)),
      ),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            children: [
              const Icon(Icons.music_note, color: AppTheme.secondary),
              const SizedBox(width: 8),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(track.title, maxLines: 1, overflow: TextOverflow.ellipsis, style: const TextStyle(fontWeight: FontWeight.w600)),
                    Text(formatDuration(track.durationSec), style: const TextStyle(color: AppTheme.textSecondary, fontSize: 12)),
                  ],
                ),
              ),
              IconButton(
                icon: Icon(state.playing ? Icons.pause_circle : Icons.play_circle, size: 36, color: AppTheme.secondary),
                onPressed: () {
                  if (state.playing) {
                    controller.pause();
                  } else {
                    controller.resume();
                  }
                },
              ),
            ],
          ),
          Padding(
            padding: const EdgeInsets.only(top: 4.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(formatMillis(position), style: const TextStyle(fontSize: 11, color: AppTheme.textSecondary)),
                Text(formatMillis(duration), style: const TextStyle(fontSize: 11, color: AppTheme.textSecondary)),
              ],
            ),
          ),
          SliderTheme(
            data: SliderTheme.of(context).copyWith(trackHeight: 2, thumbShape: const RoundSliderThumbShape(enabledThumbRadius: 6)),
            child: Slider(
              min: 0,
              max: duration.toDouble().clamp(1, double.infinity),
              value: position.toDouble(),
              onChanged: (value) => controller.seekTo(value.toInt()),
              activeColor: AppTheme.secondary,
            ),
          ),
        ],
      ),
    );
  }
}
