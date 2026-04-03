import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import 'audius_track.dart';

class AudiusRepository {
  AudiusRepository(this._dio);
  final Dio _dio;

  Future<List<AudiusTrack>> trending({int limit = 20, int offset = 0}) async {
    final res = await _dio.get('/api/v1/providers/audius/tracks/trending', queryParameters: {'limit': limit, 'offset': offset});
    final list = res.data['items'] as List<dynamic>? ?? res.data['tracks'] as List<dynamic>? ?? [];
    return list.map((e) => AudiusTrack.fromJson(e as Map<String, dynamic>)).toList();
  }
}

final audiusRepositoryProvider = Provider<AudiusRepository>((ref) {
  return AudiusRepository(ref.read(dioProvider));
});
