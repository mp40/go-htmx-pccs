package domain

type RawCharacter struct {
	Name string `json:"name"`
	Str  int    `json:"str"`
	Int  int    `json:"int"`
	Wil  int    `json:"wil"`
	Hlt  int    `json:"hlt"`
	Agi  int    `json:"agi"`
}
