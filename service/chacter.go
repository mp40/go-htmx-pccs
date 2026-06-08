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
}

func NewCharacterService(store Store) *CharacterService {
	return &CharacterService{
		store: store,
	}
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
