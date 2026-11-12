package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/render"
	"github.com/mp40/go-htmx-pccs/state"
)

type stubGearFirearmGearService struct {
	firearm *domain.EnrichedFirearmDTO
	err     error
}

func (g *stubGearFirearmGearService) GetEnrichedFirearmByID(id int) (*domain.EnrichedFirearmDTO, error) {
	return g.firearm, g.err
}

func (g *stubGearFirearmGearService) GetEquipment() ([]state.Equipment, error) {
	panic("method not used in test")
}

func (g *stubGearFirearmGearService) GetFirearms() ([]domain.FirearmDTO, error) {
	panic("method not used in test")
}

func TestGetGearFirearmHandler(t *testing.T) {
	t.Run("it returns 200 on get enriched firearm by id success", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		f := domain.EnrichedFirearmDTO{
			FirearmDTO: domain.FirearmDTO{Name: "TEST FIREARM"},
		}

		gear := stubGearFirearmGearService{firearm: &f}

		server := NewServer(nil, nil, nil, nil, &gear, r)

		request := httptest.NewRequest(http.MethodGet, "/gear/firearm/1", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()
		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusOK {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusOK)
		}
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		if !strings.Contains(string(body), "TEST FIREARM") {
			t.Errorf("want firearm rendered, got: %s", body)
		}
	})

	t.Run("it returns 400 on malformed firearm id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		gear := stubGearFirearmGearService{}

		server := NewServer(nil, nil, nil, nil, &gear, r)

		request := httptest.NewRequest(http.MethodGet, "/gear/firearm/xxxx", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusBadRequest)
		}

		if !strings.Contains(string(body), "bad request") {
			t.Errorf("want bad request error message render, got: %s", body)
		}
	})

	t.Run("it returns 500 on get enriched firearm error", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		gear := stubGearFirearmGearService{err: fmt.Errorf("sadness")}

		server := NewServer(nil, nil, nil, nil, &gear, r)

		request := httptest.NewRequest(http.MethodGet, "/gear/firearm/1", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusInternalServerError {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusInternalServerError)
		}

		if !strings.Contains(string(body), "internal server error") {
			t.Errorf("want internal error message render, got: %s", body)
		}
	})

	t.Run("it returns 404 on nil enriched firearm", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		gear := stubGearFirearmGearService{}

		server := NewServer(nil, nil, nil, nil, &gear, r)

		request := httptest.NewRequest(http.MethodGet, "/gear/firearm/1", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusNotFound {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusNotFound)
		}

		if !strings.Contains(string(body), "not found") {
			t.Errorf("want not found error message render, got: %s", body)
		}
	})
}
