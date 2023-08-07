package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type TemplateHandler struct {
	once     sync.Once
	File     string
	template *template.Template
}

func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.File)))
	})

	t.template.Execute(w, nil)
}
