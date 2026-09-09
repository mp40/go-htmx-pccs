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

type stubCharacterPageCharacterService struct {
	stubCharacterService
	character      *domain.CharacterDTO
	err            error
	spyCalls       int
	spyCharacterID uuid.UUID
	spyUserID      uuid.UUID
}

func (c *stubCharacterPageCharacterService) GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*domain.CharacterDTO, error) {
	c.spyCharacterID = characterID
	c.spyUserID = userID
	c.spyCalls++
	return c.character, c.err
}

func TestGetCharacterHandler(t *testing.T) {
	r, err := render.NewRenderService()
	if err != nil {
		t.Fatalf("render service error %v", err)
	}

	t.Run("it renders full character page on successful GET request", func(t *testing.T) {
		userID := uuid.MustParse("55600000-0000-4523-90a2-4868a5bfc90a")
		identity := stubRouteCharacterEditIndentity{userID: &userID, isSignedIn: true}
		character := domain.CharacterDTO{RawCharacter: domain.RawCharacter{Name: "TEST-CHARACTER"}}
		stubCharacterService := &stubCharacterPageCharacterService{character: &character}
		server := NewServer(nil, nil, &identity, stubCharacterService, nil, r)

		request := httptest.NewRequest(http.MethodGet, "/character/69000000-0000-4523-90a2-4868a5bfc90a", nil)
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
			t.Errorf("got user id %v want user id %v", stubCharacterService.spyUserID, userID)
		}

		if !strings.Contains(string(body), "<h1>TEST-CHARACTER</h1>") {
			t.Errorf("expected character page, got: %s", body)
		}

		if !strings.Contains(string(body), "<!doctype html>") {
			t.Errorf("expected page, got: %s", body)
		}
	})
}
