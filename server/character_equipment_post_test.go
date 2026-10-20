package server

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/render"
)

type stubPostCharacterEquipmentCharacterService struct {
	stubCharacterService
	character *domain.EnrichedCharacterDTO
	err       error
}

func (s *stubPostCharacterEquipmentCharacterService) AddCharacterEquipment(userID uuid.UUID, characterID uuid.UUID, equipmentID int, load string) (*domain.EnrichedCharacterDTO, error) {
	return s.character, s.err
}

type stubPostCharacterEquipmentIdentity struct {
	userID     *uuid.UUID
	isSignedIn bool
}

func (i *stubPostCharacterEquipmentIdentity) GetUserID(r *http.Request) *uuid.UUID {
	return i.userID
}

func (i *stubPostCharacterEquipmentIdentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestPostCharacterEquipmentHandler(t *testing.T) {
	t.Run("it returns 200 on success", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubPostCharacterEquipmentIdentity{userID: &userID}

		character := domain.EnrichedCharacterDTO{}
		character.GunCombatActions = []int{1, 1, 1, 1}
		character.HandToHandCombatActions = []int{1, 0, 1, 0}
		character.Name = "TEST-CHARACTER"
		character.Equipment.Primary.Items = []domain.Item{{ID: 1, Name: "Equipment A", Weight: 2.2, Qty: 1}}
		characterService := &stubPostCharacterEquipmentCharacterService{character: &character}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"equipmentId": {"1"},
			"load":        {"primary"},
		}

		request := httptest.NewRequest(http.MethodPost, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment", strings.NewReader(formValues.Encode()))
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

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(string(body), "Equipment A") {
			t.Errorf("expected equipment name in document, got: %s", body)
		}
	})

	t.Run("it returns 400 on non-uuid character id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")

		identity := stubPostCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPostCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"equipmentId": {"1"},
			"load":        {"primary"},
		}

		request := httptest.NewRequest(http.MethodPost, "/character/x/equipment", strings.NewReader(formValues.Encode()))
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

	t.Run("it returns 401 if user id is nil", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		identity := stubPostCharacterEquipmentIdentity{}
		characterService := &stubPostCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"equipmentId": {"1"},
			"load":        {"primary"},
		}

		request := httptest.NewRequest(http.MethodPost, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment", strings.NewReader(formValues.Encode()))
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

	t.Run("it returns 400 on invalid load", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")

		identity := stubPostCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPostCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"equipmentId": {"1"},
			"load":        {"garbage"},
		}

		request := httptest.NewRequest(http.MethodPost, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment", strings.NewReader(formValues.Encode()))
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

		identity := stubPostCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPostCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"equipmentId": {"x"},
			"load":        {"primary"},
		}

		request := httptest.NewRequest(http.MethodPost, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment", strings.NewReader(formValues.Encode()))
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

	t.Run("it returns 500 on add equipment error", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")

		identity := stubPostCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPostCharacterEquipmentCharacterService{err: fmt.Errorf("sadness")}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"equipmentId": {"1"},
			"load":        {"primary"},
		}

		request := httptest.NewRequest(http.MethodPost, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment", strings.NewReader(formValues.Encode()))
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
