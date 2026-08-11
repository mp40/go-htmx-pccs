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

func TestGetToolsShotgunSpreadHandler(t *testing.T) {
	t.Run("it returns fragment with pellet hits", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/tools/shotguns/spread?initialLocation=50&salm=10&hitCount=5", nil)
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

		if !strings.Contains(string(body), `<div class="pellet-hits">`) {
			t.Errorf("expected pellet hit fragment, got: %s", body)
		}
	})

	t.Run("it returns error on missing initial location", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/tools/shotguns/spread?salm=10&hitCount=5", nil)
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

		if !strings.Contains(string(body), `<div class="error-list">`) {
			t.Errorf("expected invalid message, got: %s", body)
		}
	})

	t.Run("it returns error on missing salm", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/tools/shotguns/spread?initialLocation=50&hitCount=5", nil)
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

		if !strings.Contains(string(body), `<div class="error-list">`) {
			t.Errorf("expected invalid message, got: %s", body)
		}
	})

	t.Run("it returns error on missing additional pellet hit count", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/tools/shotguns/spread?initialLocation=50&salm=5", nil)
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

		if !strings.Contains(string(body), `<div class="error-list">`) {
			t.Errorf("expected invalid message, got: %s", body)
		}
	})

	t.Run("it returns error on malformed initial location", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/tools/shotguns/spread?initialLocation=5x0&salm=5&hitCount=3", nil)
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

		if !strings.Contains(string(body), `<div class="error-list">`) {
			t.Errorf("expected invalid message, got: %s", body)
		}
	})

	t.Run("it returns error on malformed SALM", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/tools/shotguns/spread?initialLocation=50&salm=5XX&hitCount=3", nil)
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

		if !strings.Contains(string(body), `<div class="error-list">`) {
			t.Errorf("expected invalid message, got: %s", body)
		}
	})

	t.Run("it returns error on malformed hit count", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/tools/shotguns/spread?initialLocation=50&salm=5&hitCount=X3X", nil)
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

		if !strings.Contains(string(body), `<div class="error-list">`) {
			t.Errorf("expected invalid message, got: %s", body)
		}
	})
}
