package service

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/store"
)

type CharacterService struct {
	store Store
}

type Store interface {
	AddCharacter(character store.Character) (*store.Character, error)
	GetCharactersByUserID(userID uuid.UUID) ([]store.Character, error)
}

func NewCharacterService(store Store) *CharacterService {
	return &CharacterService{
		store: store,
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

func (cs *CharacterService) AddCharacter(userID uuid.UUID, rawCharacter domain.RawCharacter) (*store.Character, error) {
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
		GunCombatLearningPoints:  convertLevelToLearningPoints(rawCharacter.GunCombatLevel),
		HandToHandLearningPoints: convertLevelToLearningPoints(rawCharacter.HandToHandLevel),
	}

	character, err := cs.store.AddCharacter(new)

	return character, err
}

var names = []string{"Leo", "Roy", "Sam", "Joe", "Ben", "Ray", "Avi", "Ian", "Dan", "Tom"}

func generateRandomNumber(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

func getRandomName() string {
	i := generateRandomNumber(0, len(names)-1)
	return fmt.Sprintf("%s-%d", names[i], generateRandomNumber(100, 999))
}

func convertLevelToLearningPoints(level int) float32 {
	// TODO think if we want to return error if level not within map boundry, ie -1 or 21
	lp := domain.LevelToLearningPoints[level]
	return lp
}

func mapStoreCharacterToDomainCharacter(c store.Character) domain.CharacterDTO {
	dto := domain.CharacterDTO{
		RawCharacter: domain.RawCharacter{
			Name:            c.Name,
			Str:             c.Str,
			Int:             c.Int,
			Wil:             c.Wil,
			Hlt:             c.Hlt,
			Agi:             c.Agi,
			Tch:             c.Tch,
			GunCombatLevel:  convertLearningPointsToLevel(c.GunCombatLearningPoints),
			HandToHandLevel: convertLearningPointsToLevel(c.HandToHandLearningPoints),
		},
		GunCombatLearningPoints:  c.GunCombatLearningPoints,
		HandToHandLearningPoints: c.HandToHandLearningPoints,
	}
	return dto
}

func convertLearningPointsToLevel(learningPoints float32) int {
	switch {
	case learningPoints < 2:
		return 0
	case learningPoints < 4:
		return 1
	case learningPoints < 8:
		return 2
	case learningPoints < 16:
		return 3
	case learningPoints < 32:
		return 4
	case learningPoints < 56:
		return 5
	case learningPoints < 88:
		return 6
	case learningPoints < 126:
		return 7
	case learningPoints < 170:
		return 8
	case learningPoints < 218:
		return 9
	case learningPoints < 274:
		return 10
	case learningPoints < 346:
		return 11
	case learningPoints < 434:
		return 12
	case learningPoints < 542:
		return 13
	case learningPoints < 674:
		return 14
	case learningPoints < 834:
		return 15
	case learningPoints < 1026:
		return 16
	case learningPoints < 1254:
		return 17
	case learningPoints < 1552:
		return 18
	case learningPoints < 1834:
		return 19
	default:
		return 20
	}
}
