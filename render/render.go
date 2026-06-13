package render

import (
	"html/template"
	"io"
	"path/filepath"
	"runtime"

	"github.com/mp40/go-htmx-pccs/domain"
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

func (r *Render) RenderAccountPage(w io.Writer, signedIn bool, characters []domain.CharacterDTO) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("top-strap.html"), getTemplatePath("account.html"), getTemplatePath("characters.html"), getTemplatePath("character.html"))
	if err != nil {
		return err
	}

	data := struct {
		SignedIn   bool
		Characters []domain.CharacterDTO
	}{
		SignedIn:   signedIn,
		Characters: characters,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) RenderAccountFragment(w io.Writer, characters []domain.CharacterDTO) error {
	tmpl, err := template.New("page").ParseFiles(getTemplatePath("account.html"), getTemplatePath("characters.html"), getTemplatePath("character.html"))
	if err != nil {
		return err
	}

	data := struct {
		Characters []domain.CharacterDTO
	}{
		Characters: characters,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) RenderReferencePage(w io.Writer, signedIn bool) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("top-strap.html"), getTemplatePath("reference.html"))
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

func (r *Render) RenderReferenceFragment(w io.Writer) error {
	tmpl, err := template.New("page").ParseFiles(getTemplatePath("reference.html"))
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func (r *Render) RenderToolsPage(w io.Writer, signedIn bool) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("top-strap.html"), getTemplatePath("tools.html"))
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

func (r *Render) RenderToolsFragment(w io.Writer) error {
	tmpl, err := template.New("page").ParseFiles(getTemplatePath("tools.html"))
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

func (r *Render) RenderCharacter(w io.Writer, c domain.CharacterDTO) error {
	tmpl, err := template.New("character").ParseFiles(getTemplatePath("character.html"))
	if err != nil {
		return err
	}

	err = tmpl.Execute(w, c)
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
