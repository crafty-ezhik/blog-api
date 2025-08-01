# 🚀 Блог с авторизацией через JWT

**API для блога с поддержкой пользовательских статей, комментариев и безопасной системы авторизации через JWT.**

## 📌 Основные возможности

- Регистрация и вход пользователей
- Управление токенами: `access` и `refresh`
- Полный контроль над сессиями (включая выход из всех устройств)
- Управление профилем пользователя
- Создание, редактирование и удаление статей
- Комментирование статей с возможностью фильтрации
- Поддержка тестирования и мока данных через `gomock`

---

## 🛠️ Используемые технологии

| Категория        | Инструменты / Библиотеки                                                                 |
|------------------|------------------------------------------------------------------------------------------|
| Веб-фреймворк     | [Fiber](https://gofiber.io/)                                                             |
| ORM              | [GORM](https://gorm.io/)                                                                 |
| База данных      | PostgreSQL                                                                               |
| Хранение чёрного списка JWT | [Redis](https://github.com/redis/go-redis)                                                                                |
| Шифрование       | [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)                                                                                    |
| Аутентификация   | [JWT-Go](https://github.com/dgrijalva/jwt-go)                                            |
| Конфигурация     | [Viper](https://github.com/spf13/viper)                                                  |
| Логгирование     | [Zap](https://github.com/uber-go/zap)                                                    |
| Валидация        | [go-playground/validator](https://github.com/go-playground/validator)                    | |
| Тестирование     | [Testify](https://github.com/stretchr/testify), [GoMock](https://github.com/golang/mock) |

---

## 🔐 Система безопасности

- **JWT Tokens**:
  - `Access Token`: краткосрочный, используется для доступа к защищённым ресурсам.
  - `Refresh Token`: долгосрочный, хранится в `HttpOnly` cookie с параметром `SameSite=Lax`.
- **Чёрный список (Redis)**: все отозванные `refresh` токены сохраняются до истечения их срока действия.
- **Версионирование токенов**: каждый новый refresh увеличивает версию токена, чтобы предотвратить replay-атаки.
- **Безопасное хранение паролей**: использование `bcrypt` для хэширования.
- **CSRF Protection**: рекомендуется использовать middleware или проверку `SameSite` + `Origin`.

___

## 🧪 Тестирование

- Все обработчики покрыты unit-тестами с использованием библиотеки `testify`.
- Для имитации зависимостей используется `gomock`.
- Mock'и генерируются автоматически и используются при тестировании бизнес-логики.

---

## 📦 Структура API

### 1. Авторизация

| Метод | Путь             | Описание                   |
|-------|------------------|----------------------------|
| POST  | `/auth/register` | Регистрация пользователя   |
| POST  | `/auth/login`    | Авторизация                |
| POST  | `/auth/logout`   | Выход из текущей сессии     |
| POST  | `/auth/refresh`  | Обновление токенов         |

---

### 2. Пользователи

| Метод | Путь                         | Описание                          |
|-------|------------------------------|-----------------------------------|
| GET   | `/api/users/me`              | Получение информации о себе      |
| GET   | `/api/users/my/posts`        | Получение своих статей            |
| GET   | `/api/users/:id/posts`       | Получение статей пользователя     |

---

### 3. Статьи (Posts)

| Метод | Путь                         | Описание                          |
|-------|------------------------------|-----------------------------------|
| POST  | `/api/posts`                 | Создание статьи                  |
| GET   | `/api/posts`                 | Получение всех статей             |
| GET   | `/api/posts/:id`             | Получение конкретной статьи       |
| PUT   | `/api/posts/:id`             | Обновление статьи                 |
| DELETE| `/api/posts/:id`             | Удаление статьи                   |

---

### 4. Комментарии

| Метод | Путь                                   | Описание                            |
|-------|----------------------------------------|-------------------------------------|
| GET   | `/api/posts/:id/comments`              | Получение всех комментариев к статье |
| POST  | `/api/posts/:id/comments`              | Добавление комментария               |
| PUT   | `/api/posts/:id/comments/:commentId`   | Обновление комментария               |
| DELETE| `/api/posts/:id/comments/:commentId`   | Удаление комментария                 |
| GET   | `/api/users/my/posts/:postId/comments` | Получение своих комментариев к статье |
| GET   | `/api/users/:id/posts/:postId/comments`| Получение комментариев к статье по ID пользователя |

---

## 🧰 Настройка окружения

Создайте файлы конфигурации в папке `configs`:

```bash
configs/
├── dev.yaml
├── prod.yaml
└── test.yaml
```

Пример содержимого `.env` файла:

```env
# Окружение
APP_ENV=prod

# Сервер
APP_SERVER_PORT=8080
APP_SERVER_MODE=debug

# База данных
APP_DATABASE_HOST=localhost
APP_DATABASE_PORT=5432
APP_DATABASE_USERNAME=postgres
APP_DATABASE_PASSWORD=pass
APP_DATABASE_DATABASE=db_name
APP_DATABASE_SSLMODE=disable

# JWT
APP_JWT_SIGNING_KEY=my_signing_key
APP_JWT_SECRET_KEY=my_encryption_key
APP_JWT_ACCESS_TTL=5m
APP_JWT_REFRESH_TTL=48h

# Redis 
APP_REDIS_HOST=localhost
APP_REDIS_PORT=6379
```

> ⚠️ **Важно**: переменная `APP_ENV` должна совпадать с названием соответствующего YAML-файла (`dev.yaml`, `prod.yaml`).

---

## 📚 Дополнительно

- Валидация входящих данных реализована через универсальные типы (`generic`) и библиотеку `go-playground/validator`.

---

## 🧑‍💻 Разработка

Для разработки и тестирования используйте режим `dev`, который позволяет быстро перезапускать сервер и видеть логи в реальном времени.

---

## ✅ Что можно улучшить в будущем
- Добавить документацию в Swagger
- Добавить rate limiting для защиты от DDoS и злоупотребления API.
- Реализовать email-подтверждение регистрации.
- Добавить пагинацию и фильтрацию для списков статей и комментариев.
- Добавить CI/CD pipeline (GitHub Actions, GitLab CI и др.)

---