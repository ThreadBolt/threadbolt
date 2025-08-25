package framework

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/ThreadBolt/threadbolt/pkg/config"
	"github.com/ThreadBolt/threadbolt/pkg/di"
	"github.com/ThreadBolt/threadbolt/pkg/orm"
)

type App struct {
	Router    *mux.Router
	DB        *gorm.DB
	Config    *viper.Viper
	Container *di.Container
}

func LoadApp() (*App, error) {
	// Validate project structure
	if err := validateProjectStructure(); err != nil {
		return nil, fmt.Errorf("invalid project structure: %w", err)
	}

	app := &App{
		Router:    mux.NewRouter(),
		Container: di.NewContainer(),
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	app.Config = cfg

	// Initialize database
	db, err := orm.Initialize(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	app.DB = db

	// Register database in DI container
	app.Container.Register("db", db)

	// Load routes
	if err := app.loadRoutes(); err != nil {
		return nil, fmt.Errorf("failed to load routes: %w", err)
	}

	return app, nil
}

func (a *App) Start(port string) error {
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, a.Router)
}

func (a *App) RunMigrations() error {
	return orm.RunMigrations(a.DB)
}

func validateProjectStructure() error {
	requiredDirs := []string{
		"controllers",
		"models",
		"config",
		"routes",
		"internal/middleware",
		"internal/services",
		"migrations",
		"templates",
		"public",
		"tests",
	}

	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("required directory '%s' is missing", dir)
		}
	}

	requiredFiles := []string{
		"main.go",
		"go.mod",
		"config/config.yaml",
		"routes/routes.go",
	}

	for _, file := range requiredFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return fmt.Errorf("required file '%s' is missing", file)
		}
	}

	return nil
}

func (a *App) loadRoutes() error {
	// This would load routes from routes/routes.go
	// For now, we'll add a simple health check
	a.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy","framework":"threadbolt"}`)
	}).Methods("GET")

	return nil
}
