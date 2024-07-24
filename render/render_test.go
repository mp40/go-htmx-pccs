package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderHome(t *testing.T) {
	t.Run("it renders home page", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderHomePage(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		wantPage := "Home"
		wantHomeFragment := `<h1>In Development</h1>`

		if !strings.Contains(got, wantPage) {
			t.Errorf("got '%s' want '%s'", got, wantPage)
		}

		if !strings.Contains(got, wantHomeFragment) {
			t.Errorf("got '%s' want '%s'", got, wantHomeFragment)
		}
	})

	t.Run("it renders home fragment", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderHomeFragment(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := `<h1>In Development</h1>`
		if !strings.Contains(got, want) {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

func TestRenderShooting(t *testing.T) {
	t.Run("it renders shooting page", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderShootingPage(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		wantPage := "Home"
		wantHomeFragment := `<h1>Shooting</h1>`

		if !strings.Contains(got, wantPage) {
			t.Errorf("got '%s' want '%s'", got, wantPage)
		}

		if !strings.Contains(got, wantHomeFragment) {
			t.Errorf("got '%s' want '%s'", got, wantHomeFragment)
		}
	})

	t.Run("it renders shooting fragment", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderShootingFragment(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := `<h1>Shooting</h1>`
		if !strings.Contains(got, want) {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}
