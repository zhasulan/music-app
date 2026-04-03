import '../../../core/models.dart';

class AudiusTrack {
  final String id;
  final String providerTrackId;
  final String title;
  final String artistName;
  final String artistHandle;
  final String? artworkUrl;
  final int duration;
  final String? streamUrl;

  AudiusTrack({
    required this.id,
    required this.providerTrackId,
    required this.title,
    required this.artistName,
    required this.artistHandle,
    required this.artworkUrl,
    required this.duration,
    required this.streamUrl,
  });

  factory AudiusTrack.fromJson(Map<String, dynamic> json) {
    return AudiusTrack(
      id: json['id'] as String? ?? '',
      providerTrackId: json['providerTrackId'] as String? ?? '',
      title: json['title'] as String? ?? '',
      artistName: json['artistName'] as String? ?? '',
      artistHandle: json['artistHandle'] as String? ?? '',
      artworkUrl: json['artworkUrl'] as String?,
      duration: json['duration'] as int? ?? 0,
      streamUrl: json['streamUrl'] as String?,
    );
  }

  TrackModel toTrackModel() => TrackModel(
        id: id,
        title: title,
        artistId: artistHandle,
        albumId: '',
        durationSec: duration,
      );
}
