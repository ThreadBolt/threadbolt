# ThreadBolt Framework

![Version](https://img.shields.io/badge/version-v1.2.0-blue.svg)
![Status](https://img.shields.io/badge/status-beta-orange.svg)

ThreadBolt is a convention-over-configuration web framework for Go, inspired by Spring Boot. It provides MVC architecture, built-in ORM support, dependency injection, and powerful CLI tools for rapid development of enterprise-grade web applications and APIs.

## 🚀 Features

- **MVC Architecture**: Clean separation of concerns with Controllers, Models, and Views
- **Built-in ORM**: GORM integration with support for PostgreSQL, MySQL, and SQLite
- **CLI Tools**: Powerful code generation and project management commands
- **Dependency Injection**: Simple DI container for services and components
- **Convention over Configuration**: Enforced project structure and best practices
- **Middleware Support**: Built-in and custom middleware chain support
- **Configuration Management**: YAML and environment variable support with Viper
- **Hot Reload**: Development server with automatic restart
- **Testing Support**: Built-in testing utilities and patterns

## 📦 Installation

### Prerequisites

- Go 1.21 or later
- Git

### Install ThreadBolt CLI

```bash
go install github.com/ThreadBolt/threadbolt/cmd/threadbolt@v1.2.0
```

Verify installation:
```bash
threadbolt --help
```

## 🏗️ Project Structure

ThreadBolt enforces a strict project structure to maintain consistency and best practices:

```
your-app/
├── cmd/               # CLI entry points
├── config/            # Configuration files (config.yaml)
├── controllers/       # MVC controllers
├── internal/          # Internal packages
│   ├── middleware/    # Custom middleware
│   └── services/      # Business logic services
├── models/            # ORM models with GORM tags
├── migrations/        # Database migration files
├── public/            # Static assets (CSS, JS, images)
├── routes/            # Route definitions
├── templates/         # View templates (HTML)
├── tests/             # Unit and integration tests
├── go.mod             # Go modules file
├── go.sum             # Go modules checksum
└── main.go            # Application entry point
```

## 🚀 Quick Start

### 1. Create a New Application

```bash
threadbolt new my-awesome-app
cd my-awesome-app
```

This creates a complete ThreadBolt application with the standard structure and basic configuration.

### 2. Run the Application

```bash
threadbolt run
```

Or using Go directly:
```bash
go run main.go
```

Your application will start on `http://localhost:8080`

### 3. Test the API

The generated application includes health check endpoints:

```bash
# Health check
curl http://localhost:8080/health

# API status
curl http://localhost:8080/api/v1/status
```

## 🛠️ CLI Commands

### Project Management

- `threadbolt new <app-name>` - Create a new ThreadBolt application
- `threadbolt run` - Start the development server with hot reload
- `threadbolt test` - Run all tests
- `threadbolt migrate` - Run database migrations

### Code Generation

- `threadbolt generate model <ModelName>` - Generate a new model with repository
- `threadbolt generate controller <ControllerName>` - Generate a new controller with CRUD operations

### Examples

```bash
# Generate a User model
threadbolt generate model User

# Generate a Product controller
threadbolt generate controller Product
```

## 📊 Models and ORM

ThreadBolt uses GORM for database operations. Models are Go structs with tags for database mapping.

### Creating a Model

```bash
threadbolt generate model User
```

This generates `models/user.go`:

```go
package models

import (
    "gorm.io/gorm"
)

type User struct {
    BaseModel
    Name  string `gorm:"not null" json:"name"`
    Email string `gorm:"uniqueIndex;not null" json:"email"`
    Age   int    `json:"age"`
}

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

// CRUD methods...
func (r *UserRepository) Create(user *User) error { ... }
func (r *UserRepository) GetByID(id uint) (*User, error) { ... }
func (r *UserRepository) GetAll() ([]User, error) { ... }
func (r *UserRepository) Update(user *User) error { ... }
func (r *UserRepository) Delete(id uint) error { ... }
```

### Base Model

All models inherit from `BaseModel` which provides:
- `ID` (primary key)
- `CreatedAt` timestamp
- `UpdatedAt` timestamp  
- `DeletedAt` soft delete support

### Relationships

Define relationships using GORM tags:

```go
type User struct {
    BaseModel
    Name    string    `json:"name"`
    Posts   []Post    `json:"posts"` // One-to-Many
    Profile Profile   `json:"profile"` // One-to-One
}

type Post struct {
    BaseModel
    Title  string `json:"title"`
    UserID uint   `json:"user_id"`
    User   User   `json:"user"` // Belongs-to
}
```

## 🎮 Controllers

Controllers handle HTTP requests and responses following MVC patterns.

### Generating a Controller

```bash
threadbolt generate controller User
```

This generates `controllers/user_controller.go` with full CRUD operations:

```go
type UserController struct {
    repo *models.UserRepository
}

func NewUserController(db *gorm.DB) *UserController {
    return &UserController{
        repo: models.NewUserRepository(db),
    }
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) { ... }
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) { ... }
func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) { ... }
func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) { ... }
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) { ... }
```

### Custom Controllers

Create controllers manually in the `controllers/` directory:

```go
package controllers

import (
    "encoding/json"
    "net/http"
)

type CustomController struct {
    service *services.CustomService
}

func (c *CustomController) CustomEndpoint(w http.ResponseWriter, r *http.Request) {
    // Your logic here
    response := map[string]string{"message": "Hello from threadbolt!"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

## 🛣️ Routing

Define routes in `routes/routes.go`:

```go
package routes

import (
    "your-app/controllers"
    "github.com/gorilla/mux"
    "github.com/ThreadBolt/threadbolt/pkg/framework"
)

func SetupRoutes(app *framework.App) {
    // Get database from DI container
    db, _ := app.Container.Get("db")
    
    // Initialize controllers
    userController := controllers.NewUserController(db.(*gorm.DB))
    
    // API routes
    api := app.Router.PathPrefix("/api/v1").Subrouter()
    
    // User routes
    users := api.PathPrefix("/users").Subrouter()
    users.HandleFunc("", userController.CreateUser).Methods("POST")
    users.HandleFunc("", userController.GetAllUsers).Methods("GET")
    users.HandleFunc("/{id}", userController.GetUser).Methods("GET")
    users.HandleFunc("/{id}", userController.UpdateUser).Methods("PUT")
    users.HandleFunc("/{id}", userController.DeleteUser).Methods("DELETE")
}
```

## ⚙️ Configuration

ThreadBolt uses Viper for configuration management with support for YAML files and environment variables.

### Configuration File (`config/config.yaml`)

```yaml
server:
  port: 8080
  host: localhost

database:
  driver: sqlite          # postgres, mysql, sqlite
  host: localhost
  port: 5432
  name: myapp_dev
  username: myuser
  password: mypassword
  sslmode: disable

logging:
  level: info
  format: text

environment: development
```

### Environment Variables

Environment variables are prefixed with `THREADBOLT_`:

```bash
export THREADBOLT_SERVER_PORT=3000
export THREADBOLT_DATABASE_DRIVER=postgres
export THREADBOLT_DATABASE_NAME=myapp_production
```

### Accessing Configuration

```go
port := app.Config.GetString("server.port")
dbDriver := app.Config.GetString("database.driver")
isProduction := app.Config.GetString("environment") == "production"
```

## 🗄️ Database Support

ThreadBolt supports multiple databases through GORM:

### SQLite (Default)
```yaml
database:
  driver: sqlite
  name: myapp.db
```

### PostgreSQL
```yaml
database:
  driver: postgres
  host: localhost
  port: 5432
  name: myapp_dev
  username: postgres
  password: password
  sslmode: disable
```

### MySQL
```yaml
database:
  driver: mysql
  host: localhost
  port: 3306
  name: myapp_dev
  username: root
  password: password
```

## 🔧 Services and Dependency Injection

Services contain business logic and can be injected into controllers.

### Creating a Service

Create services in `internal/services/`:

```go
package services

import (
    "your-app/models"
    "gorm.io/gorm"
)

type UserService struct {
    repo *models.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
    return &UserService{
        repo: models.NewUserRepository(db),
    }
}

func (s *UserService) CreateUser(name, email string, age int) (*models.User, error) {
    // Business logic here
    user := &models.User{Name: name, Email: email, Age: age}
    return user, s.repo.Create(user)
}
```

### Dependency Injection

Register services in the DI container:

```go
// In main.go or initialization code
userService := services.NewUserService(db)
app.Container.Register("userService", userService)

// In controllers
type UserController struct {
    userService *services.UserService `inject:"userService"`
}

// Inject dependencies
app.Container.Inject(&userController)
```

## 🛡️ Middleware

ThreadBolt supports middleware chains for cross-cutting concerns.

### Built-in Middleware

CORS middleware is included in `internal/middleware/cors.go`:

```go
func CORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

### Custom Middleware

Create custom middleware:

```go
package middleware

import (
    "log"
    "net/http"
    "time"
)

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

### Applying Middleware

```go
// In routes/routes.go
import "your-app/internal/middleware"

func SetupRoutes(app *framework.App) {
    // Apply middleware
    app.Router.Use(middleware.CORS)
    app.Router.Use(middleware.Logger)
    
    // Define routes...
}
```

## 🧪 Testing

ThreadBolt provides built-in testing support and utilities.

### Running Tests

```bash
# Run all tests
threadbolt test

# Run with verbose output
go test ./... -v

# Run specific package
go test ./controllers -v
```

### Test Structure

Tests are located in the `tests/` directory:

```go
package tests

import (
    "testing"
    "your-app/models"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&models.User{})
    return db
}

func TestUserRepository_Create(t *testing.T) {
    db := setupTestDB()
    repo := models.NewUserRepository(db)
    
    user := &models.User{
        Name:  "Test User",
        Email: "test@example.com",
        Age:   25,
    }
    
    err := repo.Create(user)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    
    if user.ID == 0 {
        t.Error("Expected user ID to be set")
    }
}
```

### Integration Tests

Test HTTP endpoints:

```go
func TestUserController_CreateUser(t *testing.T) {
    db := setupTestDB()
    controller := controllers.NewUserController(db)
    
    payload := `{"name":"John Doe","email":"john@example.com","age":30}`
    req, _ := http.NewRequest("POST", "/users", strings.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    
    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(controller.CreateUser)
    handler.ServeHTTP(rr, req)
    
    if status := rr.Code; status != http.StatusCreated {
        t.Errorf("Expected status %v, got %v", http.StatusCreated, status)
    }
}
```

## 🚀 Deployment

### Building for Production

```bash
# Build binary
go build -o myapp main.go

# Build with optimizations
go build -ldflags="-s -w" -o myapp main.go
```

### Environment Configuration

Set production environment variables:

```bash
export THREADBOLT_ENVIRONMENT=production
export THREADBOLT_DATABASE_DRIVER=postgres
export THREADBOLT_DATABASE_HOST=prod-db-host
export THREADBOLT_DATABASE_NAME=myapp_production
export THREADBOLT_SERVER_PORT=8080
```

### Docker Support

Create a `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/public ./public

EXPOSE 8080
CMD ["./main"]
```

## 📚 Best Practices

### Project Organization

1. **Follow the Standard Structure**: Never deviate from the enforced directory structure
2. **One Model Per File**: Keep models focused and single-purpose
3. **Thin Controllers**: Business logic should be in services, not controllers
4. **Repository Pattern**: Use repositories for data access abstraction

### Code Quality

1. **Use Dependency Injection**: Avoid global variables and tight coupling
2. **Handle Errors Properly**: Always check and handle errors appropriately
3. **Write Tests**: Maintain good test coverage for models and controllers
4. **Use Meaningful Names**: Follow Go naming conventions

### Performance

1. **Use Database Indexes**: Add indexes for frequently queried fields
2. **Optimize Queries**: Use GORM preloading to avoid N+1 queries
3. **Connection Pooling**: Configure database connection pool settings
4. **Middleware Order**: Place expensive middleware after cheap ones

### Security

1. **Validate Input**: Always validate and sanitize user input
2. **Use HTTPS**: Enable TLS in production environments  
3. **Environment Variables**: Never commit secrets to version control
4. **CORS Configuration**: Properly configure CORS for your needs

## 🔍 Examples

### Complete CRUD Example

Here's a complete example of a blog post feature:

#### 1. Generate the Model

```bash
threadbolt generate model Post
```

#### 2. Update the Model (`models/post.go`)

```go
type Post struct {
    BaseModel
    Title   string `gorm:"not null" json:"title"`
    Content string `gorm:"type:text" json:"content"`
    UserID  uint   `json:"user_id"`
    User    User   `json:"user"`
}
```

#### 3. Generate the Controller

```bash
threadbolt generate controller Post
```

#### 4. Add Routes (`routes/routes.go`)

```go
// Post routes
posts := api.PathPrefix("/posts").Subrouter()
posts.HandleFunc("", postController.CreatePost).Methods("POST")
posts.HandleFunc("", postController.GetAllPosts).Methods("GET")
posts.HandleFunc("/{id}", postController.GetPost).Methods("GET")
posts.HandleFunc("/{id}", postController.UpdatePost).Methods("PUT")
posts.HandleFunc("/{id}", postController.DeletePost).Methods("DELETE")
```

#### 5. Test the API

```bash
# Create a post
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -d '{"title":"My First Post","content":"This is the content","user_id":1}'

# Get all posts
curl http://localhost:8080/api/v1/posts

# Get specific post
curl http://localhost:8080/api/v1/posts/1
```

## 🤝 Contributing

We welcome contributions! Please see [DEVELOPMENT.md](DEVELOPMENT.md) for development setup and contribution guidelines.

## 📄 License

ThreadBolt is released under the MIT License. See [LICENSE](LICENSE) for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/ThreadBolt/threadbolt/issues)
- **Discussions**: [GitHub Discussions](https://github.com/ThreadBolt/threadbolt/discussions)

---

**Happy coding with ThreadBolt! 🎉**