import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/models.dart';
import '../../../core/token_storage.dart';
import '../../../core/api_client.dart';
import '../data/auth_repository.dart';

class AuthState {
  final bool loading;
  final UserModel? user;
  final String? error;
  final bool authenticated;

  const AuthState({this.loading = false, this.user, this.error, this.authenticated = false});

  AuthState copyWith({bool? loading, UserModel? user, String? error, bool? authenticated}) {
    return AuthState(
      loading: loading ?? this.loading,
      user: user ?? this.user,
      error: error,
      authenticated: authenticated ?? this.authenticated,
    );
  }
}

class AuthController extends StateNotifier<AuthState> {
  AuthController(this._repo, this._storage) : super(const AuthState());

  final AuthRepository _repo;
  final TokenStorage _storage;

  Future<void> login(String email, String password) async {
    state = state.copyWith(loading: true, error: null);
    try {
      final (user, _, __) = await _repo.login(email, password);
      state = AuthState(user: user, loading: false, authenticated: true);
    } catch (e) {
      state = AuthState(error: 'Login failed', loading: false, authenticated: false);
    }
  }

  Future<void> register(String email, String password) async {
    state = state.copyWith(loading: true, error: null);
    try {
      final (user, _, __) = await _repo.register(email, password);
      state = AuthState(user: user, loading: false, authenticated: true);
    } catch (e) {
      state = AuthState(error: 'Register failed', loading: false, authenticated: false);
    }
  }

  Future<void> logout() async {
    await _repo.logout();
    await _storage.clear();
    state = const AuthState();
  }

  Future<void> loadFromStorage() async {
    final (access, _) = await _storage.loadTokens();
    state = AuthState(authenticated: access != null, user: state.user);
  }
}

final authControllerProvider = StateNotifierProvider<AuthController, AuthState>((ref) {
  final repo = ref.read(authRepositoryProvider);
  final storage = ref.read(tokenStorageProvider);
  return AuthController(repo, storage);
});
