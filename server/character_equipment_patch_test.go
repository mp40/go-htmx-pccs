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

type stubPatchCharacterEquipmentCharacterService struct {
	stubCharacterService
	character *domain.EnrichedCharacterDTO
	err       error
}

func (s *stubPatchCharacterEquipmentCharacterService) UpdateCharacterEquipment(userID uuid.UUID, characterID uuid.UUID, equipmentID int, load string, qty int) (*domain.EnrichedCharacterDTO, error) {
	return s.character, s.err
}

type stubPatchCharacterEquipmentIdentity struct {
	userID     *uuid.UUID
	isSignedIn bool
}

func (i *stubPatchCharacterEquipmentIdentity) GetUserID(r *http.Request) *uuid.UUID {
	return i.userID
}

func (i *stubPatchCharacterEquipmentIdentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestPatchCharacterEquipmentHandler(t *testing.T) {
	t.Run("it returns 200 on updating primary load success", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubPatchCharacterEquipmentIdentity{userID: &userID}

		character := domain.EnrichedCharacterDTO{}
		character.GunCombatActions = []int{1, 1, 1, 1}
		character.HandToHandCombatActions = []int{1, 0, 1, 0}
		character.Name = "TEST-CHARACTER"
		character.Equipment.Primary.Items = []domain.Item{{ID: 1, Name: "Equipment A", Weight: 2.2, Qty: 2}}
		characterService := &stubPatchCharacterEquipmentCharacterService{character: &character}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"qty": {"2"},
		}

		request := httptest.NewRequest(http.MethodPatch, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", strings.NewReader(formValues.Encode()))
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

		if !strings.Contains(string(body), "<span>2</span>") {
			t.Errorf("equipment name in document, got: %s", body)
		}
	})

	t.Run("it returns 400 on non-uuid character id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubPatchCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPatchCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"qty": {"2"},
		}

		request := httptest.NewRequest(http.MethodPatch, "/character/xxxx/equipment/1?load=primary", strings.NewReader(formValues.Encode()))
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
		identity := stubPatchCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPatchCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"qty": {"2"},
		}

		request := httptest.NewRequest(http.MethodPatch, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/xxxx?load=primary", strings.NewReader(formValues.Encode()))
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
		identity := stubPatchCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPatchCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"qty": {"2"},
		}

		request := httptest.NewRequest(http.MethodPatch, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1", strings.NewReader(formValues.Encode()))
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
		identity := stubPatchCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPatchCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"qty": {"2"},
		}

		request := httptest.NewRequest(http.MethodPatch, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=wrong", strings.NewReader(formValues.Encode()))
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

		identity := stubPatchCharacterEquipmentIdentity{}
		characterService := &stubPatchCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{
			"qty": {"2"},
		}

		request := httptest.NewRequest(http.MethodPatch, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", strings.NewReader(formValues.Encode()))
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

	t.Run("it returns 400 on malformed form value", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubPatchCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPatchCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{}

		request := httptest.NewRequest(http.MethodPatch, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", strings.NewReader(formValues.Encode()))
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

	t.Run("it returns 400 on non-integer qty form value", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubPatchCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPatchCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{"qty": {"xxxx"}}

		request := httptest.NewRequest(http.MethodPatch, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", strings.NewReader(formValues.Encode()))
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

	t.Run("it returns 400 on negative qty form value", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubPatchCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPatchCharacterEquipmentCharacterService{}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{"qty": {"-1"}}

		request := httptest.NewRequest(http.MethodPatch, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", strings.NewReader(formValues.Encode()))
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

	t.Run("it returns 500 on update error", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubPatchCharacterEquipmentIdentity{userID: &userID}
		characterService := &stubPatchCharacterEquipmentCharacterService{err: fmt.Errorf("sadness")}

		server := NewServer(nil, nil, &identity, characterService, nil, r)

		formValues := url.Values{"qty": {"2"}}

		request := httptest.NewRequest(http.MethodPatch, "/character/69000000-0000-4523-90a2-4868a5bfc90a/equipment/1?load=primary", strings.NewReader(formValues.Encode()))
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
