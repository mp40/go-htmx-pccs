package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

type StubData struct {
	count        int
	err          error
	spyCountUser int
}

func (d *StubData) CountUserByID(ID uuid.UUID) (int, error) {
	d.spyCountUser++
	return d.count, d.err
}

func TestMiddleware(t *testing.T) {
	t.Run("it should handle cookie with uuid as user id", func(t *testing.T) {
		stubData := StubData{}
		stubData.count = 1
		middleware := NewMiddlewareService(&stubData)

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

		if stubData.spyCountUser != 1 {
			t.Errorf("want 1 call to data, got: %v", stubData.spyCountUser)
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
		stubData := StubData{}
		middleware := NewMiddlewareService(&stubData)

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

		if stubData.spyCountUser != 0 {
			t.Errorf("want 1 call to data, got: %v", stubData.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle requests without cookie", func(t *testing.T) {
		stubData := StubData{}
		middleware := NewMiddlewareService(&stubData)

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

		if stubData.spyCountUser != 0 {
			t.Errorf("want 0 call to data, got: %v", stubData.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle user id with count 0", func(t *testing.T) {
		stubData := StubData{}
		stubData.count = 0
		middleware := NewMiddlewareService(&stubData)

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

		if stubData.spyCountUser != 1 {
			t.Errorf("want 1 call to data, got: %v", stubData.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle user id with count greater than 1", func(t *testing.T) {
		stubData := StubData{}
		stubData.count = 2
		middleware := NewMiddlewareService(&stubData)

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

		if stubData.spyCountUser != 1 {
			t.Errorf("want 1 call to data, got: %v", stubData.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})

	t.Run("it should handle errors from count user", func(t *testing.T) {
		stubData := StubData{}
		stubData.err = fmt.Errorf("fake-err")
		middleware := NewMiddlewareService(&stubData)

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

		if stubData.spyCountUser != 1 {
			t.Errorf("want 1 call to data, got: %v", stubData.spyCountUser)
		}

		if response.Code != http.StatusOK {
			t.Errorf("unexpected status, got: %v", response.Code)
		}

		if contextUserId != nil {
			t.Errorf("unexpected user uuid in context, got %v", *contextUserId)
		}
	})
}
