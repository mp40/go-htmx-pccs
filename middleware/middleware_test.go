package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/session"
)

type StubSession struct {
	session       *session.Session
	err           error
	spyGetSession int
}

func (s *StubSession) GetValidSessionByID(ID uuid.UUID) (*session.Session, error) {
	s.spyGetSession++
	return s.session, s.err
}

func fakeEnrichFunc(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, "fake-user-context-key", userID)
}

func TestMiddleware(t *testing.T) {
	t.Run("it should handle cookie with uuid as session id", func(t *testing.T) {
		stubSession := StubSession{}
		userID := uuid.MustParse("10000000-0000-4523-90a2-4868a5bfc90a")
		fakeSession := session.Session{UserID: userID}
		stubSession.session = &fakeSession

		middleware := NewMiddlewareService(&stubSession, fakeEnrichFunc)

		sessionID := "69000000-0000-4523-90a2-4868a5bfc90a"
		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value("fake-user-context-key").(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "sessionID", Value: sessionID})

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubSession.spyGetSession != 1 {
			t.Errorf("want 1 call to data, got: %v", stubSession.spyGetSession)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId == nil {
			t.Errorf("expected user uuid in context, got nil")
		}

		if *contextUserId != userID {
			t.Errorf("unexpected user uuid in context, got %v", contextUserId)
		}
	})

	t.Run("it should handle cookie with unparsable session id", func(t *testing.T) {
		stubSession := StubSession{}
		middleware := NewMiddlewareService(&stubSession, fakeEnrichFunc)

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value("fake-user-context-key").(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "sessionID", Value: "whoops"})

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubSession.spyGetSession != 0 {
			t.Errorf("want 0 calls to session data, got: %v", stubSession.spyGetSession)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle requests without cookie", func(t *testing.T) {
		stubSession := StubSession{}
		middleware := NewMiddlewareService(&stubSession, fakeEnrichFunc)

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value("fake-user-context-key").(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubSession.spyGetSession != 0 {
			t.Errorf("want 0 calls to session data, got: %v", stubSession.spyGetSession)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle null session", func(t *testing.T) {
		stubSession := StubSession{}
		middleware := NewMiddlewareService(&stubSession, fakeEnrichFunc)

		sessionID := "69000000-0000-4523-90a2-4868a5bfc90a"

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value("fake-user-context-key").(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "sessionID", Value: sessionID})

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubSession.spyGetSession != 1 {
			t.Errorf("want 1 call to session data, got: %v", stubSession.spyGetSession)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle errors from get session", func(t *testing.T) {
		stubSession := StubSession{}
		stubSession.err = fmt.Errorf("fake-err")
		middleware := NewMiddlewareService(&stubSession, fakeEnrichFunc)

		sessionID := "69000000-0000-4523-90a2-4868a5bfc90a"

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value("fake-user-context-key").(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "sessionID", Value: sessionID})

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubSession.spyGetSession != 1 {
			t.Errorf("want 1 call to session data, got: %v", stubSession.spyGetSession)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})
}
