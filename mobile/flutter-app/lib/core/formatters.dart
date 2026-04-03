String formatDuration(int seconds) {
  final dur = Duration(seconds: seconds);
  final m = dur.inMinutes;
  final s = dur.inSeconds % 60;
  return '${m.toString().padLeft(2, '0')}:${s.toString().padLeft(2, '0')}';
}

String formatMillis(int ms) {
  final dur = Duration(milliseconds: ms);
  final m = dur.inMinutes;
  final s = dur.inSeconds % 60;
  return '${m.toString().padLeft(2, '0')}:${s.toString().padLeft(2, '0')}';
}
