# 🚀 Проект: **API для блога с авторизацией и JWT**
## 📌 Функционал:

### 1. **Пользователи**
- Регистрация (`POST /api/auth/register`)
- Авторизация (`POST /api/auth/login`)
- Получение профиля (`GET /api/users/me` + auth middleware)

### 2. **Статьи (Posts)**
- Создание статьи (`POST /api/posts`)
- Получение всех статей (`GET /api/posts`)
- Получение конкретной статьи (`GET /api/posts/:id`)
- Обновление статьи (`PUT /api/posts/:id`)
- Удаление статьи (`DELETE /api/posts/:id`)

### 3. **Комментарии**
- Добавление комментария к статье
- Получение комментариев по ID статьи

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