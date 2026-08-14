package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/identity"
	"github.com/mp40/go-htmx-pccs/render"
	"github.com/mp40/go-htmx-pccs/store"
)

func TestGetSignUpHandler(t *testing.T) {
	t.Run("it should render the sign up modal on GET /account/sign-up", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, nil, nil, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/account/sign-up", nil)
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

		if !strings.Contains(string(body), `hx-post="/account/sign-up"`) {
			t.Errorf("unexpected modal, got: %s", body)
		}
		if strings.Contains(string(body), "<body>") {
			t.Errorf("expected fragment, got: %s", body)
		}
	})

	t.Run("it returns error if not htmx request", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, nil, nil, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/account/sign-up", nil)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusInternalServerError {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusInternalServerError)
		}
	})
}

type stubPostSignUpAuth struct {
	userID    *uuid.UUID
	spySignUp int
	err       error
}

type stubPostSignUpSession struct {
	spyAddSession int
}

func (a *stubPostSignUpAuth) SignUp(email string, password string) (*uuid.UUID, error) {
	a.spySignUp++
	return a.userID, nil
}

func (a *stubPostSignUpAuth) SignIn(email string, password string) (*store.User, error) {
	panic("method not used in test")
}

func (s *stubPostSignUpSession) AddSession(userID uuid.UUID, expiresAt time.Time) (uuid.UUID, error) {
	s.spyAddSession++
	return uuid.New(), nil
}

func (s *stubPostSignUpSession) DeleteSessionByID(ID uuid.UUID) error {
	panic("method not used in test")
}

func TestPostSignUpHandler(t *testing.T) {
	t.Run("it should redirect to Account page on successful POST request", func(t *testing.T) {
		auth := stubPostSignUpAuth{}
		session := stubPostSignUpSession{}

		newID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")
		auth.userID = &newID

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(&auth, &session, identity, nil, nil, r)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusSeeOther {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusSeeOther)
		}

		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		if got.Header.Get("Hx-Redirect") != "/account" {
			t.Errorf("got header Location %v, want header Location %v", got.Header.Get("Hx-Redirect"), "/account")
		}

		if len(got.Cookies()) != 1 {
			t.Errorf("expected cookie to be set, got %v", len(got.Cookies()))
		}

		if auth.spySignUp != 1 {
			t.Errorf("got %v calls to SignUp want 1", auth.spySignUp)
		}

		if session.spyAddSession != 1 {
			t.Errorf("got %v calls to AddSession want 1", session.spyAddSession)
		}
	})

	t.Run("it should return error message if invalid email format", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, nil, r)

		formValues := url.Values{
			"email":    {"invalid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "<span>invalid sign up:") {
			t.Errorf("expected invalid message, got: %s", body)
		}
	})

	t.Run("it should return error message if invalid password", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(nil, nil, identity, nil, nil, r)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"bad"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusOK {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusOK)
		}

		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		if !strings.Contains(string(body), "<span>invalid sign up:") {
			t.Errorf("expected invalid message, got: %s", body)
		}
	})

	t.Run("it should return 400 and user feedback unsuccessful POST request", func(t *testing.T) {
		auth := stubPostSignUpAuth{err: fmt.Errorf("fake error")}

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		identity := &identity.Identity{}
		server := NewServer(&auth, nil, identity, nil, nil, r)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusOK {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusOK)
		}

		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		if !strings.Contains(string(body), "<span>internal server error</span>") {
			t.Errorf("expected error fragment, got: %s", body)
		}
	})
}
