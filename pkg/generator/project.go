package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func CreateNewProject(appName string) error {
	// Create project directory
	if err := os.MkdirAll(appName, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Change to project directory
	if err := os.Chdir(appName); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Create directory structure
	dirs := []string{
		"cmd",
		"config",
		"controllers",
		"internal/middleware",
		"internal/services",
		"models",
		"migrations",
		"public",
		"routes",
		"templates",
		"tests",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate files
	files := map[string]string{
		"main.go":                    mainGoTemplate,
		"go.mod":                     goModTemplate,
		"config/config.yaml":         configTemplate,
		"routes/routes.go":           routesTemplate,
		"controllers/health_controller.go": healthControllerTemplate,
		"models/base.go":             baseModelTemplate,
		"internal/middleware/cors.go": corsMiddlewareTemplate,
		".gitignore":                 gitignoreTemplate,
		"README.md":                  readmeTemplate,
	}

	data := struct {
		AppName string
	}{
		AppName: appName,
	}

	for filePath, templateContent := range files {
		if err := generateFile(filePath, templateContent, data); err != nil {
			return fmt.Errorf("failed to generate %s: %w", filePath, err)
		}
	}

	return nil
}

func generateFile(filePath, templateContent string, data interface{}) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	tmpl, err := template.New(filePath).Parse(templateContent)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}