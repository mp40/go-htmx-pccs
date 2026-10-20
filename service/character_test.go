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
	equipment    []store.Equipment
	spyUniformID int
	characterErr error
	uniformErr   error
	equipmentErr error
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

func (s *stubStore) AddUniformIDByCharacterID(uniformID int, characterID uuid.UUID) error {
	s.spyUniformID = uniformID
	return s.uniformErr
}

func (s *stubStore) GetEquipmentByCharacterID(characterID uuid.UUID) ([]store.Equipment, error) {
	return s.equipment, s.equipmentErr
}

func (s *stubStore) AddEquipmentIDByCharacterID(equipmentID int, characterID uuid.UUID, load string) error {
	return s.equipmentErr
}

func (s *stubStore) UpsertEquipmentByCharacterID(equipmentID int, characterID uuid.UUID, load string, qty int) error {
	return s.equipmentErr
}

type stubState struct {
	uniform         *state.Uniform
	equipment       []state.Equipment
	Err             error
	spyUniformId    int
	spyEquipmentIds []int
	spyCalls        int
}

func (s *stubState) GetUniformByID(uniformID int) (*state.Uniform, error) {
	s.spyUniformId = uniformID
	s.spyCalls++
	return s.uniform, s.Err
}

func (s *stubState) GetEquipmentByIDs(equipmentIDs []int) ([]state.Equipment, error) {
	s.spyEquipmentIds = equipmentIDs
	return s.equipment, s.Err
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

	t.Run("it sets uniform to default Normal", func(t *testing.T) {
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

		if store.spyUniformID != defaultUniformID {
			t.Errorf("expected uniform set to id %d, got: %d", defaultUniformID, store.spyUniformID)
		}
	})

	t.Run("it does not block execution on set uniform error", func(t *testing.T) {
		store := &stubStore{uniformErr: fmt.Errorf("sadness")}
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
		zero := 0
		one := 1
		two := 2
		store := &stubStore{character: &c, uniformID: &one, equipment: []store.Equipment{{EquipmentID: 1, QtyPrimary: &one}, {EquipmentID: 2, QtyPrimary: &two, QtySecondary: &one}, {EquipmentID: 3, QtyPrimary: &zero}}}
		u := state.Uniform{ID: 1, Name: "FAKE-UNIFORM", Weight: 5.5}
		state := &stubState{uniform: &u, equipment: []state.Equipment{{ID: 1, Name: "A", Weight: 1.1}, {ID: 2, Name: "B", Weight: 2.2}, {ID: 3, Name: "C", Weight: 33.33}}}
		service := NewCharacterService(store, state)

		got, err := service.GetUserEnrichedCharacterByID(uuid.New(), uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}

		wantCombatStats := domain.CharacterCombatStats{
			BaseSpeed:               2.5,
			MaxSpeed:                5,
			SAL:                     7,
			CE:                      5,
			HandToHandDamageBonus:   1,
			KnockoutValue:           10,
			GunCombatActions:        []int{2, 1, 1, 1},
			HandToHandCombatActions: []int{1, 1, 1, 1},
		}

		wantUniform := domain.Uniform{
			Uniform:        "FAKE-UNIFORM",
			ClothingWeight: 5.5,
		}

		wantEquipment := domain.Equipment{
			Primary: domain.Load{
				Items: []domain.Item{{ID: 1, Name: "A", Weight: 1.1, Qty: 1}, {ID: 2, Name: "B", Weight: 2.2, Qty: 2}, {ID: 3, Name: "C", Weight: 33.33, Qty: 0}},
			},
			Secondary: domain.Load{
				Items: []domain.Item{{ID: 2, Name: "B", Weight: 2.2, Qty: 1}},
			},
		}

		if !cmp.Equal(got.CharacterCombatStats, wantCombatStats) {
			t.Errorf("unexpected CombatStats: %+v", cmp.Diff(got.CharacterCombatStats, wantCombatStats))
		}

		if state.spyUniformId != 1 {
			t.Errorf("unexpected uniform id as argument, got %d, want %d", state.spyUniformId, 1)
		}

		if !cmp.Equal(state.spyEquipmentIds, []int{1, 2, 3}) {
			t.Errorf("unexpected equipment ids as argument: %v", cmp.Diff(state.spyEquipmentIds, []int{1, 2, 3}))
		}

		if !cmp.Equal(got.CharacterEncumbrance.Uniform, wantUniform) {
			t.Errorf("unexpected Uniform: %v", cmp.Diff(got.CharacterEncumbrance.Uniform, wantUniform))
		}

		if !cmp.Equal(got.CharacterEncumbrance.Equipment, wantEquipment) {
			t.Errorf("unexpected Equipment: %v", cmp.Diff(got.CharacterEncumbrance.Equipment, wantEquipment))
		}
	})
}

