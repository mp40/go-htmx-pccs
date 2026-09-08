package server

import (
	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
)

type stubCharacterService struct{}

func (s *stubCharacterService) GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (s *stubCharacterService) AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (s *stubCharacterService) EditCharacter(userID uuid.UUID, characterID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (s *stubCharacterService) GetCharactersByUserID(userID uuid.UUID) ([]domain.CharacterDTO, error) {
	panic("method not used in test")
}

func (s *stubCharacterService) DeleteUserCharacterByID(userID uuid.UUID, charcaterID uuid.UUID) error {
	panic("method not used in test")
}
