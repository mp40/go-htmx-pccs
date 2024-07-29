package render

import (
	"html/template"
	"io"
	"path/filepath"
	"runtime"
)

func RenderHomePage(w io.Writer) error {
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

func RenderHomeFragment(w io.Writer) error {
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

func RenderCharactersPage(w io.Writer) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("characters.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func RenderCharactersFragment(w io.Writer) error {
	tmpl, err := template.New("page").ParseFiles(getTemplatePath("characters.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func RenderToolsPage(w io.Writer) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("tools.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func RenderToolsFragment(w io.Writer) error {
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

func RenderToolsShootingPage(w io.Writer) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("tools-shooting.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func RenderToolsShootingFragment(w io.Writer) error {
	tmpl, err := template.New("page").ParseFiles(getTemplatePath("tools-shooting.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func RenderToolsShotgunsPage(w io.Writer) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("tools-shotguns.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func RenderToolsShotgunsFragment(w io.Writer) error {
	tmpl, err := template.New("page").ParseFiles(getTemplatePath("tools-shotguns.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func RenderToolsHandToHandPage(w io.Writer) error {
	tmpl, err := template.ParseFiles(getTemplatePath("index.html"), getTemplatePath("tools-hand-to-hand.html"))

	if err != nil {
		return err
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		return err
	}

	return nil
}

func RenderToolsHandToHandFragment(w io.Writer) error {
	tmpl, err := template.New("page").ParseFiles(getTemplatePath("tools-hand-to-hand.html"))

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
