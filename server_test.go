package main

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
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/store"
)

type stubAuth struct {
	err       error
	spySignUp int
	spySignIn int
	ID        *uuid.UUID
	user      *store.User
}

type stubCharacterService struct {
	err             error
	spyAddCharacter int
	character       *store.Character
}

type stubSession struct {
	sessionID        uuid.UUID
	err              error
	spyAddSession    int
	spyDeleteSession int
}

type stubrender struct {
	renderHomePageCalls             int
	renderHomeFragmentCalls         int
	renderAccountPageCalls          int
	renderAccountFragmentCalls      int
	renderSignInModalCalls          int
	renderSignUpModalCalls          int
	renderCharacterCalls            int
	renderErrorMessageFragmentCalls int
}

type stubIdentity struct {
	userID     *uuid.UUID
	isSignedIn bool
}

func (r *stubrender) RenderHomePage(w io.Writer, signedIn bool) error {
	r.renderHomePageCalls++
	return nil
}

func (r *stubrender) RenderHomeFragment(w io.Writer) error {
	r.renderHomeFragmentCalls++
	return nil
}

func (r *stubrender) RenderAccountPage(w io.Writer, signedIn bool) error {
	r.renderAccountPageCalls++
	return nil
}

func (r *stubrender) RenderAccountFragment(w io.Writer, signedIn bool) error {
	r.renderAccountFragmentCalls++
	return nil
}

func (r *stubrender) RenderErrorMessageFragment(w io.Writer, msg string) error {
	r.renderErrorMessageFragmentCalls++
	return nil
}

func (r *stubrender) RenderSignInModal(w io.Writer) error {
	r.renderSignInModalCalls++
	return nil
}

func (r *stubrender) RenderSignUpModal(w io.Writer) error {
	r.renderSignUpModalCalls++
	return nil
}

func (r *stubrender) RenderCharcater(w io.Writer) error {
	r.renderCharacterCalls++
	return nil
}

func (a *stubAuth) SignIn(email string, password string) (*store.User, error) {
	a.spySignIn++
	return a.user, a.err
}

func (a *stubAuth) SignUp(email string, password string) (*uuid.UUID, error) {
	a.spySignUp++
	return a.ID, a.err
}

func (s *stubSession) AddSession(userID uuid.UUID, expiresAt time.Time) (uuid.UUID, error) {
	s.spyAddSession++
	return s.sessionID, s.err
}

func (s *stubSession) DeleteSessionByID(ID uuid.UUID) error {
	s.spyDeleteSession++
	return s.err
}

func (s *stubCharacterService) AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*store.Character, error) {
	s.spyAddCharacter++
	return s.character, s.err
}

func (i *stubIdentity) GetUserID(r *http.Request) *uuid.UUID {
	return i.userID
}

func (i *stubIdentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestHomeHandler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := stubrender{}
		server := NewServer(nil, nil, &stubIdentity{}, nil, &render)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.renderHomePageCalls != 1 {
			t.Errorf("want 1 call to renderHomePage, got %d", render.renderHomePageCalls)
		}
		if render.renderHomeFragmentCalls != 0 {
			t.Errorf("want 0 calls to renderHomeFragment, got %d", render.renderHomeFragmentCalls)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := stubrender{}
		server := NewServer(nil, nil, &stubIdentity{}, nil, &render)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.renderHomePageCalls != 0 {
			t.Errorf("want 0 calls to renderHomePage, got %d", render.renderHomePageCalls)
		}
		if render.renderHomeFragmentCalls != 1 {
			t.Errorf("want 1 call to renderHomeFragment, got %d", render.renderHomeFragmentCalls)
		}
	})
}

