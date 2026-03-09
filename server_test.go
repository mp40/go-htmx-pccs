package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type StubAuth struct {
	err       error
	spySignUp int
	spySignIn int
}

type StubRender struct {
	renderHomePageCalls                     int
	renderHomeFragmentCalls                 int
	renderAccountPageCalls                  int
	renderAccountFragmentCalls              int
	renderSignInModalCalls                  int
	renderSignUpModalCalls                  int
	renderAccountSignInFailureFragmentCalls int
	renderAccountSignUpFailureFragmentCalls int
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

func (r *StubRender) RenderAccountSignInFailureFragment(w io.Writer) error {
	r.renderAccountSignInFailureFragmentCalls++
	return nil
}

func (r *StubRender) RenderAccountSignUpFailureFragment(w io.Writer) error {
	r.renderAccountSignUpFailureFragmentCalls++
	return nil
}

func (r *StubRender) RenderSignInModal(w io.Writer) error {
	r.renderSignInModalCalls++
	return nil
}

func (r *StubRender) RenderSignUpModal(w io.Writer) error {
	r.renderSignUpModalCalls++
	return nil
}

func (a *StubAuth) SignIn(email string, password string) error {
	a.spySignIn++
	return a.err
}

func (a *StubAuth) SignUp(email string, password string) error {
	a.spySignUp++
	return a.err
}

func TestHomeHandler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		stubRender := StubRender{}
		server := NewServer(nil, &stubRender)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stubRender.renderHomePageCalls != 1 {
			t.Errorf("want 1 call to renderHomePage, got %d", stubRender.renderHomePageCalls)
		}
		if stubRender.renderHomeFragmentCalls != 0 {
			t.Errorf("want 0 calls to renderHomeFragment, got %d", stubRender.renderHomeFragmentCalls)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		stubRender := StubRender{}
		server := NewServer(nil, &stubRender)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stubRender.renderHomePageCalls != 0 {
			t.Errorf("want 0 calls to renderHomePage, got %d", stubRender.renderHomePageCalls)
		}
		if stubRender.renderHomeFragmentCalls != 1 {
			t.Errorf("want 1 call to renderHomeFragment, got %d", stubRender.renderHomeFragmentCalls)
		}
	})
}

func TestAccountHandler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		stubRender := StubRender{}
		server := NewServer(nil, &stubRender)

		request, _ := http.NewRequest(http.MethodGet, "/account", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stubRender.renderAccountPageCalls != 1 {
			t.Errorf("want 0 calls to renderAccountPageCalls, got %d", stubRender.renderAccountPageCalls)
		}
		if stubRender.renderAccountFragmentCalls != 0 {
			t.Errorf("want 1 call to renderAccountFragmentCalls, got %d", stubRender.renderAccountFragmentCalls)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		stubRender := StubRender{}
		server := NewServer(nil, &stubRender)

		request, _ := http.NewRequest(http.MethodGet, "/account", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stubRender.renderAccountPageCalls != 0 {
			t.Errorf("want 0 calls to renderAccountPageCalls, got %d", stubRender.renderAccountPageCalls)
		}
		if stubRender.renderAccountFragmentCalls != 1 {
			t.Errorf("want 1 call to renderAccountFragmentCalls, got %d", stubRender.renderAccountFragmentCalls)
		}
	})
}

