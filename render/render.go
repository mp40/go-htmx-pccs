package render

import (
	"html/template"
	"io"
	"path/filepath"
	"runtime"
)

func RenderIndexPage(w io.Writer) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.gohtml"), getTemplatePath("home.gohtml"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func getTemplatePath(templateName string) string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	return filepath.Join(basePath, "templates", templateName)
}
