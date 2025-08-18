package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubRender struct {
	renderHomePageCalls        int
	renderHomeFragmentCalls    int
	renderAccountPageCalls     int
	renderAccountFragmentCalls int
}

func (r *StubRender) RenderHomePage(w io.Writer) error {
	r.renderHomePageCalls++
	return nil
}

func (r *StubRender) RenderHomeFragment(w io.Writer) error {
	r.renderHomeFragmentCalls++
	return nil
}

func (r *StubRender) RenderAccountPage(w io.Writer) error {
	r.renderAccountPageCalls++
	return nil
}

func (r *StubRender) RenderAccountFragment(w io.Writer) error {
	r.renderAccountFragmentCalls++
	return nil
}

func TestHomeHandler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		stub := StubRender{}
		server := NewServer(&stub)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stub.renderHomePageCalls != 1 {
			t.Errorf("want 1 call to renderHomePage, got %d", stub.renderHomePageCalls)
		}
		if stub.renderHomeFragmentCalls != 0 {
			t.Errorf("want 0 calls to renderHomeFragment, got %d", stub.renderHomeFragmentCalls)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		stub := StubRender{}
		server := NewServer(&stub)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stub.renderHomePageCalls != 0 {
			t.Errorf("want 0 calls to renderHomePage, got %d", stub.renderHomePageCalls)
		}
		if stub.renderHomeFragmentCalls != 1 {
			t.Errorf("want 1 call to renderHomeFragment, got %d", stub.renderHomeFragmentCalls)
		}
	})
}

func TestAccountHandler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		stub := StubRender{}
		server := NewServer(&stub)

		request, _ := http.NewRequest(http.MethodGet, "/account", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stub.renderAccountPageCalls != 1 {
			t.Errorf("want 0 calls to renderHomePage, got %d", stub.renderHomePageCalls)
		}
		if stub.renderAccountFragmentCalls != 0 {
			t.Errorf("want 1 call to renderHomeFragment, got %d", stub.renderHomeFragmentCalls)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		stub := StubRender{}
		server := NewServer(&stub)

		request, _ := http.NewRequest(http.MethodGet, "/account", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stub.renderAccountPageCalls != 0 {
			t.Errorf("want 0 calls to renderHomePage, got %d", stub.renderHomePageCalls)
		}
		if stub.renderAccountFragmentCalls != 1 {
			t.Errorf("want 1 call to renderHomeFragment, got %d", stub.renderHomeFragmentCalls)
		}
	})
}

func TestSignInHandler(t *testing.T) {
	t.Run("it should redirect to Home page on successful POST request", func(t *testing.T) {
		stub := StubRender{}
		server := NewServer(&stub)

		request, _ := http.NewRequest(http.MethodPost, "/account/sign-in", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		gotStatus := response.Result().StatusCode
		gotHeaderLocation := response.Result().Header.Get("Location")

		wantStatus := http.StatusSeeOther
		wantHeaderLocation := "/"

		if gotStatus != wantStatus {
			t.Errorf("got http status %v want http status %v", gotStatus, wantStatus)
		}

		if gotHeaderLocation != wantHeaderLocation {
			t.Errorf("got header Location %v, want header Location %v", gotHeaderLocation, wantHeaderLocation)
		}
	})
}
