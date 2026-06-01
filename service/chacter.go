package service

import (
	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/store"
)

type CharacterService struct {
	store Store
}

type Store interface {
	AddCharacter(character store.Character) (*store.Character, error)
}

func NewCharacterService(store Store) *CharacterService {
	return &CharacterService{
		store: store,
	}
}

func (cs *CharacterService) AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*store.Character, error) {
	new := store.Character{
		ID:     uuid.New(),
		UserID: userID,
		Str:    rawCharacter.Str,
		Int:    rawCharacter.Int,
		Wil:    rawCharacter.Wil,
		Hlt:    rawCharacter.Hlt,
		Agi:    rawCharacter.Agi,
	}

	character, err := cs.store.AddCharacter(new)

	return character, err
}
