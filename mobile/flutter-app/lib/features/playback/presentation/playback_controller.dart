import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter/foundation.dart';
import 'package:just_audio/just_audio.dart';

import '../../../core/models.dart';
import '../../catalog/data/catalog_repository.dart';
import '../../media/data/media_repository.dart';
import '../data/playback_repository.dart';

class PlaybackViewState {
  final TrackModel? track;
  final bool playing;
  final bool loading;
  final int positionMs;
  final int durationMs;
  final String? error;

  const PlaybackViewState({
    this.track,
    this.playing = false,
    this.loading = false,
    this.positionMs = 0,
    this.durationMs = 0,
    this.error,
  });

  PlaybackViewState copyWith({
    TrackModel? track,
    bool? playing,
    bool? loading,
    int? positionMs,
    int? durationMs,
    String? error,
  }) {
    return PlaybackViewState(
      track: track ?? this.track,
      playing: playing ?? this.playing,
      loading: loading ?? this.loading,
      positionMs: positionMs ?? this.positionMs,
      durationMs: durationMs ?? this.durationMs,
      error: error,
    );
  }
}

class PlaybackController extends StateNotifier<PlaybackViewState> {
  PlaybackController(this._repo, this._mediaRepo, this._catalogRepo)
      : _player = AudioPlayer(),
        super(const PlaybackViewState()) {
    _bindStreams();
    _restoreFromServer();
  }

  final PlaybackRepository _repo;
  final MediaRepository _mediaRepo;
  final CatalogRepository _catalogRepo;
  final AudioPlayer _player;
  StreamSubscription<Duration>? _positionSub;
  StreamSubscription<Duration?>? _durationSub;

  void _bindStreams() {
    _positionSub = _player.positionStream.listen((pos) {
      state = state.copyWith(positionMs: pos.inMilliseconds);
    });
    _durationSub = _player.durationStream.listen((dur) {
      if (dur != null) {
        state = state.copyWith(durationMs: dur.inMilliseconds);
      }
    });
  }

  Future<void> _restoreFromServer() async {
    try {
      final session = await _repo.current();
      if (session == null) return;
      final track = await _catalogRepo.track(session.trackId);
      if (track == null) return;
      final source = await _mediaRepo.source(track.id);
      await _player.setUrl(source.url);
      if (session.positionMs > 0) {
        await _player.seek(Duration(milliseconds: session.positionMs));
      }
      state = state.copyWith(track: track, playing: session.state == 'playing', positionMs: session.positionMs);
      if (session.state == 'playing') {
        await _player.play();
      }
    } catch (_) {
      // ignore restore failures
    }
  }

  Future<void> playTrack(TrackModel track) async {
    state = state.copyWith(
      track: track,
      playing: true,
      loading: true,
      error: null,
      positionMs: 0,
      durationMs: 0,
    );
    try {
      await _player.stop();
      final source = await _mediaRepo.source(track.id);
      debugPrint('PLAY start track=${track.id} url=${source.url}');
      await _player.setUrl(source.url);
      await _repo.start(track.id, 0);
      await _player.play();
      state = state.copyWith(
        loading: false,
        durationMs: _player.duration?.inMilliseconds ?? 0,
      );
    } catch (e) {
      state = state.copyWith(playing: false, loading: false, error: 'Unable to play');
    }
  }

  Future<void> pause() async {
    if (state.track == null) return;
    final pos = _player.position.inMilliseconds;
    state = state.copyWith(playing: false, positionMs: pos);
    try {
      await _player.pause();
      debugPrint('PLAY pause pos=$pos');
      await _repo.pause(pos);
    } catch (_) {
      // if backend pause fails, keep local paused state
    }
  }

  Future<void> resume() async {
    if (state.track == null) return;
    debugPrint('PLAY resume track=${state.track!.id}');
    state = state.copyWith(playing: true, error: null);
    try {
      await _repo.resume();
      await _player.play();
    } catch (_) {
      state = state.copyWith(playing: false, error: 'Unable to resume');
    }
  }

  Future<void> seekTo(int positionMs) async {
    if (state.track == null) return;
    debugPrint('PLAY seek pos=$positionMs');
    await _player.seek(Duration(milliseconds: positionMs));
    await _repo.seek(positionMs);
    state = state.copyWith(positionMs: positionMs);
  }

  @override
  void dispose() {
    _positionSub?.cancel();
    _durationSub?.cancel();
    _player.dispose();
    super.dispose();
  }
}

final playbackControllerProvider = StateNotifierProvider<PlaybackController, PlaybackViewState>((ref) {
  final repo = ref.read(playbackRepositoryProvider);
  final media = ref.read(mediaRepositoryProvider);
  final catalog = ref.read(catalogRepositoryProvider);
  return PlaybackController(repo, media, catalog);
});
