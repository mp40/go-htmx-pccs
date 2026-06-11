package server

import (
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

func TestGetSignInHandler(t *testing.T) {
	t.Run("it should render the sign in modal on GET /account/sign-in", func(t *testing.T) {
		render := &render.Render{}
		server := NewServer(nil, nil, nil, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/account/sign-in", nil)
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

		if !strings.Contains(string(body), `hx-post="/account/sign-in"`) {
			t.Errorf("unexpected modal, got: %s", body)
		}
		if strings.Contains(string(body), "<body>") {
			t.Errorf("expected fragment, got: %s", body)
		}
	})

	t.Run("it returns error if not htmx request", func(t *testing.T) {
		render := &render.Render{}
		server := NewServer(nil, nil, nil, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/account/sign-in", nil)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusInternalServerError {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusInternalServerError)
		}
	})
}

type stubPostSignInAuth struct {
	user      *store.User
	spySignIn int
}

type stubPostSignInSession struct {
	spyAddSession int
}

func (a *stubPostSignInAuth) SignIn(email string, password string) (*store.User, error) {
	a.spySignIn++
	return a.user, nil
}

func (a *stubPostSignInAuth) SignUp(email string, password string) (*uuid.UUID, error) {
	panic("method not used in test")
}

func (s *stubPostSignInSession) AddSession(userID uuid.UUID, expiresAt time.Time) (uuid.UUID, error) {
	s.spyAddSession++
	return uuid.New(), nil
}

func (s *stubPostSignInSession) DeleteSessionByID(ID uuid.UUID) error {
	panic("method not used in test")
}

func TestPostSignInHandler(t *testing.T) {
	t.Run("it should redirect to Home page on successful POST request", func(t *testing.T) {
		session := stubPostSignInSession{}

		auth := stubPostSignInAuth{}
		user := store.User{}
		auth.user = &user

		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(&auth, &session, identity, nil, render)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-in", strings.NewReader(formValues.Encode()))
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

		if got.Header.Get("Hx-Redirect") != "/" {
			t.Errorf("got header Location %v, want header Location %v", got.Header.Get("Hx-Redirect"), "/")
		}

		if len(got.Cookies()) != 1 {
			t.Errorf("expected one cookie to be set, got %v", len(got.Cookies()))
		}

		if auth.spySignIn != 1 {
			t.Errorf("got %v calls to SignIn want 1", auth.spySignIn)
		}

		if session.spyAddSession != 1 {
			t.Errorf("got %v calls to AddSession want 1", session.spyAddSession)
		}
	})
}
