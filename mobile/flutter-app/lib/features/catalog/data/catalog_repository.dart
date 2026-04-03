import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import '../../../core/models.dart';

class CatalogRepository {
  CatalogRepository(this._dio);
  final Dio _dio;

  Future<List<TrackModel>> tracks() async {
    final res = await _dio.get('/api/v1/catalog/tracks');
    final list = res.data['tracks'] as List<dynamic>;
    return list.map((e) => TrackModel.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<TrackModel?> track(String id) async {
    final res = await _dio.get('/api/v1/catalog/tracks/$id');
    if (res.data['track'] == null) return null;
    return TrackModel.fromJson(res.data['track'] as Map<String, dynamic>);
  }
}

final catalogRepositoryProvider = Provider<CatalogRepository>((ref) {
  return CatalogRepository(ref.read(dioProvider));
});
