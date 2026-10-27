package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/render"
	"github.com/mp40/go-htmx-pccs/state"
)

type stubRouteCharacterEquipmentSelectGearService struct {
	equipment []state.Equipment
	err       error
}

func (g *stubRouteCharacterEquipmentSelectGearService) GetEquipment() ([]state.Equipment, error) {
	return g.equipment, g.err
}

type stubRouteCharacterEquipmentSelectIndentity struct {
	userID     *uuid.UUID
	isSignedIn bool
}

func (i *stubRouteCharacterEquipmentSelectIndentity) GetUserID(r *http.Request) *uuid.UUID {
	return i.userID
}

func (i *stubRouteCharacterEquipmentSelectIndentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestGetCharacterEquipmentSelectHandler(t *testing.T) {
	t.Run("it should render list of selectable primary equipment fragment on successful GET request", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEquipmentSelectIndentity{userID: &userID, isSignedIn: true}
		gear := &stubRouteCharacterEquipmentSelectGearService{equipment: []state.Equipment{{ID: 1, Name: "TEST EQUIPMENT", Weight: 0.5}}}

		server := NewServer(nil, nil, &identity, nil, gear, r)

		request := httptest.NewRequest(http.MethodGet, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/select?load=primary", nil)
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

		if strings.Contains(string(body), "<!doctype html>") {
			t.Errorf("expected fragment, got: %s", body)
		}

		if !strings.Contains(string(body), "TEST EQUIPMENT") {
			t.Errorf("want equipment render, got: %s", body)
		}

		if !strings.Contains(string(body), `hx-vals='{"equipmentId": 1, "load": "primary"}'`) {
			t.Errorf("want equipment id and load values in hx-vals, got: %s", body)
		}
	})

	t.Run("it should render list of selectable secondary equipment fragment on successful GET request", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEquipmentSelectIndentity{userID: &userID, isSignedIn: true}
		gear := &stubRouteCharacterEquipmentSelectGearService{equipment: []state.Equipment{{ID: 1, Name: "TEST EQUIPMENT", Weight: 0.5}}}

		server := NewServer(nil, nil, &identity, nil, gear, r)

		request := httptest.NewRequest(http.MethodGet, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/select?load=secondary", nil)
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

		if strings.Contains(string(body), "<!doctype html>") {
			t.Errorf("expected fragment, got: %s", body)
		}

		if !strings.Contains(string(body), "TEST EQUIPMENT") {
			t.Errorf("want equipment render, got: %s", body)
		}

		if !strings.Contains(string(body), `hx-vals='{"equipmentId": 1, "load": "secondary"}'`) {
			t.Errorf("want equipment id and load values in hx-vals, got: %s", body)
		}
	})

	t.Run("it returns 401 on nil user id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		identity := stubRouteCharacterEquipmentSelectIndentity{}
		gear := &stubRouteCharacterEquipmentSelectGearService{}

		server := NewServer(nil, nil, &identity, nil, gear, r)

		request := httptest.NewRequest(http.MethodGet, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/select?load=primary", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusUnauthorized {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusUnauthorized)
		}
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), "unauthorised") {
			t.Errorf("want unauthorised error message render, got: %s", body)
		}
	})

	t.Run("it returns 400 on malformed character id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEquipmentSelectIndentity{userID: &userID, isSignedIn: true}
		gear := &stubRouteCharacterEquipmentSelectGearService{}

		server := NewServer(nil, nil, &identity, nil, gear, r)

		request := httptest.NewRequest(http.MethodGet, "/character/1-2-3/equipment/select?load=primary", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusBadRequest)
		}
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), "bad request") {
			t.Errorf("want bad request error message render, got: %s", body)
		}
	})

	t.Run("it returns 400 on malformed load name", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEquipmentSelectIndentity{userID: &userID, isSignedIn: true}
		gear := &stubRouteCharacterEquipmentSelectGearService{}

		server := NewServer(nil, nil, &identity, nil, gear, r)

		request := httptest.NewRequest(http.MethodGet, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/select?load=xxxx", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusBadRequest)
		}
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), "bad request") {
			t.Errorf("want bad request error message render, got: %s", body)
		}
	})

	t.Run("it returns 500 on get equipment error", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEquipmentSelectIndentity{userID: &userID, isSignedIn: true}
		gear := &stubRouteCharacterEquipmentSelectGearService{err: fmt.Errorf("sadness")}

		server := NewServer(nil, nil, &identity, nil, gear, r)

		request := httptest.NewRequest(http.MethodGet, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/select?load=primary", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusInternalServerError {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusInternalServerError)
		}
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), "internal server error") {
			t.Errorf("want internal server error message render, got: %s", body)
		}
	})
}
