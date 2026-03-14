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

func (r *Render) RenderHomePage(w io.Writer, signedIn bool) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("top-strap.html"), getTemplatePath("home.html"))
	if err != nil {
		return err
	}

	data := struct {
		SignedIn bool
	}{
		SignedIn: signedIn,
	}

	err = tmpl.Execute(w, data)
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

func (r *Render) RenderAccountPage(w io.Writer, signedIn bool) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("top-strap.html"), getTemplatePath("account.html"))

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

func (r *Render) RenderSignUpModal(w io.Writer) error {
	tmpl, err := template.New("modal").ParseFiles(getTemplatePath("sign-up-modal.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) RenderErrorMessageFragment(w io.Writer, msg string) error {
	tmpl, err := template.New("error").ParseFiles(getTemplatePath("error-message.html"))
	if err != nil {
		return err
	}

	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: msg,
	}

	err = tmpl.Execute(w, data)
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
