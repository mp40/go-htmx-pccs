package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/store"
)

type StubAuth struct {
	err       error
	spySignUp int
	spySignIn int
	ID        *uuid.UUID
	user      *store.User
}

type StubRender struct {
	renderHomePageCalls             int
	renderHomeFragmentCalls         int
	renderAccountPageCalls          int
	renderAccountFragmentCalls      int
	renderSignInModalCalls          int
	renderSignUpModalCalls          int
	renderErrorMessageFragmentCalls int
}

func (r *StubRender) RenderHomePage(w io.Writer, signedIn bool) error {
	r.renderHomePageCalls++
	return nil
}

func (r *StubRender) RenderHomeFragment(w io.Writer) error {
	r.renderHomeFragmentCalls++
	return nil
}

func (r *StubRender) RenderAccountPage(w io.Writer, signedIn bool) error {
	r.renderAccountPageCalls++
	return nil
}

func (r *StubRender) RenderAccountFragment(w io.Writer) error {
	r.renderAccountFragmentCalls++
	return nil
}

func (r *StubRender) RenderErrorMessageFragment(w io.Writer, msg string) error {
	r.renderErrorMessageFragmentCalls++
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

func (a *StubAuth) SignIn(email string, password string) (*store.User, error) {
	a.spySignIn++
	return a.user, a.err
}

func (a *StubAuth) SignUp(email string, password string) (*uuid.UUID, error) {
	a.spySignUp++
	return a.ID, a.err
}

func TestHomeHandler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		stubRender := StubRender{}
		server := NewServer(nil, &stubRender)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
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

		request := httptest.NewRequest(http.MethodGet, "/", nil)
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

		request := httptest.NewRequest(http.MethodGet, "/account", nil)
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

		request := httptest.NewRequest(http.MethodGet, "/account", nil)
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

		request := httptest.NewRequest(http.MethodGet, "/account/sign-in", nil)
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

		user := store.User{}
		stubAuth.user = &user

		server := NewServer(&stubAuth, &stubRender)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-in", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		gotStatus := response.Result().StatusCode
		gotHeaderRedirect := response.Result().Header.Get("Hx-Redirect")

		wantStatus := http.StatusSeeOther
		wantHeaderRedirect := "/"

		gotCookies := response.Result().Cookies()

		if len(gotCookies) != 1 {
			t.Errorf("expected one cookie to be set, got %v", len(gotCookies))
		}

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

	t.Run("it should return error message if User not found", func(t *testing.T) {
		stubAuth := StubAuth{}
		stubRender := StubRender{}

		server := NewServer(&stubAuth, &stubRender)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-in", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if stubAuth.spySignIn != 1 {
			t.Errorf("got %v calls to SignIn want 1", stubAuth.spySignIn)
		}
		if stubRender.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", stubRender.renderErrorMessageFragmentCalls)
		}
	})

	t.Run("it should return error message on auth error", func(t *testing.T) {
		stubAuth := StubAuth{
			err: fmt.Errorf("boo"),
		}
		stubRender := StubRender{}
		server := NewServer(&stubAuth, &stubRender)

		stubAuth.err = fmt.Errorf("fake error")

		request := httptest.NewRequest(http.MethodPost, "/account/sign-in", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if stubAuth.spySignIn != 1 {
			t.Errorf("got %v calls to SignIn want 1", stubAuth.spySignIn)
		}
		if stubRender.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", stubRender.renderErrorMessageFragmentCalls)
		}
	})
}

func TestRenderSignUpModalHandler(t *testing.T) {
	t.Run("it should render the sign in modal on GET /account/sign-up", func(t *testing.T) {
		stubRender := StubRender{}
		server := NewServer(nil, &stubRender)

		request := httptest.NewRequest(http.MethodGet, "/account/sign-up", nil)
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

		newID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")
		stubAuth.ID = &newID

		server := NewServer(&stubAuth, &stubRender)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		gotStatus := response.Result().StatusCode
		gotHeaderRedirect := response.Result().Header.Get("Hx-Redirect")

		wantStatus := http.StatusSeeOther
		wantHeaderRedirect := "/"

		gotCookies := response.Result().Cookies()

		if len(gotCookies) != 1 {
			t.Errorf("expected one cookie to be set, got %v", len(gotCookies))
		}

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

	t.Run("it should return error message if invalid email format", func(t *testing.T) {
		stubAuth := StubAuth{}
		stubRender := StubRender{}
		server := NewServer(&stubAuth, &stubRender)

		formValues := url.Values{
			"email":    {"invalid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if stubAuth.spySignUp != 0 {
			t.Errorf("got %v calls to SignUp want 0", stubAuth.spySignUp)
		}

		if stubRender.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", stubRender.renderErrorMessageFragmentCalls)
		}
	})

	t.Run("it should return error message if invalid email format", func(t *testing.T) {
		stubAuth := StubAuth{}
		stubRender := StubRender{}
		server := NewServer(&stubAuth, &stubRender)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"bad"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if stubAuth.spySignUp != 0 {
			t.Errorf("got %v calls to SignUp want 0", stubAuth.spySignUp)
		}

		if stubRender.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", stubRender.renderErrorMessageFragmentCalls)
		}
	})

	t.Run("it should return 400 and user feedback unsuccessful POST request", func(t *testing.T) {
		stubAuth := StubAuth{
			err: fmt.Errorf("fake error"),
		}
		stubRender := StubRender{}
		server := NewServer(&stubAuth, &stubRender)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if stubRender.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", stubRender.renderErrorMessageFragmentCalls)
		}
	})
}
