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

		server.getHomeHandler(response, request)

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

		server.getHomeHandler(response, request)

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
	t.Run("it should return 200 and full page on successuful GET request", func(t *testing.T) {
		stub := StubRender{}
		server := NewServer(&stub)

		request, _ := http.NewRequest(http.MethodGet, "/sign-in", nil)
		response := httptest.NewRecorder()

		server.getAccountHandler(response, request)

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

	t.Run("it should return 200 and partial on successuful HTMX GET request", func(t *testing.T) {
		stub := StubRender{}
		server := NewServer(&stub)

		request, _ := http.NewRequest(http.MethodPost, "/sign-in", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()

		server.getAccountHandler(response, request)

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
