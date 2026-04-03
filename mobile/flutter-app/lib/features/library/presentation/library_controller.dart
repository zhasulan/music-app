import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/models.dart';
import '../data/library_repository.dart';

class LibraryState {
  final Set<String> likedIds;
  final bool loading;
  final String? error;

  const LibraryState({this.likedIds = const {}, this.loading = false, this.error});

  LibraryState copyWith({Set<String>? likedIds, bool? loading, String? error}) {
    return LibraryState(
      likedIds: likedIds ?? this.likedIds,
      loading: loading ?? this.loading,
      error: error,
    );
  }
}

class LibraryController extends StateNotifier<LibraryState> {
  LibraryController(this._repo) : super(const LibraryState()) {
    load();
  }

  final LibraryRepository _repo;

  Future<void> load() async {
    state = state.copyWith(loading: true, error: null);
    try {
      final items = await _repo.likedTracks();
      state = state.copyWith(likedIds: items.map((e) => e.trackId).toSet(), loading: false);
    } catch (e) {
      state = state.copyWith(error: 'Failed to load liked tracks', loading: false);
    }
  }

  Future<void> toggleLike(String trackId) async {
    final liked = state.likedIds.contains(trackId);
    // optimistic
    final updated = Set<String>.from(state.likedIds);
    if (liked) {
      updated.remove(trackId);
    } else {
      updated.add(trackId);
    }
    state = state.copyWith(likedIds: updated);
    try {
      if (liked) {
        await _repo.unlike(trackId);
      } else {
        await _repo.like(trackId);
      }
    } catch (_) {
      // rollback
      if (liked) {
        updated.add(trackId);
      } else {
        updated.remove(trackId);
      }
      state = state.copyWith(likedIds: updated, error: 'Failed to update like');
    }
  }
}

final libraryControllerProvider = StateNotifierProvider<LibraryController, LibraryState>((ref) {
  return LibraryController(ref.read(libraryRepositoryProvider));
});
