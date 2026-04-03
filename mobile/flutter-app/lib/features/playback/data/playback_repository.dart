import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import '../../../core/models.dart';

class PlaybackRepository {
  PlaybackRepository(this._dio);
  final Dio _dio;

  Future<PlaybackSessionModel> start(String trackId, int positionMs) async {
    final res = await _dio.post('/api/v1/playback/start', data: {
      'track_id': trackId,
      'position_ms': positionMs,
    });
    return PlaybackSessionModel.fromJson(res.data['session'] as Map<String, dynamic>);
  }

  Future<PlaybackSessionModel> pause(int positionMs) async {
    final res = await _dio.post('/api/v1/playback/pause', data: {'position_ms': positionMs});
    return PlaybackSessionModel.fromJson(res.data['session'] as Map<String, dynamic>);
  }

  Future<PlaybackSessionModel> resume() async {
    final res = await _dio.post('/api/v1/playback/resume');
    return PlaybackSessionModel.fromJson(res.data['session'] as Map<String, dynamic>);
  }

  Future<PlaybackSessionModel> seek(int positionMs) async {
    final res = await _dio.post('/api/v1/playback/seek', data: {'position_ms': positionMs});
    return PlaybackSessionModel.fromJson(res.data['session'] as Map<String, dynamic>);
  }

  Future<PlaybackSessionModel?> current() async {
    try {
      final res = await _dio.get('/api/v1/playback/current');
      if (res.data['session'] == null) return null;
      return PlaybackSessionModel.fromJson(res.data['session'] as Map<String, dynamic>);
    } on DioException catch (e) {
      if (e.response?.statusCode == 404) {
        return null; // no active session
      }
      rethrow;
    }
  }
}

final playbackRepositoryProvider = Provider<PlaybackRepository>((ref) {
  return PlaybackRepository(ref.read(dioProvider));
});
