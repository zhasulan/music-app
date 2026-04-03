import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../events/data/events_repository.dart';
import '../data/search_repository.dart';

class SearchState {
  final String query;
  final bool loading;
  final List<SearchResult> results;
  final List<String> suggestions;
  final String? error;

  const SearchState({
    this.query = '',
    this.loading = false,
    this.results = const [],
    this.suggestions = const [],
    this.error,
  });

  SearchState copyWith({
    String? query,
    bool? loading,
    List<SearchResult>? results,
    List<String>? suggestions,
    String? error,
  }) {
    return SearchState(
      query: query ?? this.query,
      loading: loading ?? this.loading,
      results: results ?? this.results,
      suggestions: suggestions ?? this.suggestions,
      error: error,
    );
  }
}

class SearchController extends StateNotifier<SearchState> {
  SearchController(this._repo, this._events) : super(const SearchState());

  final SearchRepository _repo;
  final EventsRepository _events;
  Timer? _debounce;

  void updateQuery(String q) {
    state = state.copyWith(query: q, loading: true, error: null);
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 350), () async {
      if (q.trim().isEmpty) {
        state = state.copyWith(loading: false, results: [], suggestions: []);
        return;
      }
      try {
        final suggestions = await _repo.suggest(q);
        final results = await _repo.search(q);
        unawaited(_events.searchPerformed(q));
        state = state.copyWith(loading: false, results: results, suggestions: suggestions);
      } catch (e) {
        state = state.copyWith(loading: false, error: 'Search failed: $e');
      }
    });
  }
}

final searchControllerProvider = StateNotifierProvider<SearchController, SearchState>((ref) {
  return SearchController(ref.read(searchRepositoryProvider), ref.read(eventsRepositoryProvider));
});
