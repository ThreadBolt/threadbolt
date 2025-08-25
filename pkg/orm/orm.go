package orm

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(config *viper.Viper) (*gorm.DB, error) {
	driver := config.GetString("database.driver")

	var dialector gorm.Dialector

	switch driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			config.GetString("database.host"),
			config.GetString("database.username"),
			config.GetString("database.password"),
			config.GetString("database.name"),
			config.GetString("database.port"),
			config.GetString("database.sslmode"),
		)
		dialector = postgres.Open(dsn)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.GetString("database.username"),
			config.GetString("database.password"),
			config.GetString("database.host"),
			config.GetString("database.port"),
			config.GetString("database.name"),
		)
		dialector = mysql.Open(dsn)

	case "sqlite":
		dsn := config.GetString("database.name")
		if dsn == "" {
			dsn = "springgo.db"
		}
		dialector = sqlite.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}

	// Configure GORM
	gormConfig := &gorm.Config{}

	// Set log level based on environment
	if config.GetString("environment") == "production" {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	} else {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	// Auto-migrate all models in the models directory
	modelFiles, err := filepath.Glob("models/*.go")
	if err != nil {
		return fmt.Errorf("failed to find model files: %w", err)
	}

	if len(modelFiles) == 0 {
		fmt.Println("No model files found for migration")
		return nil
	}

	// For now, we'll just run auto-migrate
	// In a real implementation, you'd parse the model files and register them
	fmt.Printf("Found %d model files for migration\n", len(modelFiles))

	// Run custom migrations from migrations directory
	return runCustomMigrations(db)
}

func runCustomMigrations(db *gorm.DB) error {
	migrationFiles := []string{}

	err := filepath.WalkDir("migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			migrationFiles = append(migrationFiles, path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	for _, file := range migrationFiles {
		fmt.Printf("Running migration: %s\n", file)
		// Here you would execute the SQL file
		// This is a simplified version
	}

	return nil
}