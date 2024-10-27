package kubegen

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func newTemplate(name string) *template.Template {
	tmpl := template.New(name).Funcs(sprig.FuncMap()).Funcs(Funcs)
	tmpl.Funcs(template.FuncMap{
		"include": func(name string, value any) (string, error) {
			buf := &bytes.Buffer{}
			err := tmpl.ExecuteTemplate(buf, name, value)
			return buf.String(), err
		},
	})
	return tmpl
}

// Executes a template located at path with the specified data
func execTemplateFile(path string, data any) ([]byte, error) {
	tmpl, err := newTemplate(filepath.Base(path)).ParseFiles(path)
	if err != nil {
		return nil, err
	}
	return execTemplate(tmpl, data)
}

// Executes a template string with the specified data
func execTemplateString(text string, data any) ([]byte, error) {
	tmpl, err := newTemplate("stdin").Parse(text)
	if err != nil {
		return nil, err
	}
	return execTemplate(tmpl, data)
}

// Helper for execTemplateFile and execTemplateString - actually executes the template
func execTemplate(tmpl *template.Template, data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
