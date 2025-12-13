# Elite Voting Application

Приложение для голосования с авторизацией через Telegram Login Widget.

## Структура проекта

- `cmd/server` - основной сервер приложения
- `cmd/migrate` - утилита для миграций базы данных
- `cmd/test_auth` - утилита для тестирования авторизации
- `internal/config` - конфигурация приложения
- `internal/database` - подключение к базе данных
- `internal/models` - модели данных
- `internal/repository` - репозитории для работы с БД
- `internal/service` - бизнес-логика
- `internal/handler` - HTTP handlers
- `internal/middleware` - middleware для авторизации
- `internal/auth` - авторизация через Telegram Login Widget
- `migrations` - миграции базы данных

## Запуск

### Локально

1. Установите зависимости:
```bash
go mod download
```

2. Создайте файл `.env`:
```bash
DATABASE_URL=postgres://elite:elite@localhost:5432/elite?sslmode=disable
TELEGRAM_BOT_TOKEN=ваш_токен_от_BotFather
JWT_SECRET=ваш_секретный_ключ
PORT=8080
```

3. Запустите PostgreSQL (можно через docker-compose):
```bash
docker-compose up -d postgres
```

4. Примените миграции:
```bash
go run cmd/migrate/main.go -command=up
```

5. Запустите сервер:
```bash
go run cmd/server/main.go
```

### Docker Compose

```bash
docker-compose up -d
```

## API Endpoints

### POST /api/auth
Авторизация через Telegram Login Widget.

**Request:**
```json
{
  "id": 123456789,
  "first_name": "John",
  "last_name": "Doe",
  "username": "johndoe",
  "photo_url": "https://...",
  "auth_date": 1234567890,
  "hash": "abc123..."
}
```

**Response:**
```json
{
  "token": "jwt_token_here"
}
```

### POST /api/vote
Создание голоса (требует авторизации).

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "nominant_id": 1
}
```

**Response:**
```json
{
  "message": "Vote created successfully"
}
```

## Telegram Login Widget

### Настройка

1. Создайте бота через [@BotFather](https://t.me/BotFather)
2. Получите токен бота
3. Добавьте токен в `.env` как `TELEGRAM_BOT_TOKEN`

### Использование на сайте

См. пример в `example_login_widget.html`:

```html
<script async src="https://telegram.org/js/telegram-widget.js?22"
    data-telegram-login="YOUR_BOT_USERNAME"
    data-size="large"
    data-onauth="onTelegramAuth(user)"
    data-request-access="write">
</script>

<script>
function onTelegramAuth(user) {
    // Отправьте данные на ваш API
    fetch('http://localhost:8080/api/auth', {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            id: user.id,
            first_name: user.first_name,
            last_name: user.last_name || '',
            username: user.username || '',
            photo_url: user.photo_url || '',
            auth_date: user.auth_date,
            hash: user.hash
        })
    })
    .then(response => response.json())
    .then(data => {
        // Сохраните токен
        localStorage.setItem('token', data.token);
    });
}
</script>
```

## База данных

### Таблицы:
- `categories` - категории
- `nominants` - номинанты
- `users` - пользователи Telegram
- `votes` - голоса (с уникальным ограничением на tg_user_id + category_id)

## Тестирование

### Через HTML страницу

1. Откройте `example_login_widget.html` в браузере
2. Укажите имя вашего бота (без @)
3. Нажмите "Login with Telegram"
4. После авторизации получите токен

### Через тестовый скрипт

```bash
go run cmd/test_auth/main.go <id> <first_name> <last_name> <username> <photo_url> <auth_date> <hash> [server_url]
```

## Документация

- [Telegram Login Widget](https://core.telegram.org/widgets/login)
- [Проверка авторизации](https://core.telegram.org/widgets/login#checking-authorization)
