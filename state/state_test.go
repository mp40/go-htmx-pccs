package state

import (
	"fmt"
	"testing"
)

func TestCreatePlaceholdersFromSliceLength(t *testing.T) {
	cases := []struct {
		length int
		want   string
	}{
		{length: 0, want: ""},
		{length: 1, want: "?"},
		{length: 2, want: "?,?"},
		{length: 3, want: "?,?,?"},
		{length: 4, want: "?,?,?,?"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("it creates placeholder sql string for length %d", tc.length), func(t *testing.T) {
			got := createPlaceholdersFromSliceLength(tc.length)
			if got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}
