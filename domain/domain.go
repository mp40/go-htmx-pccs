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
	HandToHandLevel int    `json:"hand_to_Hand_level"`
}
