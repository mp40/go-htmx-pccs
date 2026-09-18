package service

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/state"
	"github.com/mp40/go-htmx-pccs/store"
)

type stubStore struct {
	characters   []store.Character
	character    *store.Character
	uniformID    *int
	characterErr error
	uniformErr   error
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
	return s.character, s.characterErr
}

func (s *stubStore) DeleteUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) error {
	return s.characterErr
}

func (s *stubStore) GetUniformIDByCharacterID(characterID uuid.UUID) (*int, error) {
	return s.uniformID, s.uniformErr
}

type stubState struct {
	uniform      *state.Uniform
	err          error
	spyUniformId int
	spyCalls     int
}

func (s *stubState) GetUniformByID(uniformID int) (*state.Uniform, error) {
	s.spyUniformId = uniformID
	s.spyCalls++
	return s.uniform, s.err
}

func TestAddCharacter(t *testing.T) {
	t.Run("it generates random name if none provided", func(t *testing.T) {
		store := &stubStore{}
		state := &stubState{}
		service := NewCharacterService(store, state)

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
		store := &stubStore{characters: []store.Character{{Name: "TEST GRUNT"}}}
		state := &stubState{}
		service := NewCharacterService(store, state)

		got, err := service.GetCharactersByUserID(uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if len(got) != 1 {
			t.Errorf("got %d characters, want 1", len(got))
		}
	})

	t.Run("it maps gun combat learning points to levels", func(t *testing.T) {
		store := &stubStore{characters: []store.Character{{GunCombatLearningPoints: 16}}}
		state := &stubState{}
		service := NewCharacterService(store, state)

		got, err := service.GetCharactersByUserID(uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got[0].GunCombatLevel != 4 {
			t.Errorf("got level %d, want 4", got[0].GunCombatLevel)
		}
	})

	t.Run("it maps hand to hand combat learning points to levels", func(t *testing.T) {
		store := &stubStore{characters: []store.Character{{HandToHandLearningPoints: 4}}}
		state := &stubState{}
		service := NewCharacterService(store, state)

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
		state := &stubState{}
		service := NewCharacterService(store, state)

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
		state := &stubState{}
		service := NewCharacterService(store, state)

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
		state := &stubState{}
		service := NewCharacterService(store, state)

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
		store := &stubStore{characterErr: fmt.Errorf("sadness")}
		state := &stubState{}
		service := NewCharacterService(store, state)

		err := service.DeleteUserCharacterByID(uuid.New(), uuid.New())
		if err == nil {
			t.Errorf("expected error, got <nil>")
		}
	})
}

func TestGetUserEnrichedCharacterByID(t *testing.T) {
	t.Run("it returns an enriched character", func(t *testing.T) {
		c := store.Character{
			Str:                      10,
			Int:                      10,
			Wil:                      10,
			Agi:                      10,
			GunCombatLearningPoints:  4,
			HandToHandLearningPoints: 2,
		}
		id := 1
		store := &stubStore{character: &c, uniformID: &id}
		u := state.Uniform{ID: 1, Name: "FAKE-UNIFORM", Weight: 5.5}
		state := &stubState{uniform: &u}
		service := NewCharacterService(store, state)

		got, err := service.GetUserEnrichedCharacterByID(uuid.New(), uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}

		wantCombatStats := domain.CharacterCombatStats{
			BaseSpeed:               3,
			MaxSpeed:                6,
			SAL:                     7,
			CE:                      5,
			HandToHandDamageBonus:   1.5,
			KnockoutValue:           10,
			GunCombatActions:        []int{2, 1, 2, 1},
			HandToHandCombatActions: []int{2, 1, 1, 1},
		}

		wantEncumbrance := domain.CharacterEncumbrance{
			Uniform:        "FAKE-UNIFORM",
			ClothingWeight: 5.5,
		}

		if !cmp.Equal(got.CharacterCombatStats, wantCombatStats) {
			t.Errorf("unexpected CombatStats, got %+v, want %+v", got.CharacterCombatStats, wantCombatStats)
		}

		if state.spyUniformId != 1 {
			t.Errorf("unexpected uniform id as argument, got %d, want %d", state.spyUniformId, 1)
		}

		if !cmp.Equal(got.CharacterEncumbrance, wantEncumbrance) {
			t.Errorf("unexpected Encumbrance, got %+v, want %+v", got.CharacterEncumbrance, wantEncumbrance)
		}
	})

	t.Run("it defaults to no uniform if not found by character id", func(t *testing.T) {
		c := store.Character{
			Str:                      10,
			Int:                      10,
			Wil:                      10,
			Agi:                      10,
			GunCombatLearningPoints:  4,
			HandToHandLearningPoints: 2,
		}
		store := &stubStore{character: &c}
		state := &stubState{}
		service := NewCharacterService(store, state)

		got, err := service.GetUserEnrichedCharacterByID(uuid.New(), uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}

		wantEncumbrance := domain.CharacterEncumbrance{
			Uniform:        "None",
			ClothingWeight: 0,
		}

		if !cmp.Equal(got.CharacterEncumbrance, wantEncumbrance) {
			t.Errorf("unexpected Encumbrance, got %+v, want %+v", got.CharacterEncumbrance, wantEncumbrance)
		}

		if state.spyCalls != 0 {
			t.Errorf("unexpected call to state: got %d, want 0", state.spyCalls)
		}
	})

	t.Run("it defaults to no uniform if not found by uniform id", func(t *testing.T) {
		c := store.Character{
			Str:                      10,
			Int:                      10,
			Wil:                      10,
			Agi:                      10,
			GunCombatLearningPoints:  4,
			HandToHandLearningPoints: 2,
		}
		id := 1
		store := &stubStore{character: &c, uniformID: &id}
		state := &stubState{}
		service := NewCharacterService(store, state)

		got, err := service.GetUserEnrichedCharacterByID(uuid.New(), uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}

		wantEncumbrance := domain.CharacterEncumbrance{
			Uniform:        "None",
			ClothingWeight: 0,
		}

		if !cmp.Equal(got.CharacterEncumbrance, wantEncumbrance) {
			t.Errorf("unexpected Encumbrance, got %+v, want %+v", got.CharacterEncumbrance, wantEncumbrance)
		}

		if state.spyCalls != 1 {
			t.Errorf("calls to state: got %d, want 1", state.spyCalls)
		}
	})

	t.Run("it returns error on get uniform by character id error", func(t *testing.T) {
		c := store.Character{
			Str:                      10,
			Int:                      10,
			Wil:                      10,
			Agi:                      10,
			GunCombatLearningPoints:  4,
			HandToHandLearningPoints: 2,
		}
		store := &stubStore{character: &c, uniformErr: fmt.Errorf("GET UNIFORM ID SADNESS")}
		state := &stubState{}
		service := NewCharacterService(store, state)

		got, err := service.GetUserEnrichedCharacterByID(uuid.New(), uuid.New())
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if got != nil {
			t.Errorf("expected character to be nil")
		}

		if state.spyCalls != 0 {
			t.Errorf("unexpected call to state: got %d, want 0", state.spyCalls)
		}
	})

	t.Run("it returns error on get uniform by id error", func(t *testing.T) {
		c := store.Character{
			Str:                      10,
			Int:                      10,
			Wil:                      10,
			Agi:                      10,
			GunCombatLearningPoints:  4,
			HandToHandLearningPoints: 2,
		}
		id := 1
		store := &stubStore{character: &c, uniformID: &id}
		state := &stubState{err: fmt.Errorf("GET UNIFORM SADNESS")}
		service := NewCharacterService(store, state)

		got, err := service.GetUserEnrichedCharacterByID(uuid.New(), uuid.New())
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if got != nil {
			t.Errorf("expected character to be nil")
		}

		if state.spyCalls != 1 {
			t.Errorf("unexpected call to state: got %d, want 1", state.spyCalls)
		}
	})
}
