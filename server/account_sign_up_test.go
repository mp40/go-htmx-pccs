package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mp40/go-htmx-pccs/render"
)

func TestGetSignUpHandler(t *testing.T) {
	t.Run("it should render the sign up modal on GET /account/sign-up", func(t *testing.T) {
		render := &render.Render{}
		server := NewServer(nil, nil, nil, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/account/sign-up", nil)
		request.Header.Set("HX-Request", "true")

		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusOK {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusOK)
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
		render := &render.Render{}
		server := NewServer(nil, nil, nil, nil, render)

		request := httptest.NewRequest(http.MethodGet, "/account/sign-up", nil)

		response := httptest.NewRecorder()
		server.Handler.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusInternalServerError {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusInternalServerError)
		}
	})
}
