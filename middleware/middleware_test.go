package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/session"
)

type StubStore struct {
	count        int
	err          error
	spyCountUser int
}

type StubSession struct {
	session       *session.Session // down the road DTO for this and other structs will help decouple
	err           error
	spyGetSession int
}

func (s *StubStore) CountUserByID(ID uuid.UUID) (int, error) {
	s.spyCountUser++
	return s.count, s.err
}

func (s *StubSession) GetSessionByID(ID uuid.UUID) (*session.Session, error) {
	s.spyGetSession++
	return s.session, s.err
}

func TestMiddleware(t *testing.T) {
	t.Run("it should handle cookie with uuid as user id", func(t *testing.T) {
		stubStore := StubStore{}
		stubSession := StubSession{}
		stubStore.count = 1
		middleware := NewMiddlewareService(&stubStore, &stubSession)

		ID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value(userContextKey).(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "userID", Value: ID.String()})

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubStore.spyCountUser != 1 {
			t.Errorf("want 1 call to data, got: %v", stubStore.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId == nil {
			t.Errorf("expected user uuid in context, got nil")
		}

		if *contextUserId != ID {
			t.Errorf("unexpected user uuid in context, got %v", contextUserId)
		}
	})

	t.Run("it should handle cookie with unparsable user id", func(t *testing.T) {
		stubStore := StubStore{}
		stubSession := StubSession{}
		middleware := NewMiddlewareService(&stubStore, &stubSession)

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value(userContextKey).(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "userID", Value: "whoops"})

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubStore.spyCountUser != 0 {
			t.Errorf("want 1 call to data, got: %v", stubStore.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle requests without cookie", func(t *testing.T) {
		stubStore := StubStore{}
		stubSession := StubSession{}
		middleware := NewMiddlewareService(&stubStore, &stubSession)

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value(userContextKey).(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubStore.spyCountUser != 0 {
			t.Errorf("want 0 call to data, got: %v", stubStore.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle user id with count 0", func(t *testing.T) {
		stubStore := StubStore{}
		stubSession := StubSession{}
		middleware := NewMiddlewareService(&stubStore, &stubSession)

		ID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value(userContextKey).(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "userID", Value: ID.String()})

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubStore.spyCountUser != 1 {
			t.Errorf("want 1 call to data, got: %v", stubStore.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle user id with count greater than 1", func(t *testing.T) {
		stubStore := StubStore{}
		stubSession := StubSession{}
		middleware := NewMiddlewareService(&stubStore, &stubSession)

		ID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value("userID").(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "userID", Value: ID.String()})

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubStore.spyCountUser != 1 {
			t.Errorf("want 1 call to data, got: %v", stubStore.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle errors from count user", func(t *testing.T) {
		stubStore := StubStore{}
		stubSession := StubSession{}
		stubStore.err = fmt.Errorf("fake-err")
		middleware := NewMiddlewareService(&stubStore, &stubSession)

		ID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")

		var contextUserId *uuid.UUID
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw, ok := r.Context().Value(userContextKey).(uuid.UUID)
			if ok {
				contextUserId = &raw
			}
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		request := httptest.NewRequest("GET", "/", nil)
		request.AddCookie(&http.Cookie{Name: "userID", Value: ID.String()})

		response := httptest.NewRecorder()
		handler.ServeHTTP(response, request)

		if stubStore.spyCountUser != 1 {
			t.Errorf("want 1 call to data, got: %v", stubStore.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})
}
