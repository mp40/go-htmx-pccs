package server

import (
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
	character         *domain.CharacterDTO
	err               error
	spyGetCharacterID uuid.UUID
}

func (c *stubRouteCharacterEditCharacter) GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*domain.CharacterDTO, error) {
	c.spyGetCharacterID = characterID
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
	t.Run("it should render edit character page", func(t *testing.T) {
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

		if stubCharacterService.spyGetCharacterID.String() != "69000000-0000-4523-90a2-4868a5bfc90a" {
			t.Errorf("got character id %s want character id %s", stubCharacterService.spyGetCharacterID.String(), "69000000-0000-4523-90a2-4868a5bfc90a")
		}

		if !strings.Contains(string(body), "<span>Edit TEST-CHARACTER</span>") {
			t.Errorf("expected edit character fragment, got: %s", body)
		}
	})
}
