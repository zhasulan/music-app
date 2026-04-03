import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter/foundation.dart';

import 'config.dart';
import 'token_storage.dart';

final tokenStorageProvider = Provider<TokenStorage>((ref) {
  return TokenStorage(const FlutterSecureStorage());
});

final dioProvider = Provider<Dio>((ref) {
  final storage = ref.read(tokenStorageProvider);
  final dio = Dio(BaseOptions(baseUrl: apiBaseUrl));

  if (kDebugMode) {
    dio.interceptors.add(LogInterceptor(
      request: true,
      requestBody: true,
      responseBody: false,
      logPrint: (o) => debugPrint(o.toString()),
    ));
  }

  dio.interceptors.add(InterceptorsWrapper(
    onRequest: (options, handler) async {
      final (access, _) = await storage.loadTokens();
      if (access != null) {
        options.headers['Authorization'] = 'Bearer $access';
      }
      options.headers['Content-Type'] = 'application/json';
      handler.next(options);
    },
    onError: (e, handler) async {
      // Auto-clear tokens on 401 to avoid stale auth state.
      if (e.response?.statusCode == 401) {
        await storage.clear();
      }
      handler.next(e);
    },
  ));

  return dio;
});
