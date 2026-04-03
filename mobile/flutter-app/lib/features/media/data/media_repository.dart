import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import '../../../core/models.dart';

class MediaRepository {
  MediaRepository(this._dio);
  final Dio _dio;

  Future<MediaSourceModel> source(String trackId) async {
    final res = await _dio.get('/api/v1/media/tracks/$trackId/source');
    return MediaSourceModel.fromJson(res.data as Map<String, dynamic>);
  }
}

final mediaRepositoryProvider = Provider<MediaRepository>((ref) {
  return MediaRepository(ref.read(dioProvider));
});
