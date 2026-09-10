package domain

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseRawCharacter(t *testing.T) {
	t.Run("it parses form and returns character", func(t *testing.T) {
		formValues := url.Values{
			"name":                         {" Test-Man   "},
			"str":                          {" 9"},
			"int":                          {"8 "},
			"wil":                          {" 7 "},
			"hlt":                          {"6 "},
			"agi":                          {"5 "},
			"tch":                          {"4"},
			"gun_combat_learning_points":   {"  2"},
			"hand_to_hand_learning_points": {"0 "},
		}
		got, problems := ParseRawCharacter(formValues)
		if len(problems) != 0 {
			t.Errorf("got %d parsing problems, want 0", len(problems))
		}

		want := RawCharacter{
			Name:                     "Test-Man",
			Str:                      9,
			Int:                      8,
			Wil:                      7,
			Hlt:                      6,
			Agi:                      5,
			Tch:                      4,
			GunCombatLearningPoints:  2,
			HandToHandLearningPoints: 0,
		}

		if !cmp.Equal(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
	t.Run("it handles string to integer parsing errors", func(t *testing.T) {
		formValues := url.Values{
			"str":                          {"X"},
			"int":                          {"8"},
			"wil":                          {"7"},
			"hlt":                          {"6"},
			"agi":                          {"5"},
			"tch":                          {"4"},
			"gun_combat_learning_points":   {"2"},
			"hand_to_hand_learning_points": {"0 "},
		}
		_, problems := ParseRawCharacter(formValues)
		if len(problems) != 1 {
			t.Errorf("got %d parsing problems, want 1", len(problems))
		}
	})

	t.Run("it collects all string to integer parsing errors", func(t *testing.T) {
		formValues := url.Values{
			"str":                          {"X"},
			"int":                          {"Y"},
			"wil":                          {"Z"},
			"hlt":                          {"A"},
			"agi":                          {"B"},
			"tch":                          {"C"},
			"gun_combat_learning_points":   {"R"},
			"hand_to_hand_learning_points": {"G"},
		}
		_, problems := ParseRawCharacter(formValues)
		if len(problems) != 8 {
			t.Errorf("got %d parsing problems, want 8", len(problems))
		}
	})

	t.Run("it validates characteristic ranges are not below threshold", func(t *testing.T) {
		tc := []struct{ characteristic string }{
			{characteristic: "str"},
			{characteristic: "int"},
			{characteristic: "wil"},
			{characteristic: "hlt"},
			{characteristic: "agi"},
			{characteristic: "tch"},
		}

		for _, test := range tc {
			t.Run(fmt.Sprintf("it validates %s is not below minimum threshold", test.characteristic), func(t *testing.T) {
				t.Parallel()
				formValues := url.Values{
					"str":                          {"1"},
					"int":                          {"1"},
					"wil":                          {"1"},
					"hlt":                          {"1"},
					"agi":                          {"1"},
					"tch":                          {"1"},
					"gun_combat_learning_points":   {"0"},
					"hand_to_hand_learning_points": {"0"},
				}
				formValues.Set(test.characteristic, "0")

				_, problems := ParseRawCharacter(formValues)
				if !strings.Contains(problems[test.characteristic], "invalid") {
					t.Errorf("expected %s problem key to contain validation message", test.characteristic)
				}
			})
		}
	})

	t.Run("it validates characteristic ranges are not above threshold", func(t *testing.T) {
		tc := []struct{ characteristic string }{
			{characteristic: "str"},
			{characteristic: "int"},
			{characteristic: "wil"},
			{characteristic: "hlt"},
			{characteristic: "agi"},
			{characteristic: "tch"},
		}

		for _, test := range tc {
			t.Run(fmt.Sprintf("it validates %s is not above maximum threshold", test.characteristic), func(t *testing.T) {
				t.Parallel()
				formValues := url.Values{
					"str":                          {"21"},
					"int":                          {"21"},
					"wil":                          {"21"},
					"hlt":                          {"21"},
					"agi":                          {"21"},
					"tch":                          {"21"},
					"gun_combat_learning_points":   {"0"},
					"hand_to_hand_learning_points": {"0"},
				}
				formValues.Set(test.characteristic, "22")

				_, problems := ParseRawCharacter(formValues)
				if !strings.Contains(problems[test.characteristic], "invalid") {
					t.Errorf("expected %s problem key to contain validation message", test.characteristic)
				}
			})
		}
	})

	t.Run("it validates combat learning points are not below zero", func(t *testing.T) {
		tc := []struct{ level string }{
			{level: "gun_combat_learning_points"},
			{level: "hand_to_hand_learning_points"},
		}

		for _, test := range tc {
			t.Run(fmt.Sprintf("it validates %s is not below zero", test.level), func(t *testing.T) {
				t.Parallel()
				formValues := url.Values{
					"str":                          {"10"},
					"int":                          {"10"},
					"wil":                          {"10"},
					"hlt":                          {"10"},
					"agi":                          {"10"},
					"tch":                          {"10"},
					"gun_combat_learning_points":   {"0"},
					"hand_to_hand_learning_points": {"0"},
				}
				formValues.Set(test.level, "-1")

				_, problems := ParseRawCharacter(formValues)
				if !strings.Contains(problems[test.level], "invalid") {
					t.Errorf("expected %s problem key to contain validation message", test.level)
				}
			})
		}
	})

	t.Run("it validates combat learning points are not above 1834", func(t *testing.T) {
		tc := []struct{ level string }{
			{level: "gun_combat_learning_points"},
			{level: "hand_to_hand_learning_points"},
		}

		for _, test := range tc {
			t.Run(fmt.Sprintf("it validates %s is not above 1834", test.level), func(t *testing.T) {
				t.Parallel()
				formValues := url.Values{
					"str":                          {"10"},
					"int":                          {"10"},
					"wil":                          {"10"},
					"hlt":                          {"10"},
					"agi":                          {"10"},
					"tch":                          {"10"},
					"gun_combat_learning_points":   {"1834"},
					"hand_to_hand_learning_points": {"1834"},
				}
				formValues.Set(test.level, "1835")

				_, problems := ParseRawCharacter(formValues)
				if !strings.Contains(problems[test.level], "invalid") {
					t.Errorf("expected %s problem key to contain validation message", test.level)
				}
			})
		}
	})
}

func TestCalculateCharacterCombatData(t *testing.T) {
	t.Run("it calculates combat data of average grune in uniform", func(t *testing.T) {
		raw := RawCharacter{
			Str:                      14,
			Int:                      10,
			Wil:                      10,
			Hlt:                      10,
			Agi:                      12,
			Tch:                      10,
			GunCombatLearningPoints:  8,
			HandToHandLearningPoints: 4,
		}

		dto := CharacterDTO{
			RawCharacter:    raw,
			GunCombatLevel:  4,
			HandToHandLevel: 2,
		}

		enc := CharacterEncumberance{
			Uniform:        "Normal",
			ClothingWeight: 5,
		}

		character := EnchrichedCharacterDTO{
			CharacterDTO:          dto,
			CharacterEncumberance: enc,
		}

		character.enrichWithCombatData()

		wantCombatStats := CharacterCombatStats{
			BaseSpeed:               3,
			MaxSpeed:                7,
			SAL:                     10,
			CE:                      7,
			KnockoutValue:           20,
			HandToHandDamageBonus:   2.5,
			GunCombatActions:        []int{2, 1, 2, 2},
			HandToHandCombatActions: []int{2, 1, 2, 2},
		}

		if !cmp.Equal(character.CharacterCombatStats, wantCombatStats) {
			t.Errorf("got %+v, want %+v", character.CharacterCombatStats, wantCombatStats)
		}
	})
}
