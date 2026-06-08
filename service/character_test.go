package service

import (
	"fmt"
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

func TestConvertLevelToLearningPoints(t *testing.T) {
	tc := []struct {
		level int
		want  float32
	}{
		{level: 0, want: 0},
		{level: 1, want: 2},
		{level: 2, want: 4},
		{level: 3, want: 8},
		{level: 4, want: 16},
		{level: 5, want: 32},
		{level: 6, want: 56},
	}

	for _, test := range tc {
		t.Run(fmt.Sprintf("it converts level %d to %v", test.level, test.want), func(t *testing.T) {
			t.Parallel()

			got := convertLevelToLearningPoints(test.level)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
