import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/api_client.dart';
import '../../../core/models.dart';

class ProfileRepository {
  ProfileRepository(this._dio);
  final Dio _dio;

  Future<UserModel> me() async {
    final res = await _dio.get('/api/v1/users/me');
    final data = res.data['user'] as Map<String, dynamic>;
    return UserModel.fromJson(data);
  }

  Future<UserModel> update({required String name, required String country}) async {
    final res = await _dio.put('/api/v1/users/me', data: {
      'name': name,
      'country': country,
    });
    final data = res.data['user'] as Map<String, dynamic>;
    return UserModel.fromJson(data);
  }
}

final profileRepositoryProvider = Provider<ProfileRepository>((ref) {
  return ProfileRepository(ref.read(dioProvider));
});
