class UserModel {
  final int id;
  final String email;
  final String name;
  final String country;

  UserModel({required this.id, required this.email, this.name = '', this.country = ''});

  factory UserModel.fromJson(Map<String, dynamic> json) {
    return UserModel(
      id: json['id'] as int,
      email: json['email'] as String? ?? '',
      name: json['name'] as String? ?? '',
      country: json['country'] as String? ?? '',
    );
  }
}

class TrackModel {
  final String id;
  final String title;
  final String artistId;
  final String albumId;
  final int durationSec;

  TrackModel({
    required this.id,
    required this.title,
    required this.artistId,
    required this.albumId,
    required this.durationSec,
  });

  factory TrackModel.fromJson(Map<String, dynamic> json) {
    return TrackModel(
      id: json['id'] as String,
      title: json['title'] as String,
      artistId: json['artist_id'] as String? ?? '',
      albumId: json['album_id'] as String? ?? '',
      durationSec: json['duration_sec'] as int? ?? 0,
    );
  }
}

class PlaylistModel {
  final int id;
  final String name;
  final String description;

  PlaylistModel({required this.id, required this.name, this.description = ''});

  factory PlaylistModel.fromJson(Map<String, dynamic> json) {
    return PlaylistModel(
      id: json['id'] as int,
      name: json['name'] as String? ?? '',
      description: json['description'] as String? ?? '',
    );
  }
}

class PlaylistTrackModel {
  final int id;
  final String trackId;
  final int position;

  PlaylistTrackModel({required this.id, required this.trackId, required this.position});

  factory PlaylistTrackModel.fromJson(Map<String, dynamic> json) {
    return PlaylistTrackModel(
      id: json['id'] as int,
      trackId: json['track_id'] as String,
      position: json['position'] as int? ?? 0,
    );
  }
}

class LikedTrackModel {
  final String trackId;

  LikedTrackModel({required this.trackId});

  factory LikedTrackModel.fromJson(Map<String, dynamic> json) {
    return LikedTrackModel(trackId: json['track_id'] as String);
  }
}

class PlaybackSessionModel {
  final String trackId;
  final int positionMs;
  final String state;

  PlaybackSessionModel({required this.trackId, required this.positionMs, required this.state});

  factory PlaybackSessionModel.fromJson(Map<String, dynamic> json) {
    return PlaybackSessionModel(
      trackId: json['track_id'] as String,
      positionMs: json['position_ms'] as int? ?? 0,
      state: json['state'] as String? ?? 'paused',
    );
  }
}

class MediaSourceModel {
  final String url;
  MediaSourceModel({required this.url});

  factory MediaSourceModel.fromJson(Map<String, dynamic> json) {
    return MediaSourceModel(url: json['url'] as String);
  }
}
