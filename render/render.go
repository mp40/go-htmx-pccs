package render

import (
	"bytes"
	"embed"
	"html/template"
	"io"

	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/state"
)

//go:embed templates/*.html
var templatesFS embed.FS

type Render struct {
	tmpl *template.Template
}

func NewRenderService() (*Render, error) {
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Render{tmpl: tmpl}, nil
}

func (r *Render) execute(w io.Writer, name string, data any) error {
	var buf bytes.Buffer
	err := r.tmpl.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	return err
}

func (r *Render) RenderHomePage(w io.Writer, signedIn bool) error {
	var page bytes.Buffer
	err := r.RenderHomeFragment(&page)
	if err != nil {
		return err
	}

	data := struct {
		SignedIn bool
		Page     template.HTML
	}{
		SignedIn: signedIn,
		Page:     template.HTML(page.String()),
	}

	return r.execute(w, "index", data)
}

func (r *Render) RenderHomeFragment(w io.Writer) error {
	return r.execute(w, "home", nil)
}

func (r *Render) RenderAccountPage(w io.Writer, signedIn bool, characters []domain.CharacterDTO) error {
	var fragment bytes.Buffer
	err := r.RenderAccountFragment(&fragment, characters)
	if err != nil {
		return err
	}

	var page bytes.Buffer
	err = r.renderAccountShell(&page, template.HTML(fragment.String()))
	if err != nil {
		return err
	}

	data := struct {
		SignedIn bool
		Page     template.HTML
	}{
		SignedIn: signedIn,
		Page:     template.HTML(page.String()),
	}

	return r.execute(w, "index", data)
}

func (r *Render) renderAccountShell(w io.Writer, fragment template.HTML) error {
	data := struct {
		Fragment template.HTML
	}{
		Fragment: fragment,
	}

	return r.execute(w, "account", data)
}

func (r *Render) RenderAccountFragment(w io.Writer, characters []domain.CharacterDTO) error {
	var inner bytes.Buffer
	data := struct {
		Characters []domain.CharacterDTO
	}{
		Characters: characters,
	}
	if err := r.execute(&inner, "characters", data); err != nil {
		return err
	}
	return r.renderAccountShell(w, template.HTML(inner.String()))
}

func (r *Render) RenderCharacterEditPage(w io.Writer, signedIn bool, character domain.CharacterDTO) error {
	var fragment bytes.Buffer
	err := r.RenderCharacterEditFragment(&fragment, character)
	if err != nil {
		return err
	}

	var page bytes.Buffer
	err = r.renderAccountShell(&page, template.HTML(fragment.String()))
	if err != nil {
		return err
	}

	data := struct {
		SignedIn bool
		Page     template.HTML
	}{
		SignedIn: signedIn,
		Page:     template.HTML(page.String()),
	}

	return r.execute(w, "index", data)
}

func (r *Render) RenderCharacterEditFragment(w io.Writer, character domain.CharacterDTO) error {
	data := struct {
		Character domain.CharacterDTO
	}{
		Character: character,
	}

	return r.execute(w, "character-edit", data)
}

func (r *Render) RenderReferencePage(w io.Writer, signedIn bool) error {
	var page bytes.Buffer
	err := r.RenderReferenceFragment(&page)
	if err != nil {
		return err
	}

	data := struct {
		SignedIn bool
		Page     template.HTML
	}{
		SignedIn: signedIn,
		Page:     template.HTML(page.String()),
	}

	return r.execute(w, "index", data)
}

func (r *Render) RenderReferenceFragment(w io.Writer) error {
	return r.execute(w, "reference", nil)
}

func (r *Render) RenderToolsPage(w io.Writer, signedIn bool) error {
	var page bytes.Buffer
	err := r.RenderToolsFragment(&page)
	if err != nil {
		return err
	}

	data := struct {
		SignedIn bool
		Page     template.HTML
	}{
		SignedIn: signedIn,
		Page:     template.HTML(page.String()),
	}

	return r.execute(w, "index", data)
}

func (r *Render) RenderToolsFragment(w io.Writer) error {
	return r.execute(w, "tools", nil)
}

func (r *Render) RenderShotgunPelletHitsFragment(w io.Writer, hits []int) error {
	data := struct {
		Hits []int
	}{
		Hits: hits,
	}

	return r.execute(w, "pellet-hits", data)
}

func (r *Render) RenderGearPage(w io.Writer, signedIn bool, equipment []state.Equipment) error {
	var page bytes.Buffer
	err := r.RenderGearFragment(&page, equipment)
	if err != nil {
		return err
	}

	data := struct {
		SignedIn bool
		Page     template.HTML
	}{
		SignedIn: signedIn,
		Page:     template.HTML(page.String()),
	}

	return r.execute(w, "index", data)
}

func (r *Render) RenderGearFragment(w io.Writer, equipment []state.Equipment) error {
	data := struct {
		Equipment []state.Equipment
	}{
		Equipment: equipment,
	}
	return r.execute(w, "gear", data)
}

func (r *Render) RenderSignInModal(w io.Writer) error {
	return r.execute(w, "sign-in-modal", nil)
}

func (r *Render) RenderSignUpModal(w io.Writer) error {
	return r.execute(w, "sign-up-modal", nil)
}

func (r *Render) RenderCharacter(w io.Writer, c domain.CharacterDTO) error {
	return r.execute(w, "character", c)
}

func (r *Render) RenderErrorMessageFragment(w io.Writer, msg string) error {
	data := struct {
		ErrorMessage string
	}{
		ErrorMessage: msg,
	}
	return r.execute(w, "error-message", data)
}

func (r *Render) RenderErrorListFragment(w io.Writer, errors []string) error {
	data := struct {
		Errors []string
	}{
		Errors: errors,
	}
	return r.execute(w, "error-list", data)
}
