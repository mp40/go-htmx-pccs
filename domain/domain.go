package domain

type RawCharacter struct {
	Name            string `json:"name"`
	Str             int    `json:"str"`
	Int             int    `json:"int"`
	Wil             int    `json:"wil"`
	Hlt             int    `json:"hlt"`
	Agi             int    `json:"agi"`
	Tch             int    `json:"tch"`
	GunCombatLevel  int    `json:"gun_combat_level"`
	HandToHandLevel int    `json:"hand_to_hand_level"`
}

type CharacterDTO struct {
	RawCharacter
	GunCombatLearningPoints  float32 `json:"gun_combat_learning_points"`
	HandToHandLearningPoints float32 `json:"hand_to_hand_learning_points"`
}

type RawCharacterEdit struct {
	Name                     string  `json:"name"`
	Str                      int     `json:"str"`
	Int                      int     `json:"int"`
	Wil                      int     `json:"wil"`
	Hlt                      int     `json:"hlt"`
	Agi                      int     `json:"agi"`
	Tch                      int     `json:"tch"`
	GunCombatLearningPoints  float32 `json:"gun_combat_learning_points"`
	HandToHandLearningPoints float32 `json:"hand_to_hand_learning_points"`
}
