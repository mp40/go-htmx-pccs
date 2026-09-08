package service

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/store"
)

type stubStore struct {
	characters []store.Character
	character  *store.Character
	err        error
}

func (s *stubStore) AddCharacter(character store.Character) (*store.Character, error) {
	return &character, nil
}

func (s *stubStore) UpdateCharacter(character store.Character) (*store.Character, error) {
	return &character, nil
}

func (s *stubStore) GetCharactersByUserID(userID uuid.UUID) ([]store.Character, error) {
	return s.characters, nil
}

func (s *stubStore) GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*store.Character, error) {
	return s.character, s.err
}

func (s *stubStore) DeleteUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) error {
	return s.err
}

func TestAddCharacter(t *testing.T) {
	t.Run("it generates random name if none provided", func(t *testing.T) {
		store := &stubStore{}
		service := NewCharacterService(store)

		rawCharacter := domain.RawCharacter{}

		got, err := service.AddCharacter(uuid.New(), rawCharacter)
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}

		if len(got.Name) == 0 {
			t.Errorf("expected length of name to be greater than 0")
		}
	})
}

func TestGetCharactersByUserID(t *testing.T) {
	t.Run("it returns customers", func(t *testing.T) {
		s := &stubStore{}
		s.characters = []store.Character{{Name: "TEST GRUNT"}}
		service := NewCharacterService(s)

		got, err := service.GetCharactersByUserID(uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if len(got) != 1 {
			t.Errorf("got %d characters, want 1", len(got))
		}
	})

	t.Run("it maps gun combat learning points to levels", func(t *testing.T) {
		s := &stubStore{}
		s.characters = []store.Character{{GunCombatLearningPoints: 16}}
		service := NewCharacterService(s)

		got, err := service.GetCharactersByUserID(uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got[0].GunCombatLevel != 4 {
			t.Errorf("got level %d, want 4", got[0].GunCombatLevel)
		}
	})

	t.Run("it maps hand to hand combat learning points to levels", func(t *testing.T) {
		s := &stubStore{}
		s.characters = []store.Character{{HandToHandLearningPoints: 4}}
		service := NewCharacterService(s)

		got, err := service.GetCharactersByUserID(uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got[0].HandToHandLevel != 2 {
			t.Errorf("got level %d, want 2", got[0].HandToHandLevel)
		}
	})
}

func TestEditCharacter(t *testing.T) {
	t.Run("it updates character", func(t *testing.T) {
		store := &stubStore{}
		service := NewCharacterService(store)

		rawCharacter := domain.RawCharacter{}

		got, err := service.EditCharacter(uuid.New(), uuid.New(), rawCharacter)
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}
	})
}

func TestGetUserCharacterByID(t *testing.T) {
	t.Run("it gets character by character id and user id", func(t *testing.T) {
		store := &stubStore{character: &store.Character{Name: "TEST-CHARACTER"}}
		service := NewCharacterService(store)

		got, err := service.GetUserCharacterByID(uuid.New(), uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}
	})

	t.Run("it returns nil character when not found", func(t *testing.T) {
		store := &stubStore{}
		service := NewCharacterService(store)

		got, err := service.GetUserCharacterByID(uuid.New(), uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got != nil {
			t.Errorf("expected character to be nil")
		}
	})
}

func TestDeleteUserCharacterByID(t *testing.T) {
	t.Run("it handles errors", func(t *testing.T) {
		store := &stubStore{err: fmt.Errorf("sadness")}
		service := NewCharacterService(store)

		err := service.DeleteUserCharacterByID(uuid.New(), uuid.New())
		if err == nil {
			t.Errorf("expected error, got <nil>")
		}
	})
}
