# SSO Auth Project

Цей проєкт є системою автентифікації на основі Single Sign-On (SSO), реалізованою на мові програмування Go. Проєкт включає роботу з JWT для генерації токенів, управління користувачами та додатками, а також забезпечує взаємодію через gRPC.


## Основні Компоненти

### JWT

**Файл:** `internal/lib/jwt/jwt.go`

- Реалізує створення JWT токенів для користувачів.
- Основні методи:
  - `NewTokens`: Генерує новий JWT токен з інформацією про користувача, додаток та час дії токена.

### Auth

**Файл:** `internal/services/auth/auth.go`

- Забезпечує основні функції автентифікації та реєстрації користувачів.
- Основні методи:
  - `Login`: Виконує вхід користувача, перевіряє облікові дані та повертає JWT токен.
  - `Register`: Реєструє нового користувача, зберігає хеш паролю в базі даних.
  - `IsAdmin`: Перевіряє, чи є користувач адміністратором.

### Storage

**Файл:** `internal/storage/sqlite/sqlite.go`

- Відповідає за взаємодію з базою даних SQLite.
- Основні методи:
  - `SaveUser`: Зберігає інформацію про користувача в базі даних.
  - `User`: Повертає інформацію про користувача за його email.
  - `IsAdmin`: Перевіряє, чи є користувач адміністратором.
  - `App`: Повертає інформацію про додаток за його ID.

### gRPC Server

**Файли:**

- `internal/grpc/app/grpcapp.go`
- `internal/grpc/auth/server.go`

- Реалізує gRPC сервер для автентифікації.
- Основні сервіси:
  - `Login`: Вхід користувача через gRPC.
  - `Register`: Реєстрація нового користувача через gRPC.
  - `IsAdmin`: Перевірка ролі адміністратора для користувача через gRPC.

## Запуск Проєкту

1. Скомпілюйте проєкт:
    ```bash
    task run
    ```
2. Створити міграції:
    ```bash
    task migrate
    ```
3. Генерування з proto файлу:
    ```bash
    task generate
    ```

## Залежності

- Go 1.20+
- gRPC
- SQLite
- golang-jwt
З