package render

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderIndex(t *testing.T) {

	t.Run("it renders index page", func(t *testing.T) {
		buf := bytes.Buffer{}
		err := RenderIndexPage(&buf)

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
