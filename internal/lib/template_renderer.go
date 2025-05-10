package lib

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom renderer for Echo that uses Go's html/template package
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new template renderer with templates from the given directory
func NewTemplateRenderer(templatesDir string) *TemplateRenderer {
	// Parse templates from the layouts and pages directories
	templates := template.Must(template.ParseGlob(filepath.Join(templatesDir, "layouts", "*.html")))
	templates = template.Must(templates.ParseGlob(filepath.Join(templatesDir, "pages", "*.html")))

	// Add custom functions
	templates = templates.Funcs(template.FuncMap{
		// Add any custom functions here
	})

	return &TemplateRenderer{
		templates: templates,
	}
}

// Render renders a template with the given name and data
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}