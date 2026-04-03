import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import '../../../core/models.dart';

class SearchResult {
  final String type;
  final TrackModel? track;
  final String title;

  SearchResult({required this.type, required this.title, this.track});
}

class SearchRepository {
  SearchRepository(this._dio);
  final Dio _dio;

  Future<List<SearchResult>> search(String query) async {
    final res = await _dio.get('/api/v1/search', queryParameters: {'q': query});
    final list = res.data['results'] as List<dynamic>? ?? [];
    return list.map((e) => _mapResult(e as Map<String, dynamic>)).toList();
  }

  Future<List<String>> suggest(String query) async {
    final res = await _dio.get('/api/v1/search/suggest', queryParameters: {'q': query});
    final list = res.data['suggestions'] as List<dynamic>? ?? [];
    return list.map((e) => (e as Map<String, dynamic>)['text'] as String? ?? '').where((e) => e.isNotEmpty).toList();
  }

  SearchResult _mapResult(Map<String, dynamic> json) {
    final type = json['type'] as String? ?? 'track';
    if (type == 'track') {
      final track = TrackModel(
        id: json['track']?['id'] as String? ?? json['id'] as String,
        title: json['track']?['title'] as String? ?? json['title'] as String? ?? '',
        artistId: (json['track']?['artist_id'] ?? json['artist_id'] ?? '') as String,
        albumId: (json['track']?['album_id'] ?? json['album_id'] ?? '') as String,
        durationSec: json['track']?['duration_sec'] as int? ?? json['duration_sec'] as int? ?? 0,
      );
      return SearchResult(type: 'track', title: track.title, track: track);
    }
    final title = json['album']?['title'] ?? json['artist']?['name'] ?? json['title'] ?? json['name'] ?? '';
    return SearchResult(type: type, title: title.toString());
  }
}

final searchRepositoryProvider = Provider<SearchRepository>((ref) {
  return SearchRepository(ref.read(dioProvider));
});
