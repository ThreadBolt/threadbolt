package generator

const mainGoTemplate = `package main

import (
	"log"

	"github.com/ThreadBolt/threadbolt/pkg/framework"
)

func main() {
	app, err := framework.LoadApp()
	if err != nil {
		log.Fatalf("Failed to load application: %v", err)
	}

	port := app.Config.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting {{.AppName}} on port %s", port)
	if err := app.Start(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
`

const goModTemplate = `module {{.AppName}}

go 1.21

require (
	github.com/ThreadBolt/threadbolt v0.1.0
)
`

const configTemplate = `server:
  port: 8080
  host: localhost

database:
  driver: sqlite
  name: {{.AppName}}.db

logging:
  level: info
  format: text

environment: development
`

const routesTemplate = `package routes

import (
	"{{.AppName}}/controllers"
	"github.com/gorilla/mux"
	"github.com/ThreadBolt/threadbolt/pkg/framework"
)

func SetupRoutes(app *framework.App) {
	// Health check
	app.Router.HandleFunc("/health", controllers.HealthCheck).Methods("GET")

	// API routes
	api := app.Router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/status", controllers.StatusCheck).Methods("GET")
}
`

const healthControllerTemplate = `package controllers

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string    ` + "`json:\"status\"`" + `
	Timestamp time.Time ` + "`json:\"timestamp\"`" + `
	Service   string    ` + "`json:\"service\"`" + `
}

type StatusResponse struct {
	Message string ` + "`json:\"message\"`" + `
	Version string ` + "`json:\"version\"`" + `
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "{{.AppName}}",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func StatusCheck(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		Message: "{{.AppName}} API is running",
		Version: "1.0.0",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
`

const baseModelTemplate = `package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uint           ` + "`gorm:\"primarykey\" json:\"id\"`" + `
	CreatedAt time.Time      ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time      ` + "`json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`gorm:\"index\" json:\"-\"`" + `
}
`

const corsMiddlewareTemplate = `package middleware

import (
	"net/http"
)

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
`

const gitignoreTemplate = `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with go test -c
*.test

# Output of the go coverage tool
*.out

# Go workspace file
go.work

# Database files
*.db
*.sqlite

# Environment files
.env
.env.local

# IDE files
.vscode/
.idea/
*.swp
*.swo

# Logs
*.log

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db
`

const readmeTemplate = `# {{.AppName}}

A ThreadBolt application generated with the ThreadBolt framework.

## Getting Started

### Prerequisites

- Go 1.21 or later
- ThreadBolt CLI tool

### Running the Application

1. Install dependencies:
   ` + "```bash" + `
   go mod tidy
   ` + "```" + `

2. Run the application:
   ` + "```bash" + `
   threadbolt run
   ` + "```" + `

   Or use Go directly:
   ` + "```bash" + `
   go run main.go
   ` + "```" + `

3. The application will start on http://localhost:8080

### Available Endpoints

- GET /health - Health check endpoint
- GET /api/v1/status - Status endpoint

### Project Structure

` + "```" + `
{{.AppName}}/
├── cmd/               # CLI entry points
├── config/            # Configuration files
├── controllers/       # MVC controllers
├── internal/          # Internal packages
│   ├── middleware/    # Custom middleware
│   └── services/      # Business logic services
├── models/            # ORM models
├── migrations/        # Database migration files
├── public/            # Static assets
├── routes/            # Route definitions
├── templates/         # View templates
├── tests/             # Unit and integration tests
├── go.mod             # Go modules
└── main.go            # Application entry point
` + "```" + `

### CLI Commands

- ` + "`threadbolt new <app-name>`" + ` - Create a new ThreadBolt application
- ` + "`threadbolt generate model <name>`" + ` - Generate a new model
- ` + "`threadbolt generate controller <name>`" + ` - Generate a new controller
- ` + "`threadbolt migrate`" + ` - Run database migrations
- ` + "`threadbolt run`" + ` - Start the development server
- ` + "`threadbolt test`" + ` - Run tests

### Configuration

Configuration is handled through ` + "`config/config.yaml`" + ` and environment variables.
Environment variables should be prefixed with ` + "`THREADBOLT_`" + `.

Example:
- ` + "`THREADBOLT_SERVER_PORT=3000`" + `
- ` + "`THREADBOLT_DATABASE_DRIVER=postgres`" + `

## Development

### Adding a New Model

` + "```bash" + `
threadbolt generate model User
` + "```" + `

### Adding a New Controller

` + "```bash" + `
threadbolt generate controller User
` + "```" + `

### Running Tests

` + "```bash" + `
threadbolt test
` + "```" + `

### Database Migrations

` + "```bash" + `
threadbolt migrate
` + "```" + `

## License

This project is licensed under the MIT License.
`
