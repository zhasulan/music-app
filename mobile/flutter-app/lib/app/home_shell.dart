import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../features/catalog/presentation/catalog_screen.dart';
import '../features/recommendations/presentation/discover_screen.dart';
import '../features/playlists/presentation/playlists_screen.dart';
import '../features/library/presentation/library_screen.dart';
import '../features/profile/presentation/profile_screen.dart';
import '../features/playback/presentation/player_bar.dart';
import '../features/playback/presentation/playback_controller.dart';
import '../features/search/presentation/search_screen.dart';

class HomeShell extends ConsumerStatefulWidget {
  const HomeShell({super.key});

  @override
  ConsumerState<HomeShell> createState() => _HomeShellState();
}

class _HomeShellState extends ConsumerState<HomeShell> {
  int _index = 0;

  @override
  Widget build(BuildContext context) {
    final playback = ref.watch(playbackControllerProvider);
    final hasPlayer = playback.track != null;

    return Scaffold(
      body: Column(
        children: [
          Expanded(
            child: IndexedStack(
              index: _index,
              children: const [
                DiscoverScreen(),
                SearchScreen(),
                CatalogScreen(),
                PlaylistsScreen(),
                LibraryScreen(),
                ProfileScreen(),
              ],
            ),
          ),
          if (hasPlayer) const PlayerBar(),
        ],
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: _index,
        onTap: (i) => setState(() => _index = i),
        items: const [
          BottomNavigationBarItem(icon: Icon(Icons.home_filled), label: 'Home'),
          BottomNavigationBarItem(icon: Icon(Icons.search), label: 'Search'),
          BottomNavigationBarItem(icon: Icon(Icons.library_music), label: 'Catalog'),
          BottomNavigationBarItem(icon: Icon(Icons.queue_music), label: 'Playlists'),
          BottomNavigationBarItem(icon: Icon(Icons.favorite), label: 'Library'),
          BottomNavigationBarItem(icon: Icon(Icons.person), label: 'Profile'),
        ],
      ),
    );
  }
}
