# 🚀 Проект: **API для блога с авторизацией и JWT**
## 📌 Функционал:

### 1. **Авторизация**
- Регистрация 
  - (`POST /api/auth/register`)
- Авторизация 
  - (`POST /api/auth/login`)
- Выход из системы 
  - (`POST /api/auth/logout`)
- Выход из всех сессий 
  - (`POST /api/auth/logout-all`)
- Обновление refresh токена 
  - (`POST /api/auth/refresh`)

### 2. **Пользователи**
- Получение профиля 
  - (`GET /api/users/me` + auth middleware)
- Получение своих постов 
  - (`GET /api/users/my/posts`)
- Получение постов по ID пользователя 
  - (`GET /api/users/:id/posts`)

### 3. **Статьи (Posts)**
- Создание статьи 
  - (`POST /api/posts`)
- Получение всех статей 
  - (`GET /api/posts`)
- Получение конкретной статьи 
  - (`GET /api/posts/:id`)
- Обновление статьи 
  - (`PUT /api/posts/:id`)
- Удаление статьи 
  - (`DELETE /api/posts/:id`)

### 4. **Комментарии**
- Получение всех комментариев к статье
  - (`GET /api/posts/:id/comments`)
- Создание комментария к посту 
  - (`POST /api/posts/:id/comments`)
- Обновление комментария 
  - (`PUT /api/posts/:id/comments/:commentId`)
- Удаление комментария 
  - (`DELETE /api/posts/:id/comments/:commentId`)
- Получение всех своих комментариев к статье 
  - (`GET /api/users/my/posts/:postId/comments`)
- Получение всех комментариев к статье по id пользователя 
  - (`GET /api/users/:id/posts/:postId/comments`)

## 🛠 Технологии и инструменты:

| Категория         | Инструменты / Библиотеки                        |
|------------------|-------------------------------------------------|
| Фреймворк        | [Fiber](https://github.com/gofiber/fiber) |
| ORM              | [GORM](https://gorm.io/)                        |
| БД               | PostgreSQL         |
| Чёрный список JWT | Redis|
| Аутентификация   | [JWT-Go](https://github.com/dgrijalva/jwt-go)   |
| Конфигурация     | [Viper](https://github.com/spf13/viper)        |
| Логгирование     | [Zap](https://github.com/uber-go/zap)          |
| Валидация        | [go-playground/validator](https://github.com/go-playground/validator) |
| Документация     | [swaggo/gin-swagger](https://github.com/swaggo/swag) |
| Тестирование     | testify, go test                                |

---

## ⚙️ Установка и запуск
### Docker-compose
```cmd
docker-compose -up build -d
```