func TestGetCharacterUniform(t *testing.T) {
	t.Run("it defaults to no uniform if not found by character id", func(t *testing.T) {
		store := &stubStore{}
		state := &stubState{}
		service := NewCharacterService(store, state)

		got, err := service.getCharacterUniform(uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}

		want := domain.Uniform{
			Uniform:        "None",
			ClothingWeight: 0,
		}

		if !cmp.Equal(*got, want) {
			t.Errorf("got %+v, want %+v", *got, want)
		}

		if state.spyCalls != 0 {
			t.Errorf("unexpected call to state: got %d, want 0", state.spyCalls)
		}
	})

	t.Run("it defaults to no uniform if not found by uniform id", func(t *testing.T) {
		id := 1
		store := &stubStore{uniformID: &id}
		state := &stubState{}
		service := NewCharacterService(store, state)

		got, err := service.getCharacterUniform(uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected character not to be nil")
		}

		want := domain.Uniform{
			Uniform:        "None",
			ClothingWeight: 0,
		}

		if !cmp.Equal(*got, want) {
			t.Errorf("got %+v, want %+v", *got, want)
		}

		if state.spyCalls != 1 {
			t.Errorf("calls to state: got %d, want 1", state.spyCalls)
		}
	})

	t.Run("it returns error on get uniform by character id error", func(t *testing.T) {
		store := &stubStore{uniformErr: fmt.Errorf("GET UNIFORM ID SADNESS")}
		state := &stubState{}
		service := NewCharacterService(store, state)

		got, err := service.getCharacterUniform(uuid.New())
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if got != nil {
			t.Errorf("expected uniform to be nil")
		}

		if state.spyCalls != 0 {
			t.Errorf("unexpected call to state: got %d, want 0", state.spyCalls)
		}
	})

	t.Run("it returns error on get uniform by id error", func(t *testing.T) {
		id := 1
		store := &stubStore{uniformID: &id}
		state := &stubState{Err: fmt.Errorf("GET UNIFORM SADNESS")}
		service := NewCharacterService(store, state)

		got, err := service.getCharacterUniform(uuid.New())
		if err == nil {
			t.Errorf("expected error, got nil")
		}

		if got != nil {
			t.Errorf("expected uniform to be nil")
		}

		if state.spyCalls != 1 {
			t.Errorf("unexpected call to state: got %d, want 1", state.spyCalls)
		}
	})
}

func TestGetCharacterEquipment(t *testing.T) {
	t.Run("it defaults to no equipment if not found by character id", func(t *testing.T) {
		store := &stubStore{equipment: []store.Equipment{}}
		state := &stubState{}
		service := NewCharacterService(store, state)

		got, err := service.getCharacterEquipment(uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected equipment not to be nil")
		}

		want := domain.Equipment{
			Primary:   domain.Load{Items: []domain.Item{}},
			Secondary: domain.Load{Items: []domain.Item{}},
		}

		if !cmp.Equal(*got, want) {
			t.Errorf("got %+v, want %+v", *got, want)
		}
	})

	t.Run("it defaults to no equipment item if not found by equipment id", func(t *testing.T) {
		one := 1
		store := &stubStore{equipment: []store.Equipment{{EquipmentID: 1, QtyPrimary: &one}, {EquipmentID: 2, QtyPrimary: &one, QtySecondary: &one}}}
		state := &stubState{equipment: []state.Equipment{{ID: 1, Name: "A", Weight: 1.1}}}
		service := NewCharacterService(store, state)

		got, err := service.getCharacterEquipment(uuid.New())
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		if got == nil {
			t.Errorf("expected equipment not to be nil")
		}

		want := domain.Equipment{
			Primary:   domain.Load{Items: []domain.Item{{ID: 1, Name: "A", Weight: 1.1, Qty: 1}}},
			Secondary: domain.Load{Items: []domain.Item{}},
		}

		if !cmp.Equal(*got, want) {
			t.Errorf("unexpected Equipment: %v", cmp.Diff(*got, want))
		}
	})
}

func TestAddCharacterEquipment(t *testing.T) {
	t.Run("it returns error on db error", func(t *testing.T) {
		store := &stubStore{equipmentErr: fmt.Errorf("sadness")}
		service := NewCharacterService(store, nil)
		_, err := service.AddCharacterEquipment(uuid.New(), uuid.New(), 1, "primary")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}

func TestUpdateCharacterEquipment(t *testing.T) {
	t.Run("it returns error on db error", func(t *testing.T) {
		store := &stubStore{equipmentErr: fmt.Errorf("sadness")}
		service := NewCharacterService(store, nil)
		_, err := service.UpdateCharacterEquipment(uuid.New(), uuid.New(), 1, "primary", 3)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
	})
}
