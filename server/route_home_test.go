package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mp40/go-htmx-pccs/identity"
	"github.com/mp40/go-htmx-pccs/render"
)

type stubHomePageRender struct {
	spyRenderHomePage     int
	spyRenderHomeFragment int
}

type stubGetHomePageIdentity struct {
	isSignedIn bool
}

func (r *stubHomePageRender) RenderHomePage(w io.Writer, signedIn bool) error {
	r.spyRenderHomePage++
	return nil
}

func (r *stubHomePageRender) RenderHomeFragment(w io.Writer) error {
	r.spyRenderHomeFragment++
	return nil
}

func (i *stubGetHomePageIdentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestGetHomeHandler_Route(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}

func TestHomeHandler_Handler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := stubHomePageRender{}
		identity := stubGetHomePageIdentity{}

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		handler := getHomeHandler(&render, &identity)
		handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.spyRenderHomePage != 1 {
			t.Errorf("want 1 call to renderHomePage, got %d", render.spyRenderHomePage)
		}
		if render.spyRenderHomeFragment != 0 {
			t.Errorf("want 0 calls to renderHomeFragment, got %d", render.spyRenderHomeFragment)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := stubHomePageRender{}
		identity := stubGetHomePageIdentity{}

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()

		handler := getHomeHandler(&render, &identity)
		handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.spyRenderHomePage != 0 {
			t.Errorf("want 0 calls to renderHomePage, got %d", render.spyRenderHomePage)
		}
		if render.spyRenderHomeFragment != 1 {
			t.Errorf("want 1 call to renderHomeFragment, got %d", render.spyRenderHomeFragment)
		}
	})
}
