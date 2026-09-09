package domain

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type CharacterDTO struct {
	RawCharacter
	GunCombatLevel  int `json:"gun_combat_level"`
	HandToHandLevel int `json:"hand_to_hand_level"`
}

type RawCharacter struct {
	ID                       uuid.UUID `json:"id"`
	Name                     string    `json:"name"`
	Str                      int       `json:"str"`
	Int                      int       `json:"int"`
	Wil                      int       `json:"wil"`
	Hlt                      int       `json:"hlt"`
	Agi                      int       `json:"agi"`
	Tch                      int       `json:"tch"`
	GunCombatLearningPoints  float32   `json:"gun_combat_learning_points"`
	HandToHandLearningPoints float32   `json:"hand_to_hand_learning_points"`
}

func ParseRawCharacter(form url.Values) (RawCharacter, map[string]string) {
	problems := map[string]string{}

	name := strings.TrimSpace(form.Get("name"))
	rawStr := strings.TrimSpace(form.Get("str"))
	rawItel := strings.TrimSpace(form.Get("int"))
	rawWil := strings.TrimSpace(form.Get("wil"))
	rawHlt := strings.TrimSpace(form.Get("hlt"))
	rawAgi := strings.TrimSpace(form.Get("agi"))
	rawTch := strings.TrimSpace(form.Get("tch"))
	rawGunCombatLearningPoints := strings.TrimSpace(form.Get("gun_combat_learning_points"))
	rawHandToHandLearningPoints := strings.TrimSpace(form.Get("hand_to_hand_learning_points"))

	str, err := strconv.Atoi(rawStr)
	if err != nil {
		problems["str"] = "is not a number"
	} else if str < 1 || str > 21 {
		problems["str"] = "is invalid (must be between 1 and 21)"
	}

	intel, err := strconv.Atoi(rawItel)
	if err != nil {
		problems["int"] = "is not a number"
	} else if intel < 1 || intel > 21 {
		problems["int"] = "is invalid (must be between 1 and 21)"
	}

	wil, err := strconv.Atoi(rawWil)
	if err != nil {
		problems["wil"] = "is not a number"
	} else if wil < 1 || wil > 21 {
		problems["wil"] = "is invalid (must be between 1 and 21)"
	}

	hlt, err := strconv.Atoi(rawHlt)
	if err != nil {
		problems["hlt"] = "is not a number"
	} else if hlt < 1 || hlt > 21 {
		problems["hlt"] = "is invalid (must be between 1 and 21)"
	}

	agi, err := strconv.Atoi(rawAgi)
	if err != nil {
		problems["agi"] = "is not a number"
	} else if agi < 1 || agi > 21 {
		problems["agi"] = "is invalid (must be between 1 and 21)"
	}

	tch, err := strconv.Atoi(rawTch)
	if err != nil {
		problems["tch"] = "is not a number"
	} else if tch < 1 || tch > 21 {
		problems["tch"] = "is invalid (must be between 1 and 21)"
	}

	gunCombatLearningPoints, err := strconv.ParseFloat(rawGunCombatLearningPoints, 32)
	if err != nil {
		problems["gun_combat_learning_points"] = "is not a number"
	} else if gunCombatLearningPoints < 0 || gunCombatLearningPoints > 1834 {
		problems["gun_combat_learning_points"] = "is invalid (must be between 0 and 1834)"
	}

	handToHandLearningPoints, err := strconv.ParseFloat(rawHandToHandLearningPoints, 32)
	if err != nil {
		problems["hand_to_hand_learning_points"] = "is not a number"
	} else if handToHandLearningPoints < 0 || handToHandLearningPoints > 1834 {
		problems["hand_to_hand_learning_points"] = "is invalid (must be between 0 and 1834)"
	}

	c := RawCharacter{
		Name:                     name,
		Str:                      str,
		Int:                      intel,
		Wil:                      wil,
		Hlt:                      hlt,
		Agi:                      agi,
		Tch:                      tch,
		GunCombatLearningPoints:  float32(gunCombatLearningPoints),
		HandToHandLearningPoints: float32(handToHandLearningPoints),
	}

	return c, problems
}