func TestAccountHandler(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := stubrender{}
		server := NewServer(nil, nil, &stubIdentity{}, nil, &render)

		request := httptest.NewRequest(http.MethodGet, "/account", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.renderAccountPageCalls != 1 {
			t.Errorf("want 0 calls to renderAccountPageCalls, got %d", render.renderAccountPageCalls)
		}
		if render.renderAccountFragmentCalls != 0 {
			t.Errorf("want 1 call to renderAccountFragmentCalls, got %d", render.renderAccountFragmentCalls)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := stubrender{}
		server := NewServer(nil, nil, &stubIdentity{}, nil, &render)

		request := httptest.NewRequest(http.MethodGet, "/account", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.renderAccountPageCalls != 0 {
			t.Errorf("want 0 calls to renderAccountPageCalls, got %d", render.renderAccountPageCalls)
		}
		if render.renderAccountFragmentCalls != 1 {
			t.Errorf("want 1 call to renderAccountFragmentCalls, got %d", render.renderAccountFragmentCalls)
		}
	})
}

func TestRenderSignInModalHandler(t *testing.T) {
	t.Run("it should render the sign in modal on GET /account/sign-in", func(t *testing.T) {
		render := stubrender{}
		server := NewServer(nil, nil, &stubIdentity{}, nil, &render)

		request := httptest.NewRequest(http.MethodGet, "/account/sign-in", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.renderSignInModalCalls != 1 {
			t.Errorf("want 1 call to renderSignInModalCalls, got %d", render.renderSignInModalCalls)
		}
	})
}

func TestSignInHandler(t *testing.T) {
	t.Run("it should redirect to Home page on successful POST request", func(t *testing.T) {
		auth := stubAuth{}
		session := stubSession{}
		render := stubrender{}

		user := store.User{}
		auth.user = &user

		server := NewServer(&auth, &session, &stubIdentity{}, nil, &render)

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

		if auth.spySignIn != 1 {
			t.Errorf("got %v calls to SignIn want 1", auth.spySignIn)
		}

		if session.spyAddSession != 1 {
			t.Errorf("got %v calls to AddSession want 1", session.spyAddSession)
		}

		if gotStatus != wantStatus {
			t.Errorf("got http status %v want http status %v", gotStatus, wantStatus)
		}

		if gotHeaderRedirect != wantHeaderRedirect {
			t.Errorf("got header Location %v, want header Location %v", gotHeaderRedirect, wantHeaderRedirect)
		}
	})

	t.Run("it should return error message if User not found", func(t *testing.T) {
		auth := stubAuth{}
		session := stubSession{}
		render := stubrender{}

		server := NewServer(&auth, &session, &stubIdentity{}, nil, &render)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-in", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if auth.spySignIn != 1 {
			t.Errorf("got %v calls to SignIn want 1", auth.spySignIn)
		}
		if render.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", render.renderErrorMessageFragmentCalls)
		}
	})

	t.Run("it should return error message on auth error", func(t *testing.T) {
		auth := stubAuth{
			err: fmt.Errorf("boo"),
		}
		session := stubSession{}
		render := stubrender{}

		server := NewServer(&auth, &session, &stubIdentity{}, nil, &render)

		auth.err = fmt.Errorf("fake error")

		request := httptest.NewRequest(http.MethodPost, "/account/sign-in", nil)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if auth.spySignIn != 1 {
			t.Errorf("got %v calls to SignIn want 1", auth.spySignIn)
		}
		if render.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", render.renderErrorMessageFragmentCalls)
		}
	})
}

