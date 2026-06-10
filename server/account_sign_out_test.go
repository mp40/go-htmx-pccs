package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/identity"
	"github.com/mp40/go-htmx-pccs/render"
)

type stubDeleteSignOutSession struct {
	spyDeleteSession int
}

func (s *stubDeleteSignOutSession) DeleteSessionByID(ID uuid.UUID) error {
	s.spyDeleteSession++
	return nil
}

func (s *stubDeleteSignOutSession) AddSession(userID uuid.UUID, expiresAt time.Time) (uuid.UUID, error) {
	panic("method not used in test")
}

func TestDeleteSignOutHandler(t *testing.T) {
	t.Run("it should redirect to Home page on successful sign out", func(t *testing.T) {
		session := stubDeleteSignOutSession{}

		render := &render.Render{}
		identity := &identity.Identity{}
		server := NewServer(nil, &session, identity, nil, render)

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
