<div align="center">
  <img src="./docs/ledgerly.png" width="100%" alt="Ledgerly" />
</div>

## Tech Stack

| Library / Tool      | Purpose                      |
| ------------------- | ---------------------------- |
| **Go**              | Programming language         |
| **Gin**             | HTTP web framework           |
| **GORM**            | ORM for database operations  |
| **SQLite**          | Embedded database            |
| **JWT**             | Authentication tokens        |
| **bcrypt**          | Password hashing             |
| **Swagger**         | API documentation            |
| **godotenv**        | Environment configuration    |


## Up and Running

```bash
# Clone repo
git clone <repo-url>

cd ledgerly

# Install dependencies
go mod tidy

# Configure environment
cp .env.example .env

# Run in development mode
go run cmd/main.go

# Build for production
go build -o ledgerly cmd/main.go

# Run production binary
./ledgerly
```

---

## API Documentation

Swagger UI available at: **`http://localhost:8080/swagger/index.html`**

### Endpoints

| Method | Endpoint                    | Description              | Auth |
| ------ | --------------------------- | ------------------------ | ---- |
| POST   | `/auth/login`               | User login               | ❌   |
| POST   | `/petty-cash`               | Create transaction       | ✅   |
| GET    | `/petty-cash`               | List transactions        | ✅   |
| GET    | `/petty-cash/balance`       | Get balance              | ✅   |
| POST   | `/expenses`                 | Create expense           | ✅   |
| GET    | `/expenses`                 | List expenses            | ✅   |
| GET    | `/reports/expenses-summary` | Expense report           | ✅   |
| GET    | `/reports/petty-cash-summary` | Petty cash report      | ✅   |

---

## Data Models

- **User**: Authentication & profile
- **Expense**: Transaction records with categories
- **PettyCash**: Cash flow & balance tracking

---

## Docker

```bash
# Build image
docker build -t ledgerly .

# Run container
docker run -p 8080:8080 -e JWT_SECRET=your_secret ledgerly
```

---

## Architecture

```
ledgerly/
├── cmd/           # Application entrypoint
├── db/            # Database connection
├── docs/          # Swagger documentation
├── handlers/      # HTTP handlers
├── middleware/    # Auth, RBAC, Rate limiter
├── models/        # Data models & permissions
├── routes/        # Route definitions
└── services/      # Business logic
```