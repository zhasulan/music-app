# Browser Smoke Flow (Phase 1)

1. Start backend  
   ```bash
   make up
   curl http://127.0.0.1:8080/health
   ```
2. Run Flutter Web (Chrome)  
   ```bash
   make flutter-run-web
   ```
3. In the app (Chrome):
   - Register a new account
   - Login
   - Open Profile tab (edit/save name and country)
   - Open Catalog tab (view tracks)
   - Logout via the AppBar icon

Notes
- API_BASE_URL for web defaults to `http://127.0.0.1:8080` via dart-define.
- If backend ports are busy, adjust in `deploy/docker-compose/docker-compose.yml` and pass a matching `API_BASE_URL`.
