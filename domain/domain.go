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

var LevelToLearningPoints = map[int]float32{
	0:  0,
	1:  2,
	2:  4,
	3:  8,
	4:  16,
	5:  32,
	6:  56,
	7:  88,
	8:  126,
	9:  170,
	10: 218,
	11: 274,
	12: 346,
	13: 434,
	14: 542,
	15: 674,
	16: 834,
	17: 1026,
	18: 1254,
	19: 1552,
	20: 1834,
}
