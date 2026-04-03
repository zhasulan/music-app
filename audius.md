# Log In with Audius - краткая документация

## Что это
`Log In with Audius` позволяет:
- получить профиль пользователя Audius
- при необходимости выполнять действия от его имени
- не запрашивать у пользователя пароль Audius напрямую

Реализация основана на:
- OAuth 2.0 Authorization Code Flow with PKCE

Важно:
- backend и client secret не обязательны при использовании SDK

---

## 1. Получить API Key

Нужно создать developer app и получить API key.

Где:
- `audius.co/settings -> Developer Apps`
- `api.audius.co/plans`

---

## 2. Зарегистрировать Redirect URI

Нужно добавить redirect URI, который будет использовать приложение.

Примеры:
- Web: `https://yourapp.com/callback`
- Mobile: `myapp://oauth/callback`
- Local dev: `http://localhost:PORT`

Audius валидирует redirect URI в каждом auth-запросе.

---

## 3. Инициализация SDK

```ts
import { sdk } from '@audius/sdk'

const audiusSdk = sdk({
  appName: 'My App',
  apiKey: 'YOUR_API_KEY',
  redirectUri: 'https://yourapp.com/callback',
})
React Native / Expo

Нужно установить зависимости:

npx expo install expo-web-browser @react-native-async-storage/async-storage
4. Логин пользователя
await audiusSdk.oauth.login({ scope: 'write' })

const user = await audiusSdk.oauth.getUser()

console.log('Signed in as', user.name)
Scope
read - если нужен только профиль и read-only доступ
write - если нужны действия от имени пользователя:
upload
favorite
другие write-операции

Важно:

write не дает доступ к DMs
write не дает доступ к wallet
5. Callback page (только web)

На callback page нужно вызвать:

const audiusSdk = sdk({ appName: 'My App', apiKey: 'YOUR_API_KEY' })

await audiusSdk.oauth.handleRedirect()
Что делает handleRedirect()

Поддерживает оба сценария:

popup flow
full-page redirect
Особенности
в popup SDK передаст код в родительское окно и закроет popup
в full-page redirect SDK завершит token exchange
Mobile

На mobile handleRedirect() вызывается автоматически внутри login().

6. Восстановление существующей сессии
if (await audiusSdk.oauth.isAuthenticated()) {
  const user = await audiusSdk.oauth.getUser()
  // restore UI
}
Примеры use cases
Write scope
загрузка треков в аккаунт пользователя
сохранение треков в библиотеку пользователя
Read-only scope
логин/регистрация через Audius
привязка аккаунта Audius к твоему приложению
получение данных пользователя
проверка, является ли пользователь verified artist
Ограничение

Этот flow:

не управляет сессией логина внутри твоего приложения автоматически
Manual OAuth Implementation

Если SDK использовать нельзя, можно реализовать OAuth 2.0 Authorization Code Flow with PKCE вручную.

Базовый API:

https://api.audius.co/v1
Шаг 1. Сгенерировать PKCE параметры

Нужно создать:

code_verifier
code_challenge = BASE64URL(SHA256(code_verifier))
state

Пример:

async function generatePkce() {
  const array = new Uint8Array(32)
  crypto.getRandomValues(array)

  const codeVerifier = btoa(String.fromCharCode(...array))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '')

  const encoded = new TextEncoder().encode(codeVerifier)
  const digest = await crypto.subtle.digest('SHA-256', encoded)

  const codeChallenge = btoa(String.fromCharCode(...new Uint8Array(digest)))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '')

  return { codeVerifier, codeChallenge }
}
Шаг 2. Открыть consent screen

Endpoint:

https://api.audius.co/v1/oauth/authorize

Пример:

https://api.audius.co/v1/oauth/authorize
  ?response_type=code
  &scope=read
  &api_key=YOUR_API_KEY
  &redirect_uri=https://mydemoapp.com/callback
  &state=YOUR_STATE
  &code_challenge=YOUR_CODE_CHALLENGE
  &code_challenge_method=S256
Обязательные параметры
response_type=code
scope=read|write
api_key=YOUR_API_KEY
redirect_uri=...
state=...
code_challenge=...
code_challenge_method=S256
Дополнительные параметры
response_mode=fragment|query
display=popup|fullScreen
Требования к redirect_uri
только http или https
raw IP запрещен, кроме localhost
нельзя использовать #
нельзя userinfo
нельзя path traversal
Шаг 3. Получить authorization code

После подтверждения Audius редиректит на callback URL.

Пример:

https://mydemoapp.com/callback#code=AUTH_CODE&state=YOUR_STATE

Нужно проверить:

state в ответе должен совпадать с отправленным

Если не совпадает:

прервать flow
возможна CSRF-атака
Шаг 4. Обменять code на токены

Endpoint:

POST https://api.audius.co/v1/oauth/token
Content-Type: application/json

Body:

{
  "grant_type": "authorization_code",
  "code": "AUTH_CODE",
  "code_verifier": "YOUR_CODE_VERIFIER",
  "client_id": "YOUR_API_KEY",
  "redirect_uri": "https://mydemoapp.com/callback"
}

Успешный ответ:

{
  "access_token": "...",
  "refresh_token": "..."
}

Что делать:

сохранить access_token
сохранить refresh_token
Шаг 5. Получить профиль пользователя

Endpoint:

GET https://api.audius.co/v1/me
Authorization: Bearer ACCESS_TOKEN

Пример ответа:

{
  "userId": 123,
  "name": "Artist Name",
  "handle": "artist_handle",
  "verified": true,
  "profilePicture": {
    "150x150": "...",
    "480x480": "...",
    "1000x1000": "...",
    "mirrors": []
  }
}
Шаг 6. Обновить access token

Endpoint:

POST https://api.audius.co/v1/oauth/token
Content-Type: application/json

Body:

{
  "grant_type": "refresh_token",
  "refresh_token": "YOUR_REFRESH_TOKEN",
  "client_id": "YOUR_API_KEY"
}

Ответ:

новый access_token
новый refresh_token
Шаг 7. Отозвать токен (logout)

Endpoint:

POST https://api.audius.co/v1/oauth/revoke
Content-Type: application/json

Body:

{
  "token": "YOUR_REFRESH_TOKEN",
  "client_id": "YOUR_API_KEY"
}

После этого:

удалить токены локально
даже если revoke вернул ошибку, локально токены все равно нужно очистить
Ключевые выводы
Когда нужен OAuth

Используй OAuth, если нужно:

логинить пользователя через Audius
получать профиль текущего пользователя
делать действия от имени пользователя
Когда OAuth не нужен

Если нужно только:

искать музыку
получать карточки треков
стримить музыку

то обычно хватает обычного REST API без user login.

Основные OAuth endpoints
Авторизация
GET /v1/oauth/authorize
Обмен code на token
POST /v1/oauth/token
Обновление token
POST /v1/oauth/token
Профиль текущего пользователя
GET /v1/me
Logout / revoke
POST /v1/oauth/revoke

Могу следом сделать **второй MD-блок именно под твой кейс** - без OAuth, только:
- поиск музыки
- стрим
- provider integration для Codex.