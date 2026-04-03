import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/models.dart';
import '../data/playlist_repository.dart';

class PlaylistsState {
  final List<PlaylistModel> playlists;
  final bool loading;
  final String? error;

  const PlaylistsState({this.playlists = const [], this.loading = false, this.error});

  PlaylistsState copyWith({List<PlaylistModel>? playlists, bool? loading, String? error}) {
    return PlaylistsState(
      playlists: playlists ?? this.playlists,
      loading: loading ?? this.loading,
      error: error,
    );
  }
}

class PlaylistsController extends StateNotifier<PlaylistsState> {
  PlaylistsController(this._repo) : super(const PlaylistsState()) {
    load();
  }
  final PlaylistRepository _repo;

  Future<void> load() async {
    state = state.copyWith(loading: true, error: null);
    try {
      final data = await _repo.list();
      state = state.copyWith(playlists: data, loading: false);
    } catch (e) {
      state = state.copyWith(error: 'Failed to load playlists: $e', loading: false);
    }
  }

  Future<void> create(String name, String description) async {
    try {
      await _repo.create(name, description);
      await load();
    } catch (e) {
      state = state.copyWith(error: 'Failed to create playlist: $e');
    }
  }

  Future<void> delete(int id) async {
    try {
      await _repo.delete(id);
      await load();
    } catch (e) {
      state = state.copyWith(error: 'Failed to delete playlist: $e');
    }
  }
}

final playlistsControllerProvider = StateNotifierProvider<PlaylistsController, PlaylistsState>((ref) {
  return PlaylistsController(ref.read(playlistRepositoryProvider));
});
