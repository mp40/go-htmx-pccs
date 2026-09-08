package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/render"
)

type stubDeleteCharacter struct {
	stubCharacterService
	err            error
	spyCharacterID uuid.UUID
}

type stubDeleteCharacterIndentity struct {
	userID *uuid.UUID
}

func (s *stubDeleteCharacter) DeleteUserCharacterByID(userID uuid.UUID, charcaterID uuid.UUID) error {
	s.spyCharacterID = charcaterID
	return s.err
}

func (i *stubDeleteCharacterIndentity) GetUserID(r *http.Request) *uuid.UUID {
	return i.userID
}

func (i *stubDeleteCharacterIndentity) IsSignedIn(r *http.Request) bool {
	panic("method not used in test")
}

func TestDeleteCharacterHandler(t *testing.T) {
	t.Run("it should delete character", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		identity := stubDeleteCharacterIndentity{}
		userID := uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a")
		identity.userID = &userID

		characterService := stubDeleteCharacter{}

		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/account/characters/11000000-0000-4523-90a2-4868a5bfc90a", nil)
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		ctx := request.Context()
		ctx = context.WithValue(ctx, "userID", "69000000-0000-4523-90a2-4868a5bfc90a")
		request = request.WithContext(ctx)

		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)

		got := response.Result()

		if got.StatusCode != http.StatusOK {
			t.Errorf("got %v want %v", got.StatusCode, http.StatusOK)
		}

		if characterService.spyCharacterID.String() != "11000000-0000-4523-90a2-4868a5bfc90a" {
			t.Errorf("got %v want %v", characterService.spyCharacterID.String(), "11000000-0000-4523-90a2-4868a5bfc90a")
		}
	})

	t.Run("it should return error message on delete error", func(t *testing.T) {
		r, err := render.NewRenderService()
		if err != nil {
			t.Fatalf("render service error %v", err)
		}

		identity := stubDeleteCharacterIndentity{}
		userID := uuid.MustParse("69000000-0000-4523-90a2-4868a5bfc90a")
		identity.userID = &userID

		characterService := stubDeleteCharacter{err: fmt.Errorf("sadness")}

		server := NewServer(nil, nil, &identity, &characterService, nil, r)

		request := httptest.NewRequest(http.MethodDelete, "/account/characters/11000000-0000-4523-90a2-4868a5bfc90a", nil)
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
	})
}
