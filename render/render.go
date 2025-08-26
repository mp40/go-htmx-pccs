package render

import (
	"html/template"
	"io"
	"path/filepath"
	"runtime"
)

type Render struct{}

func NewRenderService() *Render {
	return &Render{}
}

func (r *Render) RenderHomePage(w io.Writer) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("home.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) RenderHomeFragment(w io.Writer) error {
	tmpl, err := template.New("page").ParseFiles(getTemplatePath("home.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) RenderAccountPage(w io.Writer) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("account.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) RenderAccountFragment(w io.Writer) error {
	tmpl, err := template.New("page").ParseFiles(getTemplatePath("account.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) RenderSignInModal(w io.Writer) error {
	tmpl, err := template.New("modal").ParseFiles(getTemplatePath("sign-in-modal.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) RenderAccountSignInFailureFragment(w io.Writer) error {
	return nil
}

func getTemplatePath(templateName string) string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	return filepath.Join(basePath, "templates", templateName)
}