func TestRenderSignUpModalHandler(t *testing.T) {
	t.Run("it should render the sign in modal on GET /account/sign-up", func(t *testing.T) {
		render := stubrender{}
		server := NewServer(nil, nil, &stubIdentity{}, nil, &render)

		request := httptest.NewRequest(http.MethodGet, "/account/sign-up", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result().StatusCode
		want := http.StatusOK

		if got != want {
			t.Errorf("got %v want %v", got, want)
		}

		if render.renderSignUpModalCalls != 1 {
			t.Errorf("want 1 call to renderSignUpModalCalls, got %d", render.renderSignUpModalCalls)
		}
	})
}

func TestSignUpHandler(t *testing.T) {
	t.Run("it should redirect to Home page on successful POST request", func(t *testing.T) {
		auth := stubAuth{}
		session := stubSession{}
		characterService := stubCharacterService{}
		render := stubrender{}

		newID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")
		auth.ID = &newID

		server := NewServer(&auth, &session, &stubIdentity{}, &characterService, &render)

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

		if auth.spySignUp != 1 {
			t.Errorf("got %v calls to SignUp want 1", auth.spySignUp)
		}

		if session.spyAddSession != 1 {
			t.Errorf("got %v calls to AddSession want 1", session.spyAddSession)
		}

		if gotStatus != wantStatus {
			t.Errorf("got http status %v want http status %v", gotStatus, wantStatus)
		}

		if gotHeaderRedirect != wantHeaderRedirect {
			t.Errorf("got header Location %v, want header Location %v", gotHeaderRedirect, wantHeaderRedirect)
		}
	})

	t.Run("it should return error message if invalid email format", func(t *testing.T) {
		auth := stubAuth{}
		session := stubSession{}
		characterService := stubCharacterService{}
		render := stubrender{}

		server := NewServer(&auth, &session, &stubIdentity{}, &characterService, &render)

		formValues := url.Values{
			"email":    {"invalid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if auth.spySignUp != 0 {
			t.Errorf("got %v calls to SignUp want 0", auth.spySignUp)
		}

		if render.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", render.renderErrorMessageFragmentCalls)
		}
	})

	t.Run("it should return error message if invalid email format", func(t *testing.T) {
		auth := stubAuth{}
		session := stubSession{}
		render := stubrender{}

		server := NewServer(&auth, &session, &stubIdentity{}, nil, &render)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"bad"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if auth.spySignUp != 0 {
			t.Errorf("got %v calls to SignUp want 0", auth.spySignUp)
		}

		if render.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", render.renderErrorMessageFragmentCalls)
		}
	})

	t.Run("it should return 400 and user feedback unsuccessful POST request", func(t *testing.T) {
		auth := stubAuth{
			err: fmt.Errorf("fake error"),
		}
		session := stubSession{}
		render := stubrender{}

		server := NewServer(&auth, &session, &stubIdentity{}, nil, &render)

		formValues := url.Values{
			"email":    {"762@valid.com"},
			"password": {"fake-password"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/sign-up", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		if render.renderErrorMessageFragmentCalls != 1 {
			t.Errorf("want 1 call to renderErrorMessageFragment, got %d", render.renderErrorMessageFragmentCalls)
		}
	})
}

func TestSignOutHandler(t *testing.T) {
	t.Run("it should redirect to Home page on successful sign out", func(t *testing.T) {
		auth := stubAuth{}
		session := stubSession{}
		render := stubrender{}

		server := NewServer(&auth, &session, &stubIdentity{}, nil, &render)

		cookie := getSecureCookie(uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a"))

		request := httptest.NewRequest(http.MethodDelete, "/account/sign-out", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.AddCookie(&cookie)
		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		gotStatus := response.Result().StatusCode
		gotHeaderRedirect := response.Result().Header.Get("Hx-Redirect")

		wantStatus := http.StatusSeeOther
		wantHeaderRedirect := "/"

		gotCookies := response.Result().Cookies()

		if len(gotCookies) != 1 {
			t.Errorf("expected cookie to be set, got %v", len(gotCookies))
		}

		if gotCookies[0].MaxAge >= 0 {
			t.Errorf("expected cookie max age to be less than 0, got %v", len(gotCookies))
		}

		if session.spyDeleteSession != 1 {
			t.Errorf("got %v calls to DeleteSessionByID want 1", session.spyDeleteSession)
		}

		if gotStatus != wantStatus {
			t.Errorf("got http status %v want http status %v", gotStatus, wantStatus)
		}

		if gotHeaderRedirect != wantHeaderRedirect {
			t.Errorf("got header Location %v, want header Location %v", gotHeaderRedirect, wantHeaderRedirect)
		}
	})
}

func TestPostCharacterHandler(t *testing.T) {
	t.Run("it should add character and return character", func(t *testing.T) {
		auth := stubAuth{}
		session := stubSession{}
		characterService := stubCharacterService{}
		render := stubrender{}

		userID := uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a")
		stubID := stubIdentity{userID: &userID, isSignedIn: true}

		server := NewServer(&auth, &session, &stubID, &characterService, &render)

		formValues := url.Values{
			"str": {"9"},
			"int": {"8"},
			"wil": {"7"},
			"hlt": {"6"},
			"agi": {"5"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/characters", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		gotStatus := response.Result().StatusCode
		wantStatus := http.StatusCreated

		if gotStatus != wantStatus {
			t.Errorf("got http status %v want http status %v", gotStatus, wantStatus)
		}

		if characterService.spyAddCharacter != 1 {
			t.Errorf("got %v calls to AddCharacter want 1", characterService.spyAddCharacter)
		}

		if render.renderCharacterCalls != 1 {
			t.Errorf("got %v calls to render character want 1", render.renderCharacterCalls)
		}
	})
}
