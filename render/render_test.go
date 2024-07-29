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

func TestRenderCharacters(t *testing.T) {
	t.Run("it renders characters page", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderCharactersPage(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		wantPage := "Home"
		wantHomeFragment := `<h1>Characters</h1>`

		if !strings.Contains(got, wantPage) {
			t.Errorf("got '%s' want '%s'", got, wantPage)
		}

		if !strings.Contains(got, wantHomeFragment) {
			t.Errorf("got '%s' want '%s'", got, wantHomeFragment)
		}
	})

	t.Run("it renders characters fragment", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderCharactersFragment(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := `<h1>Characters</h1>`
		if !strings.Contains(got, want) {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

func TestRenderTools(t *testing.T) {
	t.Run("it renders tools page", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderToolsPage(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		wantPage := "Home"
		wantHomeFragment := `<h1>Tools</h1>`

		if !strings.Contains(got, wantPage) {
			t.Errorf("got '%s' want '%s'", got, wantPage)
		}

		if !strings.Contains(got, wantHomeFragment) {
			t.Errorf("got '%s' want '%s'", got, wantHomeFragment)
		}
	})

	t.Run("it renders tools fragment", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderToolsFragment(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := `<h1>Tools</h1>`
		if !strings.Contains(got, want) {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

func TestRenderToolsShooting(t *testing.T) {
	t.Run("it renders tools shooting page", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderToolsShootingPage(&buf)

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

	t.Run("it renders tools shooting fragment", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderToolsShootingFragment(&buf)

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

func TestRenderToolsShotguns(t *testing.T) {
	t.Run("it renders tools shotguns page", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderToolsShotgunsPage(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		wantPage := "Home"
		wantHomeFragment := `<h1>Shotguns</h1>`

		if !strings.Contains(got, wantPage) {
			t.Errorf("got '%s' want '%s'", got, wantPage)
		}

		if !strings.Contains(got, wantHomeFragment) {
			t.Errorf("got '%s' want '%s'", got, wantHomeFragment)
		}
	})

	t.Run("it renders tools shotguns fragment", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderToolsShotgunsFragment(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := `<h1>Shotguns</h1>`
		if !strings.Contains(got, want) {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}

func TestRenderToolsHandToHand(t *testing.T) {
	t.Run("it renders tools hand to hand page", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderToolsHandToHandPage(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		wantPage := "Home"
		wantHomeFragment := `<h1>Hand To Hand</h1>`

		if !strings.Contains(got, wantPage) {
			t.Errorf("got '%s' want '%s'", got, wantPage)
		}

		if !strings.Contains(got, wantHomeFragment) {
			t.Errorf("got '%s' want '%s'", got, wantHomeFragment)
		}
	})

	t.Run("it renders tools hand to hand fragment", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderToolsHandToHandFragment(&buf)

		if err != nil {
			t.Fatal(err)
		}

		got := buf.String()
		want := `<h1>Hand To Hand</h1>`
		if !strings.Contains(got, want) {
			t.Errorf("got '%s' want '%s'", got, want)
		}
	})
}
