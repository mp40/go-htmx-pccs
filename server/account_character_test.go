package server

import (
	"context"
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

type stubPostCharacter struct {
	character       *domain.CharacterDTO
	err             error
	spyAddCharacter int
}

type stubPostCharacterIndentity struct {
	userID *uuid.UUID
}

func (c *stubPostCharacter) AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error) {
	c.spyAddCharacter++
	return c.character, c.err
}

func (c *stubPostCharacter) EditCharacter(userID uuid.UUID, characterID uuid.UUID, rawCharacter domain.RawCharacterEdit) (*domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (c *stubPostCharacter) GetCharactersByUserID(userID uuid.UUID) ([]domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (i *stubPostCharacterIndentity) GetUserID(r *http.Request) *uuid.UUID {
	return i.userID
}

func (i *stubPostCharacterIndentity) IsSignedIn(r *http.Request) bool {
	panic("method not used in test")
}

func TestPostCharacterHandler(t *testing.T) {
	t.Run("it should add character and return character", func(t *testing.T) {
		identity := stubPostCharacterIndentity{}
		userID := uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a")
		identity.userID = &userID

		characterService := stubPostCharacter{}
		character := domain.CharacterDTO{RawCharacter: domain.RawCharacter{Name: "TEST-CHARACTER"}}
		characterService.character = &character

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		formValues := url.Values{
			"str":                {"9"},
			"int":                {"8"},
			"wil":                {"7"},
			"hlt":                {"6"},
			"agi":                {"5"},
			"tch":                {"4"},
			"gun_combat_level":   {"1"},
			"hand_to_hand_level": {"0"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/characters", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := request.Context()
		ctx = context.WithValue(ctx, "userID", "69000000-0000-4523-90a2-4868a5bfc90a")
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusCreated {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusCreated)
		}
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		if characterService.spyAddCharacter != 1 {
			t.Errorf("got %v calls to AddCharacter want 1", characterService.spyAddCharacter)
		}

		if !strings.Contains(string(body), "<span>TEST-CHARACTER</span>") {
			t.Errorf("expected character fragment, got: %s", body)
		}
	})

	t.Run("it should return error message when error parsing raw character", func(t *testing.T) {
		identity := stubPostCharacterIndentity{}
		userID := uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a")
		identity.userID = &userID

		characterService := stubPostCharacter{}

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		formValues := url.Values{
			"str":                {"garbage"},
			"int":                {"8"},
			"wil":                {"7"},
			"hlt":                {"6"},
			"agi":                {"5"},
			"tch":                {"4"},
			"gun_combat_level":   {"1"},
			"hand_to_hand_level": {"0"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/characters", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := request.Context()
		ctx = context.WithValue(ctx, "userID", "69000000-0000-4523-90a2-4868a5bfc90a")
		request = request.WithContext(ctx)

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
		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		if characterService.spyAddCharacter != 0 {
			t.Errorf("got %v calls to AddCharacter want 0", characterService.spyAddCharacter)
		}

		if !strings.Contains(string(body), "<span>str is not a number</span>") {
			t.Errorf("expected character fragment, got: %s", body)
		}
	})

	t.Run("it should return error on add character error", func(t *testing.T) {
		identity := stubPostCharacterIndentity{}
		userID := uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a")
		identity.userID = &userID

		characterService := stubPostCharacter{}
		characterService.err = fmt.Errorf("FAKE ERROR")

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		formValues := url.Values{
			"str":                {"9"},
			"int":                {"8"},
			"wil":                {"7"},
			"hlt":                {"6"},
			"agi":                {"5"},
			"tch":                {"4"},
			"gun_combat_level":   {"1"},
			"hand_to_hand_level": {"0"},
		}

		request := httptest.NewRequest(http.MethodPost, "/account/characters", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := request.Context()
		ctx = context.WithValue(ctx, "userID", "69000000-0000-4523-90a2-4868a5bfc90a")
		request = request.WithContext(ctx)

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

		if characterService.spyAddCharacter != 1 {
			t.Errorf("got %v calls to AddCharacter want 1", characterService.spyAddCharacter)
		}

		if !strings.Contains(string(body), "<span>internal server error</span>") {
			t.Errorf("expected character fragment, got: %s", body)
		}
	})
}

type stubPutCharacter struct {
	character        *domain.CharacterDTO
	err              error
	spyEditCharacter int
	spyCharacterID   uuid.UUID
}

type stubPutCharacterIndentity struct {
	userID *uuid.UUID
}

func (c *stubPutCharacter) EditCharacter(userID uuid.UUID, characterID uuid.UUID, rawCharacter domain.RawCharacterEdit) (*domain.CharacterDTO, error) {
	c.spyCharacterID = characterID
	c.spyEditCharacter++
	return c.character, c.err
}

func (c *stubPutCharacter) AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (c *stubPutCharacter) GetCharactersByUserID(userID uuid.UUID) ([]domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (i *stubPutCharacterIndentity) GetUserID(r *http.Request) *uuid.UUID {
	return i.userID
}

func (i *stubPutCharacterIndentity) IsSignedIn(r *http.Request) bool {
	panic("method not used in test")
}

func TestPutCharacterHandler(t *testing.T) {
	t.Run("it should take raw edited character and return updated character", func(t *testing.T) {
		identity := stubPutCharacterIndentity{}
		userID := uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a")
		identity.userID = &userID

		wantCharacterID := uuid.MustParse("11000000-0000-4523-90a2-4868a5bfc90a")

		character := domain.CharacterDTO{RawCharacter: domain.RawCharacter{Name: "TEST-CHARACTER"}}
		characterService := stubPutCharacter{character: &character}

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		formValues := url.Values{
			"name":                         {"UPDATED-CHARACTER"},
			"str":                          {"9"},
			"int":                          {"8"},
			"wil":                          {"7"},
			"hlt":                          {"6"},
			"agi":                          {"5"},
			"tch":                          {"4"},
			"gun_combat_learning_points":   {"4.9"},
			"hand_to_hand_learning_points": {"2.1"},
		}

		request := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/account/characters/%s", wantCharacterID.String()), strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := request.Context()
		ctx = context.WithValue(ctx, "userID", "69000000-0000-4523-90a2-4868a5bfc90a")
		request = request.WithContext(ctx)

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

		if characterService.spyEditCharacter != 1 {
			t.Errorf("got %v calls to EditCharacter want 1", characterService.spyEditCharacter)
		}

		if characterService.spyCharacterID != wantCharacterID {
			t.Errorf("got %v, want %v", characterService.spyCharacterID, wantCharacterID)
		}

		if !strings.Contains(string(body), "<span>TEST-CHARACTER</span>") {
			t.Errorf("expected character fragment, got: %s", body)
		}
	})

	t.Run("it should return error on malformed character id", func(t *testing.T) {
		identity := stubPutCharacterIndentity{}
		characterService := stubPutCharacter{}

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		formValues := url.Values{
			"str":                          {"9"},
			"int":                          {"8"},
			"wil":                          {"7"},
			"hlt":                          {"6"},
			"agi":                          {"5"},
			"tch":                          {"4"},
			"gun_combat_learning_points":   {"4.9"},
			"hand_to_hand_learning_points": {"2.1"},
		}

		request := httptest.NewRequest(http.MethodPut, "/account/characters/bad-id", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusBadRequest)
		}

		if characterService.spyEditCharacter != 0 {
			t.Errorf("got %v calls to EditCharacter want 0", characterService.spyEditCharacter)
		}
	})

	t.Run("it should return error if user id is nil", func(t *testing.T) {
		identity := stubPutCharacterIndentity{}
		characterService := stubPutCharacter{}

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		formValues := url.Values{
			"str":                          {"9"},
			"int":                          {"8"},
			"wil":                          {"7"},
			"hlt":                          {"6"},
			"agi":                          {"5"},
			"tch":                          {"4"},
			"gun_combat_learning_points":   {"4.9"},
			"hand_to_hand_learning_points": {"2.1"},
		}

		request := httptest.NewRequest(http.MethodPut, "/account/characters/11000000-0000-4523-90a2-4868a5bfc90a", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := request.Context()
		ctx = context.WithValue(ctx, "userID", "69000000-0000-4523-90a2-4868a5bfc90a")
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusUnauthorized {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusUnauthorized)
		}

		if characterService.spyEditCharacter != 0 {
			t.Errorf("got %v calls to EditCharacter want 0", characterService.spyEditCharacter)
		}
	})

	t.Run("it should return error when form does not parse", func(t *testing.T) {
		identity := stubPutCharacterIndentity{}
		userID := uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a")
		identity.userID = &userID
		characterService := stubPutCharacter{}

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		formValues := url.Values{
			"name":                         {"UPDATED-CHARACTER"},
			"str":                          {"ABC"},
			"int":                          {"8"},
			"wil":                          {"7"},
			"hlt":                          {"6"},
			"agi":                          {"5"},
			"tch":                          {"4"},
			"gun_combat_learning_points":   {"4.9"},
			"hand_to_hand_learning_points": {"2.1"},
		}

		request := httptest.NewRequest(http.MethodPut, "/account/characters/11000000-0000-4523-90a2-4868a5bfc90a", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := request.Context()
		ctx = context.WithValue(ctx, "userID", "69000000-0000-4523-90a2-4868a5bfc90a")
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusBadRequest)
		}

		if characterService.spyEditCharacter != 0 {
			t.Errorf("got %v calls to EditCharacter want 0", characterService.spyEditCharacter)
		}
	})

	t.Run("it should return error on edit character error", func(t *testing.T) {
		identity := stubPutCharacterIndentity{}
		userID := uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a")
		identity.userID = &userID

		characterService := stubPutCharacter{err: fmt.Errorf("fake character service error")}

		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}
		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		formValues := url.Values{
			"name":                         {"UPDATED-CHARACTER"},
			"str":                          {"9"},
			"int":                          {"8"},
			"wil":                          {"7"},
			"hlt":                          {"6"},
			"agi":                          {"5"},
			"tch":                          {"4"},
			"gun_combat_learning_points":   {"4.9"},
			"hand_to_hand_learning_points": {"2.1"},
		}

		request := httptest.NewRequest(http.MethodPut, "/account/characters/11000000-0000-4523-90a2-4868a5bfc90a", strings.NewReader(formValues.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := request.Context()
		ctx = context.WithValue(ctx, "userID", "69000000-0000-4523-90a2-4868a5bfc90a")
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusInternalServerError {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusInternalServerError)
		}

		if got.Header.Get("Content-Type") != "text/html; charset=utf-8" {
			t.Errorf("got Content-Type %q, want %q", got.Header.Get("Content-Type"), "text/html; charset=utf-8")
		}

		if characterService.spyEditCharacter != 1 {
			t.Errorf("got %v calls to EditCharacter want 1", characterService.spyEditCharacter)
		}
	})
}
