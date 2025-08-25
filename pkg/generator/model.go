package generator

import (
	"fmt"
	"strings"
)

func GenerateModel(modelName string) error {
	fileName := fmt.Sprintf("models/%s.go", strings.ToLower(modelName))

	template := `package models

import (
	"gorm.io/gorm"
)

type {{.ModelName}} struct {
	BaseModel
	Name string ` + "`gorm:\"not null\" json:\"name\"`" + `
	// Add your fields here
}

// {{.ModelName}}Repository provides data access methods for {{.ModelName}}
type {{.ModelName}}Repository struct {
	db *gorm.DB
}

// New{{.ModelName}}Repository creates a new repository instance
func New{{.ModelName}}Repository(db *gorm.DB) *{{.ModelName}}Repository {
	return &{{.ModelName}}Repository{db: db}
}

// Create creates a new {{.ModelName}}
func (r *{{.ModelName}}Repository) Create({{.ModelNameLower}} *{{.ModelName}}) error {
	return r.db.Create({{.ModelNameLower}}).Error
}

// GetByID retrieves a {{.ModelName}} by ID
func (r *{{.ModelName}}Repository) GetByID(id uint) (*{{.ModelName}}, error) {
	var {{.ModelNameLower}} {{.ModelName}}
	err := r.db.First(&{{.ModelNameLower}}, id).Error
	if err != nil {
		return nil, err
	}
	return &{{.ModelNameLower}}, nil
}

// GetAll retrieves all {{.ModelName}}s
func (r *{{.ModelName}}Repository) GetAll() ([]{{.ModelName}}, error) {
	var {{.ModelNameLower}}s []{{.ModelName}}
	err := r.db.Find(&{{.ModelNameLower}}s).Error
	return {{.ModelNameLower}}s, err
}

// Update updates a {{.ModelName}}
func (r *{{.ModelName}}Repository) Update({{.ModelNameLower}} *{{.ModelName}}) error {
	return r.db.Save({{.ModelNameLower}}).Error
}

// Delete deletes a {{.ModelName}}
func (r *{{.ModelName}}Repository) Delete(id uint) error {
	return r.db.Delete(&{{.ModelName}}{}, id).Error
}
`

	data := struct {
		ModelName      string
		ModelNameLower string
	}{
		ModelName:      modelName,
		ModelNameLower: strings.ToLower(modelName),
	}

	return generateFile(fileName, template, data)
}
