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

type stubToolsPageRender struct {
	spyRenderToolsPage     int
	spyRenderToolsFragment int
}

type stubGetToolsPageIdentity struct {
	isSignedIn bool
}

func (r *stubToolsPageRender) RenderToolsPage(w io.Writer, signedIn bool) error {
	r.spyRenderToolsPage++
	return nil
}

func (r *stubToolsPageRender) RenderToolsFragment(w io.Writer) error {
	r.spyRenderToolsFragment++
	return nil
}

func (i *stubGetToolsPageIdentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestGetToolsHandler_Route(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/tools", nil)
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

		if !strings.Contains(string(body), "<h1>Tools</h1>") {
			t.Errorf("unexpected page, got: %s", body)
		}
		if !strings.Contains(string(body), "<body>") {
			t.Errorf("expected full page, got: %s", body)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/tools", nil)
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

		if !strings.Contains(string(body), "<h1>Tools</h1>") {
			t.Errorf("unexpected body: %s", body)
		}
		if strings.Contains(string(body), "<body>") {
			t.Errorf("expected fragment, got: %s", body)
		}
	})
}

func TestToolsHandler_Handler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := stubToolsPageRender{}
		identity := stubGetToolsPageIdentity{}

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		handler := getToolsHandler(&render, &identity)
		handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.spyRenderToolsPage != 1 {
			t.Errorf("want 1 call to renderToolsPage, got %d", render.spyRenderToolsPage)
		}
		if render.spyRenderToolsFragment != 0 {
			t.Errorf("want 0 calls to renderToolsFragment, got %d", render.spyRenderToolsFragment)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := stubToolsPageRender{}
		identity := stubGetToolsPageIdentity{}

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()

		handler := getToolsHandler(&render, &identity)
		handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.spyRenderToolsPage != 0 {
			t.Errorf("want 0 calls to renderToolsPage, got %d", render.spyRenderToolsPage)
		}
		if render.spyRenderToolsFragment != 1 {
			t.Errorf("want 1 call to renderToolsFragment, got %d", render.spyRenderToolsFragment)
		}
	})
}
