import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import '../../../core/models.dart';

class LibraryRepository {
  LibraryRepository(this._dio);
  final Dio _dio;

  Future<List<LikedTrackModel>> likedTracks() async {
    final res = await _dio.get('/api/v1/library/liked-tracks');
    final list = res.data['liked_tracks'] as List<dynamic>? ?? [];
    return list.map((e) => LikedTrackModel.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<void> like(String trackId) async {
    await _dio.put('/api/v1/library/liked-tracks/$trackId');
  }

  Future<void> unlike(String trackId) async {
    await _dio.delete('/api/v1/library/liked-tracks/$trackId');
  }
}

final libraryRepositoryProvider = Provider<LibraryRepository>((ref) {
  return LibraryRepository(ref.read(dioProvider));
});
