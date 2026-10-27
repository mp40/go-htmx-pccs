package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/render"
)

type stubDeleteCharacterEquipmentCharacterService struct {
	stubCharacterService
	character *domain.EnrichedCharacterDTO
	err       error
}

func (s *stubDeleteCharacterEquipmentCharacterService) DeleteCharacterEquipment(userID uuid.UUID, characterID uuid.UUID, equipmentID int, load string) (*domain.EnrichedCharacterDTO, error) {
	return s.character, s.err
}

type stubDeleteCharacterEquipmentIdentity struct {
	userID     *uuid.UUID
	isSignedIn bool
}

func (i *stubDeleteCharacterEquipmentIdentity) GetUserID(r *http.Request) *uuid.UUID {
	return i.userID
}

func (i *stubDeleteCharacterEquipmentIdentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestDeleteCharacterEquipmentHandler(t *testing.T) {
	t.Run("it returns 200 on delete primary load success", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubDeleteCharacterEquipmentIdentity{userID: &userID}

		character := domain.EnrichedCharacterDTO{}
		character.Name = "TEST-CHARACTER"
		character.Equipment.Secondary.Items = []domain.Item{}
		characterService := &stubDeleteCharacterEquipmentCharacterService{character: &character}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusOK {
			t.Errorf("got %v, want %v", got.StatusCode, http.StatusOK)
		}
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}
	})

	t.Run("it returns 400 on non-uuid character id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubDeleteCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubDeleteCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/character/xxxx/equipment/1?load=primary", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got %v, want %v", got.StatusCode, http.StatusBadRequest)
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

	t.Run("it returns 400 on non-integer equipment id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubDeleteCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubDeleteCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/xxxx?load=primary", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got %v, want %v", got.StatusCode, http.StatusBadRequest)
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

	t.Run("it returns 400 on missing load query parameter", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubDeleteCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubDeleteCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got %v, want %v", got.StatusCode, http.StatusBadRequest)
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

	t.Run("it returns 400 on malformed load query parameter", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubDeleteCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubDeleteCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=wrong", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got %v, want %v", got.StatusCode, http.StatusBadRequest)
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

	t.Run("it returns 401 on nil user id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		identity := stubDeleteCharacterEquipmentIdentity{}
		characterService := &stubDeleteCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusUnauthorized {
			t.Errorf("got %v, want %v", got.StatusCode, http.StatusUnauthorized)
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

	t.Run("it returns 404 on nil character", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubDeleteCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubDeleteCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusNotFound {
			t.Errorf("got %v, want %v", got.StatusCode, http.StatusNotFound)
		}

		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), "not found") {
			t.Errorf("want not found error message render, got: %s", body)
		}
	})

	t.Run("it returns 500 on delete error", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubDeleteCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubDeleteCharacterEquipmentCharacterService{err: fmt.Errorf("sadness")}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusInternalServerError {
			t.Errorf("got %v, want %v", got.StatusCode, http.StatusInternalServerError)
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
