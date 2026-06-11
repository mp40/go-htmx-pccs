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
}

func (s *stubStore) AddCharacter(character store.Character) (*store.Character, error) {
	return &character, nil
}

func (s *stubStore) GetCharactersByUserID(userID uuid.UUID) ([]store.Character, error) {
	return s.characters, nil
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
		{level: 7, want: 88},
		{level: 8, want: 126},
		{level: 9, want: 170},
		{level: 10, want: 218},
		{level: 11, want: 274},
		{level: 12, want: 346},
		{level: 13, want: 434},
		{level: 14, want: 542},
		{level: 15, want: 674},
		{level: 16, want: 834},
		{level: 17, want: 1026},
		{level: 18, want: 1254},
		{level: 19, want: 1552},
		{level: 20, want: 1834},
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

func TestLearningPointsToLevel(t *testing.T) {
	tc := []struct {
		level      int
		lowerLimit float32
		upperLimit float32
	}{
		{level: 0, lowerLimit: 0, upperLimit: 1.9999},
		{level: 1, lowerLimit: 2, upperLimit: 3.9999},
		{level: 2, lowerLimit: 4, upperLimit: 7.9999},
		{level: 3, lowerLimit: 8, upperLimit: 15.9999},
		{level: 4, lowerLimit: 16, upperLimit: 31.9999},
		{level: 5, lowerLimit: 32, upperLimit: 55.9999},
		{level: 6, lowerLimit: 56, upperLimit: 87.9999},
		{level: 7, lowerLimit: 88, upperLimit: 125.9999},
		{level: 8, lowerLimit: 126, upperLimit: 169.9999},
		{level: 9, lowerLimit: 170, upperLimit: 217.9999},
		{level: 10, lowerLimit: 218, upperLimit: 273.9999},
		{level: 11, lowerLimit: 274, upperLimit: 345.9999},
		{level: 12, lowerLimit: 346, upperLimit: 433.9999},
		{level: 13, lowerLimit: 434, upperLimit: 541.9999},
		{level: 14, lowerLimit: 542, upperLimit: 673.9999},
		{level: 15, lowerLimit: 674, upperLimit: 833.9999},
		{level: 16, lowerLimit: 834, upperLimit: 1025.9999},
		{level: 17, lowerLimit: 1026, upperLimit: 1253.9999},
		{level: 18, lowerLimit: 1254, upperLimit: 1551.9999},
		{level: 19, lowerLimit: 1552, upperLimit: 1833.9999},
		{level: 20, lowerLimit: 1834, upperLimit: 9999},
	}

	for _, test := range tc {
		t.Run(fmt.Sprintf("it converts lower limit threshold learning points %f to level %d", test.lowerLimit, test.level), func(t *testing.T) {
			t.Parallel()

			got := convertLearningPointsToLevel(test.lowerLimit)
			if got != test.level {
				t.Errorf("got %v, want %v", got, test.level)
			}
		})
		t.Run(fmt.Sprintf("it converts upper limit threshold learning points %f to level %d", test.lowerLimit, test.level), func(t *testing.T) {
			t.Parallel()

			got := convertLearningPointsToLevel(test.upperLimit)
			if got != test.level {
				t.Errorf("got %v, want %v", got, test.level)
			}
		})
	}
}
