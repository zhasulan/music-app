import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';

class EventsRepository {
  EventsRepository(this._dio);
  final Dio _dio;

  Future<void> trackPlayed(String trackId) => _post('/api/v1/events/track-played', {'track_id': trackId});
  Future<void> trackPaused(String trackId) => _post('/api/v1/events/track-paused', {'track_id': trackId});
  Future<void> trackLiked(String trackId) => _post('/api/v1/events/track-liked', {'track_id': trackId});
  Future<void> playlistCreated(int playlistId) => _post('/api/v1/events/playlist-created', {'playlist_id': playlistId});
  Future<void> playlistTrackAdded(int playlistId, String trackId) =>
      _post('/api/v1/events/playlist-track-added', {'playlist_id': playlistId, 'track_id': trackId});
  Future<void> searchPerformed(String query) => _post('/api/v1/events/search', {'query': query});

  Future<void> _post(String path, Map<String, dynamic> body) async {
    try {
      await _dio.post(path, data: body);
    } catch (_) {
      // Fire-and-forget; ignore failures for UX stability.
    }
  }
}

final eventsRepositoryProvider = Provider<EventsRepository>((ref) {
  return EventsRepository(ref.read(dioProvider));
});
