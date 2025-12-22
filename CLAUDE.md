# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based e-commerce shop management system with the following key components:
- **Shop Management**: Create, manage shops and products
- **Order Processing**: Handle order creation and status updates
- **User Authentication**: JWT-based auth with role-based access
- **Search Functionality**: Product and order search capabilities
- **Comment System**: Shop comments and reviews

## Development Commands

### Running the Application
```bash
# Start the application server
go run ./cmd

# The server runs on http://localhost:8080
```

### Testing
```bash
# Run all unit tests
go test ./test -v

# Run specific tests
go test ./test -v -run TestCreateShop
go test ./test -v -run TestShop*

# Run integration tests
cd test && go run integration_test.go
```

### Dependency Management
```bash
# Install dependencies
go mod tidy

# Generate Wire dependency injection code
go generate ./cmd
```

### Docker Services
```bash
# Start required services (PostgreSQL + Redis)
docker-compose up -d

# Stop services
docker-compose down
```

## Architecture

### Project Structure
```
├── cmd/                 # Application entry points
│   ├── main.go         # Main application bootstrap
│   ├── api.go          # HTTP handlers and routing
│   └── wire.go         # Dependency injection setup
├── internal/           # Private application code
│   ├── Auth/           # Authentication services
│   ├── Comment/        # Comment system
│   ├── Coordinator/    # Order coordination logic
│   ├── Order/          # Order management
│   ├── Shop/           # Shop and product management
│   ├── User/           # User management
│   ├── search/         # Search functionality
│   └── config/         # Configuration management
├── pkg/                # Public library code
│   ├── database/       # Database utilities
│   ├── middleware/     # HTTP middleware (JWT, CORS)
│   └── utils/          # Utility functions
├── test/               # Test files and integration tests
└── app.yaml           # Application configuration
```

### Dependency Injection (Wire)
The project uses Google Wire for dependency injection. Key points:
- `cmd/wire.go` defines the injection setup
- Run `go generate ./cmd` to regenerate `wire_gen.go`
- Each module has its own `ProviderSet` in `wire_set.go` files

### Database Layer
- **ORM**: GORM with PostgreSQL
- **Migrations**: Auto-migration runs on startup via `provideDB()` in `cmd/wire.go`
- **Models**: Located in each module's `Model.go` files

### API Versioning
- **v0**: Public endpoints (auth, user registration)
- **v1**: Order management (JWT protected)
- **v2**: Shop/product management with role-based access
- **v3**: Comment system

### Authentication & Authorization
- JWT-based authentication using middleware.JWTAuthMiddleware()
- Role-based access for merchant operations via middleware.MerchantAuthMiddleware()
- JWT secret configurable in app.yaml or via SHOP_JWT_SECRET env var

### Configuration
- Primary config: `app.yaml`
- Environment variables prefixed with `SHOP_` (e.g., `SHOP_DATABASE_HOST`)
- Config struct defined in `internal/config/config.go`

## Development Notes

### Database Schema
The system auto-migrates the following models:
- Users, Shops, Products, Categories
- Orders, OrderItems
- Comments

### Redis Usage
- Session management via middleware.RedisStore
- Configured in app.yaml or SHOP_REDIS_* env vars

### Search Implementation
- Separate search services in `internal/search/`
- Product and order search capabilities

### Wire Code Generation
When adding new dependencies:
1. Update provider sets in respective `wire_set.go` files
2. Run `go generate ./cmd` to regenerate injection code
3. Ensure new providers are added to `InitializeApp()` in `cmd/wire.go`

### Testing Strategy
- Unit tests in `test/api_test.go`
- Integration tests in `test/integration_test.go`
- Postman collection available in `test/postman_collection.json`

### CORS Configuration
Currently configured to allow all origins in development mode. Adjust cors.Config in `cmd/api.go` for production.