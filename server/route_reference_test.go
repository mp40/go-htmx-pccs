package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mp40/go-htmx-pccs/identity"
	"github.com/mp40/go-htmx-pccs/render"
)

type stubReferencePageRender struct {
	spyRenderReferencePage     int
	spyRenderReferenceFragment int
}

type stubGetReferencePageIdentity struct {
	isSignedIn bool
}

func (r *stubReferencePageRender) RenderReferencePage(w io.Writer, signedIn bool) error {
	r.spyRenderReferencePage++
	return nil
}

func (r *stubReferencePageRender) RenderReferenceFragment(w io.Writer) error {
	r.spyRenderReferenceFragment++
	return nil
}

func (i *stubGetReferencePageIdentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestGetRefrenceHandler_Route(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/reference", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusOK {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusOK)
		}

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "<h1>Reference</h1>") {
			t.Errorf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "<body>") {
			t.Errorf("expected fragment, got: %s", body)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/reference", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusOK {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusOK)
		}

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "<h1>Reference</h1>") {
			t.Errorf("unexpected body: %s", body)
		}
		if strings.Contains(string(body), "<body>") {
			t.Errorf("expected fragment, got: %s", body)
		}
	})
}

func TestReferenceHandler_Handler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := stubReferencePageRender{}
		identity := stubGetReferencePageIdentity{}

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		handler := getReferenceHandler(&render, &identity)
		handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.spyRenderReferencePage != 1 {
			t.Errorf("want 1 call to renderReferencePage, got %d", render.spyRenderReferencePage)
		}
		if render.spyRenderReferenceFragment != 0 {
			t.Errorf("want 0 calls to renderReferenceFragment, got %d", render.spyRenderReferenceFragment)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := stubReferencePageRender{}
		identity := stubGetReferencePageIdentity{}

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()

		handler := getReferenceHandler(&render, &identity)
		handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.spyRenderReferencePage != 0 {
			t.Errorf("want 0 calls to renderReferencePage, got %d", render.spyRenderReferencePage)
		}
		if render.spyRenderReferenceFragment != 1 {
			t.Errorf("want 1 call to renderReferenceFragment, got %d", render.spyRenderReferenceFragment)
		}
	})
}
