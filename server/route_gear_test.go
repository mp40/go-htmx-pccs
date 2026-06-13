package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mp40/go-htmx-pccs/identity"
	"github.com/mp40/go-htmx-pccs/render"
	"github.com/mp40/go-htmx-pccs/state"
)

type stubGearPageGearService struct {
	equipment []state.Equipment
	err       error
}

func (g *stubGearPageGearService) GetEquipment() ([]state.Equipment, error) {
	return g.equipment, g.err
}

func TestGetGearHandler_Route(t *testing.T) {
	t.Run("it should return 200 and full page on successful GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}

		gear := stubGearPageGearService{}
		equipment := []state.Equipment{{ID: 1, Name: "TEST EQUIPMENT", Weight: 0.5}}
		gear.equipment = equipment

		server := NewServer(nil, nil, identity, nil, &gear, render)

		request := httptest.NewRequest(http.MethodGet, "/gear", nil)
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

		if !strings.Contains(string(body), "<h1>Gear</h1>") {
			t.Errorf("unexpected body: %s", body)
		}
		if !strings.Contains(string(body), "<body>") {
			t.Errorf("expected page, got: %s", body)
		}

		if !strings.Contains(string(body), "TEST EQUIPMENT") {
			t.Errorf("want equipment render, got: %s", body)
		}
	})

	t.Run("it should return 200 and partial on successful HTMX GET request", func(t *testing.T) {
		render := &render.Render{}
		identity := &identity.Identity{}

		gear := stubGearPageGearService{}
		equipment := []state.Equipment{{ID: 1, Name: "TEST EQUIPMENT", Weight: 0.5}}
		gear.equipment = equipment

		server := NewServer(nil, nil, identity, nil, &gear, render)

		request := httptest.NewRequest(http.MethodGet, "/gear", nil)
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

		if !strings.Contains(string(body), "<h1>Gear</h1>") {
			t.Errorf("unexpected body: %s", body)
		}
		if strings.Contains(string(body), "<body>") {
			t.Errorf("expected fragment, got: %s", body)
		}

		if !strings.Contains(string(body), "TEST EQUIPMENT") {
			t.Errorf("want equipment render, got: %s", body)
		}
	})
}
