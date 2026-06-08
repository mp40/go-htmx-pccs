package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/store"
)

type stubStore struct{}

func (s *stubStore) AddCharacter(character store.Character) (*store.Character, error) {
	return &character, nil
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

	t.Run("it converts gun combat level to learning point total", func(t *testing.T) {
		store := &stubStore{}
		service := NewCharacterService(store)

		rawCharacter := domain.RawCharacter{
			GunCombatLevel: 4,
		}

		got, err := service.AddCharacter(uuid.New(), rawCharacter)
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}

		if got.GunCombatLearningPoints != 16 {
			t.Errorf("got %v, want 16", got.GunCombatLearningPoints)
		}
	})

	t.Run("it converts hand to hand level to learning point total", func(t *testing.T) {
		store := &stubStore{}
		service := NewCharacterService(store)

		rawCharacter := domain.RawCharacter{
			HandToHandLevel: 2,
		}

		got, err := service.AddCharacter(uuid.New(), rawCharacter)
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}

		if got.HandToHandLearningPoints != 4 {
			t.Errorf("got %v, want 4", got.HandToHandLearningPoints)
		}
	})
}
