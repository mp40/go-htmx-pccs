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

type stubRouteCharacterEditCharacter struct {
	character      *domain.CharacterDTO
	err            error
	spyCalls       int
	spyCharacterID uuid.UUID
	spyUserID      uuid.UUID
}

func (c *stubRouteCharacterEditCharacter) GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*domain.CharacterDTO, error) {
	c.spyCharacterID = characterID
	c.spyUserID = userID
	c.spyCalls++
	return c.character, c.err
}

func (c *stubRouteCharacterEditCharacter) AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (c *stubRouteCharacterEditCharacter) EditCharacter(userID uuid.UUID, characterID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (c *stubRouteCharacterEditCharacter) GetCharactersByUserID(userID uuid.UUID) ([]domain.CharacterDTO, error) {
	panic("method not used in test")
}

type stubRouteCharacterEditIndentity struct {
	userID     *uuid.UUID
	isSignedIn bool
}

func (i *stubRouteCharacterEditIndentity) GetUserID(r *http.Request) *uuid.UUID {
	return i.userID
}

func (i *stubRouteCharacterEditIndentity) IsSignedIn(r *http.Request) bool {
	return i.isSignedIn
}

func TestGetCharacterEditHandler_Route(t *testing.T) {
	t.Run("it should render edit character page on successful GET request", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEditIndentity{userID: &userID, isSignedIn: true}
		character := domain.CharacterDTO{RawCharacter: domain.RawCharacter{Name: "TEST-CHARACTER"}}
		stubCharacterService := &stubRouteCharacterEditCharacter{character: &character}
		server := NewServer(nil, nil, &identity, stubCharacterService, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/account/characters/69000000-0000-4523-90a2-4868a5bfc90a/edit", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusOK {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusOK)
		}

		if stubCharacterService.spyCharacterID.String() != "69000000-0000-4523-90a2-4868a5bfc90a" {
			t.Errorf("got character id %s want character id %s", stubCharacterService.spyCharacterID.String(), "69000000-0000-4523-90a2-4868a5bfc90a")
		}

		if stubCharacterService.spyUserID != userID {
			t.Errorf("got character id %v want character id %v", stubCharacterService.spyUserID, userID)
		}

		if !strings.Contains(string(body), "<h2>Edit TEST-CHARACTER</h2>") {
			t.Errorf("expected edit character page, got: %s", body)
		}

		if !strings.Contains(string(body), "<body>") {
			t.Errorf("expected page, got: %s", body)
		}
	})

	t.Run("it should render edit character fragment on HTMX GET request", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEditIndentity{userID: &userID, isSignedIn: true}
		character := domain.CharacterDTO{RawCharacter: domain.RawCharacter{Name: "TEST-CHARACTER"}}
		stubCharacterService := &stubRouteCharacterEditCharacter{character: &character}
		server := NewServer(nil, nil, &identity, stubCharacterService, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/account/characters/69000000-0000-4523-90a2-4868a5bfc90a/edit", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusOK {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusOK)
		}

		if stubCharacterService.spyCharacterID.String() != "69000000-0000-4523-90a2-4868a5bfc90a" {
			t.Errorf("got character id %s want character id %s", stubCharacterService.spyCharacterID.String(), "69000000-0000-4523-90a2-4868a5bfc90a")
		}

		if stubCharacterService.spyUserID != userID {
			t.Errorf("got character id %v want character id %v", stubCharacterService.spyUserID, userID)
		}

		if !strings.Contains(string(body), "<h2>Edit TEST-CHARACTER</h2") {
			t.Errorf("expected edit character page, got: %s", body)
		}

		if strings.Contains(string(body), "<body>") {
			t.Errorf("expected fragment, got: %s", body)
		}
	})

	t.Run("it returns bad request message on malformed character id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEditIndentity{userID: &userID, isSignedIn: true}
		stubCharacterService := &stubRouteCharacterEditCharacter{}
		server := NewServer(nil, nil, &identity, stubCharacterService, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/account/characters/bad-id/edit", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusBadRequest {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusBadRequest)
		}

		if stubCharacterService.spyCalls != 0 {
			t.Errorf("got %d calls to character service want 0", stubCharacterService.spyCalls)
		}

		if !strings.Contains(string(body), "<span>bad request: malformed character id</span>") {
			t.Errorf("expected bad request error message, got: %s", body)
		}
	})

	t.Run("it returns unauthorized message on nil user id", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		identity := stubRouteCharacterEditIndentity{isSignedIn: true}
		stubCharacterService := &stubRouteCharacterEditCharacter{}
		server := NewServer(nil, nil, &identity, stubCharacterService, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/account/characters/69000000-0000-4523-90a2-4868a5bfc90a/edit", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusUnauthorized {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusUnauthorized)
		}

		if stubCharacterService.spyCalls != 0 {
			t.Errorf("got %d calls to character service want 0", stubCharacterService.spyCalls)
		}

		if !strings.Contains(string(body), "<span>unauthorised: sign in to edit character</span>") {
			t.Errorf("expected unauthorised error message, got: %s", body)
		}
	})

	t.Run("it should return internal server error message on get character error", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEditIndentity{userID: &userID, isSignedIn: true}
		stubCharacterService := &stubRouteCharacterEditCharacter{err: fmt.Errorf("sadness")}
		server := NewServer(nil, nil, &identity, stubCharacterService, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/account/characters/69000000-0000-4523-90a2-4868a5bfc90a/edit", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusInternalServerError {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusInternalServerError)
		}

		if stubCharacterService.spyCalls != 1 {
			t.Errorf("got %d calls to character service want 1", stubCharacterService.spyCalls)
		}

		if !strings.Contains(string(body), "<span>internal server error</span>") {
			t.Errorf("expected internal server error message, got: %s", body)
		}
	})

	t.Run("it should return not found message on nil character", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEditIndentity{userID: &userID, isSignedIn: true}
		stubCharacterService := &stubRouteCharacterEditCharacter{}
		server := NewServer(nil, nil, &identity, stubCharacterService, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/account/characters/69000000-0000-4523-90a2-4868a5bfc90a/edit", nil)
		request.Header.Set("HX-Request", "true")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		body, err := io.ReadAll(got.Body)
		if err != nil {
			t.Fatal(err)
		}

		if got.StatusCode != http.StatusNotFound {
			t.Errorf("got http status %v want http status %v", got.StatusCode, http.StatusNotFound)
		}

		if stubCharacterService.spyCalls != 1 {
			t.Errorf("got %d calls to character service want 1", stubCharacterService.spyCalls)
		}

		if !strings.Contains(string(body), "<span>character not found</span>") {
			t.Errorf("expected not found message, got: %s", body)
		}
	})
}
