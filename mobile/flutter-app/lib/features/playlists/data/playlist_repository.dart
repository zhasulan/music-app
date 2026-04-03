import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import '../../../core/models.dart';

class PlaylistRepository {
  PlaylistRepository(this._dio);
  final Dio _dio;

  Future<List<PlaylistModel>> list() async {
    final res = await _dio.get('/api/v1/playlists');
    final list = res.data['playlists'] as List<dynamic>? ?? [];
    return list.map((e) => PlaylistModel.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<PlaylistModel> create(String name, String description) async {
    final res = await _dio.post('/api/v1/playlists', data: {'name': name, 'description': description});
    return PlaylistModel.fromJson(res.data['playlist'] as Map<String, dynamic>);
  }

  Future<PlaylistModel> update(int id, String name, String description) async {
    final res = await _dio.patch('/api/v1/playlists/$id', data: {'name': name, 'description': description});
    return PlaylistModel.fromJson(res.data['playlist'] as Map<String, dynamic>);
  }

  Future<void> delete(int id) async {
    await _dio.delete('/api/v1/playlists/$id');
  }

  Future<(PlaylistModel, List<PlaylistTrackModel>)> details(int id) async {
    final res = await _dio.get('/api/v1/playlists/$id');
    final playlist = PlaylistModel.fromJson(res.data['playlist'] as Map<String, dynamic>);
    final tracks = (res.data['tracks'] as List<dynamic>? ?? [])
        .map((e) => PlaylistTrackModel.fromJson(e as Map<String, dynamic>))
        .toList();
    return (playlist, tracks);
  }

  Future<List<PlaylistTrackModel>> addTrack(int playlistId, String trackId) async {
    final res = await _dio.post('/api/v1/playlists/$playlistId/tracks', data: {'track_id': trackId});
    final tracks = (res.data['tracks'] as List<dynamic>? ?? [])
        .map((e) => PlaylistTrackModel.fromJson(e as Map<String, dynamic>))
        .toList();
    return tracks;
  }

  Future<List<PlaylistTrackModel>> removeTrack(int playlistId, String trackId) async {
    final res = await _dio.delete('/api/v1/playlists/$playlistId/tracks/$trackId');
    final tracks = (res.data['tracks'] as List<dynamic>? ?? [])
        .map((e) => PlaylistTrackModel.fromJson(e as Map<String, dynamic>))
        .toList();
    return tracks;
  }
}

final playlistRepositoryProvider = Provider<PlaylistRepository>((ref) {
  return PlaylistRepository(ref.read(dioProvider));
});