func TestRenderSignInModalHandler(t *testing.T) {
	t.Run("it should render the sign in modal on GET /account/sign-in", func(t *testing.T) {
		stubRender := StubRender{}
		server := NewServer(nil, &stubRender)

		request, _ := http.NewRequest(http.MethodGet, "/account/sign-in", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stubRender.renderSignInModalCalls != 1 {
			t.Errorf("want 1 call to renderSignInModalCalls, got %d", stubRender.renderSignInModalCalls)
		}
	})
}

func TestSignInHandler(t *testing.T) {
	t.Run("it should redirect to Home page on successful POST request", func(t *testing.T) {
		stubAuth := StubAuth{}
		stubRender := StubRender{}
		server := NewServer(&stubAuth, &stubRender)

		request, _ := http.NewRequest(http.MethodPost, "/account/sign-in", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		gotStatus := response.Result().StatusCode
		gotHeaderRedirect := response.Result().Header.Get("Hx-Redirect")

		wantStatus := http.StatusSeeOther
		wantHeaderRedirect := "/"

		if stubAuth.spySignIn != 1 {
			t.Errorf("got %v calls to SignIn want 1", stubAuth.spySignIn)
		}

		if gotStatus != wantStatus {
			t.Errorf("got http status %v want http status %v", gotStatus, wantStatus)
		}

		if gotHeaderRedirect != wantHeaderRedirect {
			t.Errorf("got header Location %v, want header Location %v", gotHeaderRedirect, wantHeaderRedirect)
		}
	})

	t.Run("it should return 400 and user feedback unsuccessful POST request", func(t *testing.T) {
		stubAuth := StubAuth{
			err: fmt.Errorf("boo"),
		}
		stubRender := StubRender{}
		server := NewServer(&stubAuth, &stubRender)

		stubAuth.err = fmt.Errorf("fake error")

		request, _ := http.NewRequest(http.MethodPost, "/account/sign-in", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusBadRequest

		if got != want {
			t.Errorf("got http status %v want http status %v", got, want)
		}
		if stubRender.renderAccountSignInFailureFragmentCalls != 1 {
			t.Errorf("want 1 call to renderHomeFragment, got %d", stubRender.renderAccountSignInFailureFragmentCalls)
		}
	})
}

func TestRenderSignUpModalHandler(t *testing.T) {
	t.Run("it should render the sign in modal on GET /account/sign-up", func(t *testing.T) {
		stubRender := StubRender{}
		server := NewServer(nil, &stubRender)

		request, _ := http.NewRequest(http.MethodGet, "/account/sign-up", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if stubRender.renderSignUpModalCalls != 1 {
			t.Errorf("want 1 call to renderSignUpModalCalls, got %d", stubRender.renderSignUpModalCalls)
		}
	})
}

func TestSignUpHandler(t *testing.T) {
	t.Run("it should redirect to Home page on successful POST request", func(t *testing.T) {
		stubAuth := StubAuth{}
		stubRender := StubRender{}
		server := NewServer(&stubAuth, &stubRender)

		request, _ := http.NewRequest(http.MethodPost, "/account/sign-up", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		gotStatus := response.Result().StatusCode
		gotHeaderRedirect := response.Result().Header.Get("Hx-Redirect")

		wantStatus := http.StatusSeeOther
		wantHeaderRedirect := "/"

		if stubAuth.spySignUp != 1 {
			t.Errorf("got %v calls to SignUp want 1", stubAuth.spySignUp)
		}

		if gotStatus != wantStatus {
			t.Errorf("got http status %v want http status %v", gotStatus, wantStatus)
		}

		if gotHeaderRedirect != wantHeaderRedirect {
			t.Errorf("got header Location %v, want header Location %v", gotHeaderRedirect, wantHeaderRedirect)
		}
	})

	t.Run("it should return 400 and user feedback unsuccessful POST request", func(t *testing.T) {
		stubAuth := StubAuth{
			err: fmt.Errorf("boo"),
		}
		stubRender := StubRender{}
		server := NewServer(&stubAuth, &stubRender)

		stubAuth.err = fmt.Errorf("fake error")

		request, _ := http.NewRequest(http.MethodPost, "/account/sign-up", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusBadRequest

		if got != want {
			t.Errorf("got http status %v want http status %v", got, want)
		}
		if stubRender.renderAccountSignUpFailureFragmentCalls != 1 {
			t.Errorf("want 1 call to renderAccountSignUpFailureFragment, got %d", stubRender.renderAccountSignUpFailureFragmentCalls)
		}
	})
}
