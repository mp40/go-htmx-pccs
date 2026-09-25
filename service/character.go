package service

import (
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/domain/pccs"
	"github.com/mp40/go-htmx-pccs/state"
	"github.com/mp40/go-htmx-pccs/store"
)

var defaultUniformID = 2

type CharacterService struct {
	store characterStore
	state characterState
}

type characterStore interface {
	AddCharacter(character store.Character) (*store.Character, error)
	UpdateCharacter(character store.Character) (*store.Character, error)
	GetCharactersByUserID(userID uuid.UUID) ([]store.Character, error)
	GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*store.Character, error)
	DeleteUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) error
	GetUniformIDByCharacterID(characterID uuid.UUID) (*int, error)
	AddUniformIDByCharacterID(uniformID int, characterID uuid.UUID) error
}

type characterState interface {
	GetUniformByID(uniformID int) (*state.Uniform, error)
}

func NewCharacterService(store characterStore, state characterState) *CharacterService {
	return &CharacterService{
		store: store,
		state: state,
	}
}

func (cs *CharacterService) GetCharactersByUserID(userID uuid.UUID) ([]domain.CharacterDTO, error) {
	characters, err := cs.store.GetCharactersByUserID(userID)
	if err != nil {
		return []domain.CharacterDTO{}, err
	}

	result := make([]domain.CharacterDTO, len(characters))
	for i, c := range characters {
		dto := mapStoreCharacterToDomainCharacter(c)
		result[i] = dto
	}
	return result, err
}

func (cs *CharacterService) AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*domain.CharacterDTO, error) {
	n := rawCharacter.Name
	if len(n) == 0 {
		n = getRandomName()
	}
	new := store.Character{
		ID:                       uuid.New(),
		UserID:                   userID,
		Name:                     n,
		Str:                      rawCharacter.Str,
		Int:                      rawCharacter.Int,
		Wil:                      rawCharacter.Wil,
		Hlt:                      rawCharacter.Hlt,
		Agi:                      rawCharacter.Agi,
		Tch:                      rawCharacter.Tch,
		GunCombatLearningPoints:  rawCharacter.GunCombatLearningPoints,
		HandToHandLearningPoints: rawCharacter.HandToHandLearningPoints,
	}

	character, err := cs.store.AddCharacter(new)
	if character == nil || err != nil {
		return nil, err
	}

	uniformErr := cs.store.AddUniformIDByCharacterID(defaultUniformID, character.ID)
	if uniformErr != nil {
		slog.Error("could not set default uniform by id", "err", uniformErr)
	}

	dto := mapStoreCharacterToDomainCharacter(*character)
	return &dto, err
}

func (cs *CharacterService) EditCharacter(userID uuid.UUID, characterID uuid.UUID, rawCharacterEdit domain.RawCharacter) (*domain.CharacterDTO, error) {
	updated := store.Character{
		ID:                       characterID,
		UserID:                   userID,
		Name:                     rawCharacterEdit.Name,
		Str:                      rawCharacterEdit.Str,
		Int:                      rawCharacterEdit.Int,
		Wil:                      rawCharacterEdit.Wil,
		Hlt:                      rawCharacterEdit.Hlt,
		Agi:                      rawCharacterEdit.Agi,
		Tch:                      rawCharacterEdit.Tch,
		GunCombatLearningPoints:  rawCharacterEdit.GunCombatLearningPoints,
		HandToHandLearningPoints: rawCharacterEdit.HandToHandLearningPoints,
	}

	character, err := cs.store.UpdateCharacter(updated)
	if character == nil || err != nil {
		return nil, err
	}
	dto := mapStoreCharacterToDomainCharacter(*character)
	return &dto, err
}

func (cs *CharacterService) GetUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*domain.CharacterDTO, error) {
	character, err := cs.store.GetUserCharacterByID(userID, characterID)
	if character == nil || err != nil {
		return nil, err
	}
	dto := mapStoreCharacterToDomainCharacter(*character)
	return &dto, err
}

func (cs *CharacterService) DeleteUserCharacterByID(userID uuid.UUID, characterID uuid.UUID) error {
	return cs.store.DeleteUserCharacterByID(userID, characterID)
}

func (cs *CharacterService) GetUserEnrichedCharacterByID(userID uuid.UUID, characterID uuid.UUID) (*domain.EnrichedCharacterDTO, error) {
	character, err := cs.store.GetUserCharacterByID(userID, characterID)
	if character == nil || err != nil {
		return nil, err
	}

	uniformID, err := cs.store.GetUniformIDByCharacterID(characterID)
	if err != nil {
		return nil, err
	}

	var u *state.Uniform
	if uniformID == nil {
		slog.Warn("character has no uniform assigned", "character_id", characterID)
	} else {
		u, err = cs.state.GetUniformByID(*uniformID)
		if err != nil {
			slog.Warn("could not get uniform by id", "uniform_id", *uniformID, "err", err)
			return nil, err
		}
	}

	dto := mapStoreCharacterToDomainCharacter(*character)
	encumbrance := mapStateUniformToDomainEncumbrance(u)

	enrichedCharacter := domain.EnrichCharacter(dto, encumbrance)

	return &enrichedCharacter, nil
}

var names = []string{"Leo", "Roy", "Sam", "Joe", "Ben", "Ray", "Avi", "Ian", "Dan", "Tom"}

func generateRandomNumber(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func getRandomName() string {
	i := generateRandomNumber(0, len(names)-1)
	return fmt.Sprintf("%s-%d", names[i], generateRandomNumber(100, 999))
}

func mapStoreCharacterToDomainCharacter(c store.Character) domain.CharacterDTO {
	dto := domain.CharacterDTO{
		RawCharacter: domain.RawCharacter{
			ID:                       c.ID,
			Name:                     c.Name,
			Str:                      c.Str,
			Int:                      c.Int,
			Wil:                      c.Wil,
			Hlt:                      c.Hlt,
			Agi:                      c.Agi,
			Tch:                      c.Tch,
			GunCombatLearningPoints:  c.GunCombatLearningPoints,
			HandToHandLearningPoints: c.HandToHandLearningPoints,
		},
		GunCombatLevel:  pccs.ConvertLearningPointsToLevel(c.GunCombatLearningPoints),
		HandToHandLevel: pccs.ConvertLearningPointsToLevel(c.HandToHandLearningPoints),
	}
	return dto
}

func mapStateUniformToDomainEncumbrance(uniform *state.Uniform) domain.CharacterEncumbrance {
	encumbrance := domain.CharacterEncumbrance{
		Uniform:        "None",
		ClothingWeight: 0,
	}
	if uniform == nil {
		return encumbrance
	}

	encumbrance.Uniform = uniform.Name
	encumbrance.ClothingWeight = uniform.Weight

	return encumbrance
}
