import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import '../../../core/models.dart';
import '../../../core/token_storage.dart';

class AuthRepository {
  AuthRepository(this._dio, this._storage);

  final Dio _dio;
  final TokenStorage _storage;

  Future<(UserModel, String, String)> login(String email, String password) async {
    final res = await _dio.post('/api/v1/auth/login', data: {
      'email': email,
      'password': password,
    });
    final data = res.data as Map<String, dynamic>;
    final user = UserModel.fromJson(data['user'] as Map<String, dynamic>);
    final access = data['access_token'] as String;
    final refresh = data['refresh_token'] as String;
    await _storage.saveTokens(access, refresh);
    return (user, access, refresh);
  }

  Future<(UserModel, String, String)> register(String email, String password) async {
    final res = await _dio.post('/api/v1/auth/register', data: {
      'email': email,
      'password': password,
    });
    final data = res.data as Map<String, dynamic>;
    final user = UserModel.fromJson(data['user'] as Map<String, dynamic>);
    final access = data['access_token'] as String;
    final refresh = data['refresh_token'] as String;
    await _storage.saveTokens(access, refresh);
    return (user, access, refresh);
  }

  Future<void> logout() async {
    await _storage.clear();
  }
}

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepository(ref.read(dioProvider), ref.read(tokenStorageProvider));
});
