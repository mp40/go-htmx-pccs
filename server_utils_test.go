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
			"name": {" Test-Man   "},
			"str":  {" 9"},
			"int":  {"8 "},
			"wil":  {" 7 "},
			"hlt":  {"6 "},
			"agi":  {"5 "},
		}
		got, problems := parseRawCharacter(formValues)
		if len(problems) != 0 {
			t.Errorf("got %d parsing problems, want 0", len(problems))
		}

		want := domain.RawCharacter{
			Name: "Test-Man",
			Str:  9,
			Int:  8,
			Wil:  7,
			Hlt:  6,
			Agi:  5,
		}

		if !cmp.Equal(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
	t.Run("it handles string to integer parsing errors", func(t *testing.T) {
		formValues := url.Values{
			"str": {"X"},
			"int": {"8"},
			"wil": {"7"},
			"hlt": {"6"},
			"agi": {"5"},
		}
		_, problems := parseRawCharacter(formValues)
		if len(problems) != 1 {
			t.Errorf("got %d parsing problems, want 1", len(problems))
		}
	})

	t.Run("it collects all string to integer parsing errors", func(t *testing.T) {
		formValues := url.Values{
			"str": {"X"},
			"int": {"Y"},
			"wil": {"Z"},
			"hlt": {"A"},
			"agi": {"B"},
		}
		_, problems := parseRawCharacter(formValues)
		if len(problems) != 5 {
			t.Errorf("got %d parsing problems, want 5", len(problems))
		}
	})

	t.Run("it validates attribute ranges are not below threshold", func(t *testing.T) {
		tc := []struct{ attribute string }{
			{attribute: "str"},
			{attribute: "int"},
			{attribute: "wil"},
			{attribute: "hlt"},
			{attribute: "agi"},
		}

		for _, test := range tc {
			t.Run(fmt.Sprintf("it validates %s is not below minimum threshold", test.attribute), func(t *testing.T) {
				t.Parallel()
				formValues := url.Values{
					"str": {"9"},
					"int": {"8"},
					"wil": {"7"},
					"hlt": {"6"},
					"agi": {"5"},
				}
				formValues.Set(test.attribute, "0")

				_, problems := parseRawCharacter(formValues)
				if !strings.Contains(problems[test.attribute], "invalid") {
					t.Errorf("expected %s problem key to contain validation message", test.attribute)
				}
			})
		}
	})

	t.Run("it validates attribute ranges are not above threshold", func(t *testing.T) {
		tc := []struct{ attribute string }{
			{attribute: "str"},
			{attribute: "int"},
			{attribute: "wil"},
			{attribute: "hlt"},
			{attribute: "agi"},
		}

		for _, test := range tc {
			t.Run(fmt.Sprintf("it validates %s is not above maximum threshold", test.attribute), func(t *testing.T) {
				t.Parallel()
				formValues := url.Values{
					"str": {"9"},
					"int": {"8"},
					"wil": {"7"},
					"hlt": {"6"},
					"agi": {"5"},
				}
				formValues.Set(test.attribute, "22")

				_, problems := parseRawCharacter(formValues)
				if !strings.Contains(problems[test.attribute], "invalid") {
					t.Errorf("expected %s problem key to contain validation message", test.attribute)
				}
			})
		}
	})
}
