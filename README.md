# Go RESTful API Boilerplate

A production-ready, enterprise-grade RESTful API built with Go, following **Clean Architecture** principles. This boilerplate provides a robust foundation for building scalable and maintainable web services.

## 🚀 Features

- **Clean Architecture**: Decoupled layers (Domain, Usecase, Repository, Delivery) for better maintainability and testability.
- **Web Framework**: [Fiber v2](https://gofiber.io/) - An Express-inspired web framework for Go, designed for high performance.
- **ORM**: [GORM](https://gorm.io/) - Fantastic ORM library for Go, providing a rich set of features.
- **Database**: PostgreSQL support out of the box.
- **Authentication**: [PASETO](https://paseto.io/) (Platform-Agnostic Security Tokens) - A safer alternative to JWT.
- **Configuration**: [Viper](https://github.com/spf13/viper) - Complete configuration solution for Go applications.
- **Logging**: [Zap](https://github.com/uber-go/zap) - Blazing fast, structured, leveled logging.
- **Validation**: [Validator v10](https://github.com/go-playground/validator) - Go Struct and Field validation.
- **Documentation**: [Swagger](https://swagger.io/) (via [swag](https://github.com/swaggo/swag)) - Automatically generated API documentation.
- **Dependency Injection**: Manual dependency injection with a centralized Bootstrap configuration.
- **Graceful Shutdown**: Handles OS signals for clean application termination.
- **Containerization**: Docker and Docker Compose support.

## 🛠 Tech Stack

- **Language**: Go 1.25+
- **Infrastructure**: Fiber, GORM, PostgreSQL, Zap, Viper, Paseto
- **Dev Tools**: Air (Live Reload), Swag (Swagger generation)

## 📁 Project Structure

```text
.
├── cmd/api             # Application entry point
├── internal/
│   ├── config          # Application bootstrap and DI
│   ├── delivery/http   # HTTP layer (Handlers, Routes, Middleware)
│   ├── domain/         # Core business logic (Entities, Models, Interfaces)
│   ├── infrastructure/ # External tools (DB, Logger, Validator, etc.)
│   ├── repository/     # Data access layer
│   └── usecase/        # Business logic implementation
├── docs/               # Swagger documentation files
├── test/               # Integration and Unit tests
├── .env.example        # Environment variable template
├── docker-compose.yaml # Docker orchestration
└── go.mod              # Go dependencies
```

## 🏁 Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL
- Docker & Docker Compose (optional)

### Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/yoviarpauzi/go-restful-api-boilerplate.git
    cd go-restful-api-boilerplate
    ```

2.  **Setup Environment Variables**:
    Copy the example environment file and update it with your credentials:
    ```bash
    cp .env.example .env
    ```

3.  **Download Dependencies**:
    ```bash
    go mod download
    ```

### Running the Application

#### Local Development
```bash
go run cmd/api/main.go
```

#### Using Docker
```bash
docker-compose up -d
```

## 📚 API Documentation

The API documentation is powered by Swagger. Once the server is running, you can access it at:
`http://localhost:8080/swagger/index.html`

### Core Endpoints

| Method | Endpoint | Description | Auth Required |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/health` | Health check endpoint | No |
| `POST` | `/api/v1/auth/register` | Register a new user | No |
| `POST` | `/api/v1/auth/login` | Login and get tokens | No |
| `POST` | `/api/v1/auth/reset-password` | Reset forgotten password | No |
| `POST` | `/api/v1/auth/change-password` | Change current password | Yes |
| `GET` | `/api/v1/users` | Get all users | Yes |
| `GET` | `/api/v1/users/:id` | Get user by ID | Yes |
| `PUT` | `/api/v1/users/:id` | Update user by ID | Yes |
| `DELETE` | `/api/v1/users/:id` | Delete user by ID | Yes |

### Re-generating Swagger Docs
If you've made changes to the API annotations, re-generate the docs using:
```bash
swag init -g cmd/api/main.go
```

## 🧪 Testing

### Running Tests
To run all tests:
```bash
go test ./...
```

### Integration Tests
Integration tests are located in `test/integration`. They require a running database (configured in `.env`).

## 🔐 Authentication

This project uses **PASETO v4 Local** tokens for authentication.
- **Access Token**: Returned in the login/register response body.
- **Refresh Token**: Set as an `HttpOnly` cookie for secure session management.

## 📜 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
