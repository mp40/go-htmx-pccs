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

type stubAccountPageRender struct {
	spyRenderAccountPage     int
	spyRenderAccountFragment int
}

type stubGetAccountPageIdentity struct {
	isSignedIn bool
}

func (r *stubAccountPageRender) RenderHomePage(w io.Writer, signedIn bool) error {
	r.spyRenderAccountPage++
	return nil
}

func (r *stubAccountPageRender) RenderHomeFragment(w io.Writer) error {
	r.spyRenderAccountFragment++
	return nil
}

func (i *stubGetAccountPageIdentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestGetAccountHandler_Route(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/account", nil)
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

		if !strings.Contains(string(body), "<h1>Account</h1>") {
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

		request := httptest.NewRequest(http.MethodGet, "/account", nil)
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

		if !strings.Contains(string(body), "<h1>Account</h1>") {
			t.Errorf("unexpected body: %s", body)
		}
		if strings.Contains(string(body), "<body>") {
			t.Errorf("expected fragment, got: %s", body)
		}
	})
}

func TestAccountHandler_Handler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := stubAccountPageRender{}
		identity := stubGetAccountPageIdentity{}

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		handler := getHomeHandler(&render, &identity)
		handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.spyRenderAccountPage != 1 {
			t.Errorf("want 1 call to renderAccountPage, got %d", render.spyRenderAccountPage)
		}
		if render.spyRenderAccountFragment != 0 {
			t.Errorf("want 0 calls to renderAccountFragment, got %d", render.spyRenderAccountFragment)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := stubAccountPageRender{}
		identity := stubGetAccountPageIdentity{}

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

		if render.spyRenderAccountPage != 0 {
			t.Errorf("want 0 calls to renderAccountPage, got %d", render.spyRenderAccountPage)
		}
		if render.spyRenderAccountFragment != 1 {
			t.Errorf("want 1 call to renderAccountFragment, got %d", render.spyRenderAccountFragment)
		}
	})
}
