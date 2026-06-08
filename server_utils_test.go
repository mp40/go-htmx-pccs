package main

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mp40/go-htmx-pccs/domain"
)

func TestParseRawCharacter(t *testing.T) {
	t.Run("it parses form and returns character", func(t *testing.T) {
		formValues := url.Values{
			"name":               {" Test-Man   "},
			"str":                {" 9"},
			"int":                {"8 "},
			"wil":                {" 7 "},
			"hlt":                {"6 "},
			"agi":                {"5 "},
			"tch":                {"4"},
			"gun_combat_level":   {"  1"},
			"hand_to_hand_level": {"0 "},
		}
		got, problems := parseRawCharacter(formValues)
		if len(problems) != 0 {
			t.Errorf("got %d parsing problems, want 0", len(problems))
		}

		want := domain.RawCharacter{
			Name:            "Test-Man",
			Str:             9,
			Int:             8,
			Wil:             7,
			Hlt:             6,
			Agi:             5,
			Tch:             4,
			GunCombatLevel:  1,
			HandToHandLevel: 0,
		}

		if !cmp.Equal(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
	t.Run("it handles string to integer parsing errors", func(t *testing.T) {
		formValues := url.Values{
			"str":                {"X"},
			"int":                {"8"},
			"wil":                {"7"},
			"hlt":                {"6"},
			"agi":                {"5"},
			"tch":                {"4"},
			"gun_combat_level":   {"1"},
			"hand_to_hand_level": {"0"},
		}
		_, problems := parseRawCharacter(formValues)
		if len(problems) != 1 {
			t.Errorf("got %d parsing problems, want 1", len(problems))
		}
	})

	t.Run("it collects all string to integer parsing errors", func(t *testing.T) {
		formValues := url.Values{
			"str":                {"X"},
			"int":                {"Y"},
			"wil":                {"Z"},
			"hlt":                {"A"},
			"agi":                {"B"},
			"tch":                {"C"},
			"gun_combat_level":   {"R"},
			"hand_to_hand_level": {"G"},
		}
		_, problems := parseRawCharacter(formValues)
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
					"str":                {"1"},
					"int":                {"1"},
					"wil":                {"1"},
					"hlt":                {"1"},
					"agi":                {"1"},
					"tch":                {"1"},
					"gun_combat_level":   {"0"},
					"hand_to_hand_level": {"0"},
				}
				formValues.Set(test.characteristic, "0")

				_, problems := parseRawCharacter(formValues)
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
					"str":                {"22"},
					"int":                {"22"},
					"wil":                {"22"},
					"hlt":                {"22"},
					"agi":                {"22"},
					"tch":                {"22"},
					"gun_combat_level":   {"0"},
					"hand_to_hand_level": {"0"},
				}
				formValues.Set(test.characteristic, "22")

				_, problems := parseRawCharacter(formValues)
				if !strings.Contains(problems[test.characteristic], "invalid") {
					t.Errorf("expected %s problem key to contain validation message", test.characteristic)
				}
			})
		}
	})

	t.Run("it validates combat level is not below zero", func(t *testing.T) {
		tc := []struct{ level string }{
			{level: "gun_combat_level"},
			{level: "hand_to_hand_level"},
		}

		for _, test := range tc {
			t.Run(fmt.Sprintf("it validates %s is not below zero", test.level), func(t *testing.T) {
				t.Parallel()
				formValues := url.Values{
					"str":                {"10"},
					"int":                {"10"},
					"wil":                {"10"},
					"hlt":                {"10"},
					"agi":                {"10"},
					"tch":                {"10"},
					"gun_combat_level":   {"0"},
					"hand_to_hand_level": {"0"},
				}
				formValues.Set(test.level, "-1")

				_, problems := parseRawCharacter(formValues)
				if !strings.Contains(problems[test.level], "invalid") {
					t.Errorf("expected %s problem key to contain validation message", test.level)
				}
			})
		}
	})

	t.Run("it validates combat level is not above 20", func(t *testing.T) {
		tc := []struct{ level string }{
			{level: "gun_combat_level"},
			{level: "hand_to_hand_level"},
		}

		for _, test := range tc {
			t.Run(fmt.Sprintf("it validates %s is not below zero", test.level), func(t *testing.T) {
				t.Parallel()
				formValues := url.Values{
					"str":                {"10"},
					"int":                {"10"},
					"wil":                {"10"},
					"hlt":                {"10"},
					"agi":                {"10"},
					"tch":                {"10"},
					"gun_combat_level":   {"20"},
					"hand_to_hand_level": {"20"},
				}
				formValues.Set(test.level, "21")

				_, problems := parseRawCharacter(formValues)
				if !strings.Contains(problems[test.level], "invalid") {
					t.Errorf("expected %s problem key to contain validation message", test.level)
				}
			})
		}
	})
}
