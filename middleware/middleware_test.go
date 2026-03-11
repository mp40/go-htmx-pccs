package middleware

import (
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

func TestSignUp(t *testing.T) {
	t.Run("it should not panic", func(t *testing.T) {
		stubData := StubData{}
		stubData.count = 1
		middleware := NewMiddlewareService(&stubData)

		ID := uuid.MustParse("5ea69240-823c-4523-90a2-4868a5bfc90a")

		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		handler := middleware.AuthMiddleware(nextHandler)
		s := httptest.NewServer(handler)
		defer s.Close()

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
	})
}
