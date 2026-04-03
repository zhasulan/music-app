import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

class AppTheme {
  static const Color primary = Color(0xFF8F7BFF);
  static const Color secondary = Color(0xFF1ED760);
  static const Color surface = Color(0xFF121212);
  static const Color surfaceAlt = Color(0xFF1E1E22);
  static const Color textPrimary = Colors.white;
  static const Color textSecondary = Color(0xFFB9B9C5);

  static ThemeData dark() {
    final base = ThemeData.dark(useMaterial3: true);
    final colorScheme = ColorScheme.fromSeed(
      seedColor: primary,
      brightness: Brightness.dark,
      primary: primary,
      secondary: secondary,
      surface: surface,
      background: surface,
    );
    return base.copyWith(
      colorScheme: colorScheme,
      scaffoldBackgroundColor: surface,
      appBarTheme: const AppBarTheme(
        backgroundColor: surface,
        elevation: 0,
        centerTitle: false,
        titleTextStyle: TextStyle(fontSize: 20, fontWeight: FontWeight.w700, color: textPrimary),
      ),
      textTheme: GoogleFonts.interTextTheme(base.textTheme).apply(
        bodyColor: textPrimary,
        displayColor: textPrimary,
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          backgroundColor: primary,
          foregroundColor: Colors.white,
          minimumSize: const Size.fromHeight(48),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        ),
      ),
      filledButtonTheme: FilledButtonThemeData(
        style: FilledButton.styleFrom(
          backgroundColor: secondary,
          foregroundColor: Colors.black,
          minimumSize: const Size.fromHeight(48),
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: surfaceAlt,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
        labelStyle: const TextStyle(color: textSecondary),
      ),
      cardTheme: CardThemeData(
        color: surfaceAlt,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        elevation: 0,
        margin: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      ),
      listTileTheme: const ListTileThemeData(
        contentPadding: EdgeInsets.symmetric(horizontal: 12, vertical: 2),
        iconColor: textSecondary,
        textColor: textPrimary,
        subtitleTextStyle: TextStyle(color: textSecondary),
      ),
      dialogTheme: DialogThemeData(
        backgroundColor: surfaceAlt,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        titleTextStyle: const TextStyle(fontSize: 18, fontWeight: FontWeight.w700, color: textPrimary),
        contentTextStyle: const TextStyle(color: textPrimary),
      ),
      bottomSheetTheme: const BottomSheetThemeData(
        backgroundColor: surfaceAlt,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.vertical(top: Radius.circular(16))),
      ),
      dividerTheme: const DividerThemeData(color: Colors.white10, thickness: 1),
      progressIndicatorTheme: const ProgressIndicatorThemeData(color: secondary),
      bottomNavigationBarTheme: const BottomNavigationBarThemeData(
        backgroundColor: surfaceAlt,
        selectedItemColor: secondary,
        unselectedItemColor: textSecondary,
        type: BottomNavigationBarType.fixed,
      ),
      snackBarTheme: const SnackBarThemeData(behavior: SnackBarBehavior.floating, backgroundColor: surfaceAlt),
    );
  }
}
