package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/identity"
	"github.com/mp40/go-htmx-pccs/render"
)

type stubRouteAccountIndentity struct {
	userID *uuid.UUID
}

func (i *stubRouteAccountIndentity) IsSignedIn(r *http.Request) bool {
	return true
}

func (i *stubRouteAccountIndentity) GetUserID(r *http.Request) *uuid.UUID {
	panic("method not used in test")
}

func TestGetAccountHandler_Route(t *testing.T) {
	t.Run("it should redirect to home if not signed in", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/account", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusSeeOther {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusSeeOther)
		}

		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		if got.Header.Get("Hx-Redirect") != "/" {
			t.Errorf("got header Location %v, want header Location %v", got.Header.Get("Hx-Redirect"), "/")
		}
	})

	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &stubRouteAccountIndentity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/account", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusOK {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusOK)
		}
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
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
		identity := &stubRouteAccountIndentity{}
		server := NewServer(nil, nil, identity, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/account", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusOK {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusOK)
		}
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
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
