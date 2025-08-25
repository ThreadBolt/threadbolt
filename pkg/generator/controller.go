package generator

import (
	"fmt"
	"strings"
)

func GenerateController(controllerName string) error {
	fileName := fmt.Sprintf("controllers/%s_controller.go", strings.ToLower(controllerName))
	
	template := `package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	
	"{{.AppName}}/models"
)

type {{.ControllerName}}Controller struct {
	repo *models.{{.ControllerName}}Repository
}

func New{{.ControllerName}}Controller(db *gorm.DB) *{{.ControllerName}}Controller {
	return &{{.ControllerName}}Controller{
		repo: models.New{{.ControllerName}}Repository(db),
	}
}

// Create{{.ControllerName}} handles POST /{{.ControllerNameLower}}s
func (c *{{.ControllerName}}Controller) Create{{.ControllerName}}(w http.ResponseWriter, r *http.Request) {
	var {{.ControllerNameLower}} models.{{.ControllerName}}
	
	if err := json.NewDecoder(r.Body).Decode(&{{.ControllerNameLower}}); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := c.repo.Create(&{{.ControllerNameLower}}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode({{.ControllerNameLower}})
}

// Get{{.ControllerName}} handles GET /{{.ControllerNameLower}}s/{id}
func (c *{{.ControllerName}}Controller) Get{{.ControllerName}}(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	{{.ControllerNameLower}}, err := c.repo.GetByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "{{.ControllerName}} not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode({{.ControllerNameLower}})
}

// GetAll{{.ControllerName}}s handles GET /{{.ControllerNameLower}}s
func (c *{{.ControllerName}}Controller) GetAll{{.ControllerName}}s(w http.ResponseWriter, r *http.Request) {
	{{.ControllerNameLower}}s, err := c.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode({{.ControllerNameLower}}s)
}

// Update{{.ControllerName}} handles PUT /{{.ControllerNameLower}}s/{id}
func (c *{{.ControllerName}}Controller) Update{{.ControllerName}}(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var {{.ControllerNameLower}} models.{{.ControllerName}}
	if err := json.NewDecoder(r.Body).Decode(&{{.ControllerNameLower}}); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	{{.ControllerNameLower}}.ID = uint(id)
	if err := c.repo.Update(&{{.ControllerNameLower}}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode({{.ControllerNameLower}})
}

// Delete{{.ControllerName}} handles DELETE /{{.ControllerNameLower}}s/{id}
func (c *{{.ControllerName}}Controller) Delete{{.ControllerName}}(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := c.repo.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
`

	data := struct {
		ControllerName      string
		ControllerNameLower string
		AppName             string
	}{
		ControllerName:      controllerName,
		ControllerNameLower: strings.ToLower(controllerName),
		AppName:             "example-app", // This should be dynamic
	}

	return generateFile(fileName, template, data)
}