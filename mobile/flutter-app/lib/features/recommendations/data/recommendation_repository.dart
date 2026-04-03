import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import '../../../core/models.dart';

class RecommendationRepository {
  RecommendationRepository(this._dio);
  final Dio _dio;

  Future<List<TrackModel>> trending() async {
    final res = await _dio.get('/api/v1/recommendations/trending');
    return _parseTracks(res.data['tracks'] as List<dynamic>? ?? []);
  }

  Future<List<TrackModel>> recentlyPlayed() async {
    final res = await _dio.get('/api/v1/recommendations/recently-played');
    return _parseTracks(res.data['tracks'] as List<dynamic>? ?? []);
  }

  Future<List<TrackModel>> forYou() async {
    final res = await _dio.get('/api/v1/recommendations/for-you');
    return _parseTracks(res.data['tracks'] as List<dynamic>? ?? []);
  }

  List<TrackModel> _parseTracks(List<dynamic> list) {
    return list
        .map((e) => TrackModel.fromJson(e as Map<String, dynamic>))
        .toList();
  }
}

final recommendationRepositoryProvider = Provider<RecommendationRepository>((ref) {
  return RecommendationRepository(ref.read(dioProvider));
});
