# Todo App - Clean Architecture

A Todo application built with Go, Gin, PostgreSQL, and JWT authentication following clean architecture principles.

## Features

- User registration and authentication with JWT
- Role-based access control (User and Admin roles)
- Users can CRUD their own todo lists
- Admins can view all users and todos
- Raw SQL queries (no ORM)
- Clean architecture with dependency injection

## API Endpoints

### Authentication
- `POST /api/v1/register` - User registration
- `POST /api/v1/login` - User login

### User Profile
- `GET /api/v1/profile` - Get user profile (requires auth)

### Todos
- `POST /api/v1/todos` - Create todo (requires auth)
- `GET /api/v1/todos` - Get user's todos (requires auth)
- `GET /api/v1/todos/:id` - Get specific todo (requires auth)
- `PUT /api/v1/todos/:id` - Update todo (requires auth)
- `DELETE /api/v1/todos/:id` - Delete todo (requires auth)
- `PATCH /api/v1/todos/:id/status` - Update todo status (requires auth)

### Admin Only
- `GET /api/v1/users` - Get all users (admin only)
- `GET /api/v1/admin/todos` - Get all todos (admin only)

## Project Structure

```
cmd/
  web/
    main.go                 # Application entry point
internal/
  config/
    config.go              # Configuration management
  delivery/
    http/
      middleware/
        auth.go            # JWT authentication middleware
      route/
        route.go           # Route definitions
      auth_handler.go      # Authentication handlers
      todo_handler.go      # Todo handlers
      user_handler.go      # User handlers
  entity/
    user.go               # User entity
    todo.go               # Todo entity
  model/
    model.go              # Request/response models
    converter/
      converter.go        # Entity-model converters
  repository/
    user_repository.go    # User database operations
    todo_repository.go    # Todo database operations
  usecase/
    auth_usecase.go       # Authentication business logic
    todo_usecase.go       # Todo business logic
    user_usecase.go       # User business logic
```

## Setup and Installation

### Prerequisites
- Go 1.21+
- PostgreSQL
- Git

### Environment Variables

Create a `.env` file or set the following environment variables:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todoapp
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key
JWT_ISSUER=todo-app

# Server
SERVER_PORT=8080
```

### Database Setup

1. Create a PostgreSQL database:
```sql
CREATE DATABASE todoapp;
```

2. The application will automatically create the required tables on startup.

### Installation and Running

1. Clone the repository:
```bash
git clone <repository-url>
cd otel-propagation-monorepo
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run cmd/web/main.go
```

The server will start on port 8080 (or the port specified in SERVER_PORT environment variable).

## Usage Examples

### Register a user
```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username": "john", "password": "password123"}'
```

### Create a todo (requires authentication)
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title": "Buy groceries", "description": "Milk, bread, eggs"}'
```

### Get user todos
```bash
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## User Roles

- **User**: Can register, login, and CRUD their own todos
- **Admin**: Can do everything users can do, plus view all users and todos

To create an admin user, you need to manually update the user's role in the database:

```sql
UPDATE users SET role = 'admin' WHERE username = 'admin_username';
```

## Architecture

This project follows Clean Architecture principles:

1. **Entity Layer**: Core business entities (`User`, `Todo`)
2. **Use Case Layer**: Business logic and rules
3. **Interface Adapters**: Controllers, presenters, and gateways
4. **Framework Layer**: Database, web framework, external interfaces

### Dependencies Flow
- **Entity** ← **Use Case** ← **Interface Adapters** ← **Frameworks**
- Inner layers don't depend on outer layers
- Dependency injection is used to invert dependencies

## Security

- Passwords are hashed using bcrypt
- JWT tokens are used for authentication
- Role-based access control
- Input validation on all endpoints
- SQL injection protection through parameterized queries

## Development

### Adding New Features

1. Define entities in `internal/entity/`
2. Create use cases in `internal/usecase/`
3. Implement repositories in `internal/repository/`
4. Add handlers in `internal/delivery/http/`
5. Update routes in `internal/delivery/http/route/`

### Testing

Run tests with:
```bash
go test ./...
```

## License

This project is licensed under the MIT License.