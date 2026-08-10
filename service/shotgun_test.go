package service

import (
	"testing"
)

func TestConvertSALMToHitLocationSpacing(t *testing.T) {
	cases := []struct {
		salm    int
		spacing int
	}{
		{salm: -99, spacing: 1},
		{salm: -13, spacing: 1},
		{salm: -12, spacing: 2},
		{salm: -11, spacing: 2},
		{salm: -10, spacing: 2},
		{salm: -9, spacing: 2},
		{salm: -8, spacing: 2},
		{salm: -7, spacing: 2},
		{salm: -6, spacing: 3},
		{salm: -5, spacing: 3},
		{salm: -4, spacing: 4},
		{salm: -3, spacing: 4},
		{salm: -2, spacing: 6},
		{salm: -1, spacing: 6},
		{salm: 0, spacing: 8},
		{salm: 1, spacing: 8},
		{salm: 2, spacing: 11},
		{salm: 3, spacing: 11},
		{salm: 4, spacing: 14},
		{salm: 5, spacing: 14},
		{salm: 6, spacing: 19},
		{salm: 7, spacing: 19},
		{salm: 8, spacing: 25},
		{salm: 9, spacing: 25},
		{salm: 10, spacing: 34},
		{salm: 11, spacing: 34},
		{salm: 12, spacing: 45},
		{salm: 13, spacing: 45},
		{salm: 14, spacing: 60},
		{salm: 15, spacing: 60},
		{salm: 16, spacing: 79},
		{salm: 17, spacing: 79},
		{salm: 18, spacing: 100},
		{salm: 99, spacing: 100},
	}

	for _, tc := range cases {
		got := convertSALMToHitLocationSpacing(tc.salm)
		if got != tc.spacing {
			t.Errorf("for salm %d got %d, want %d", tc.salm, got, tc.spacing)
		}
	}
}

func TestGenerateRandomHitsWithinSpread(t *testing.T) {
	t.Run("it gets 3 locations for mid location with tight spread with dice range 100", func(t *testing.T) {
		got := generateRandomHitsWithinSpread(50, -10, 3, RollRangeBasic)

		wantLowest := 48
		wantHighest := 52
		if len(got) != 3 {
			t.Errorf("got %d results, want 3", len(got))
		}

		for _, hit := range got {
			if hit < wantLowest {
				t.Errorf("got %d, want lowest %d", hit, wantLowest)
			}
			if hit > wantHighest {
				t.Errorf("got %d, want highest %d", hit, wantHighest)
			}
		}
	})

	t.Run("it gets 100 locations respecting minimum location 0 for low digit location with wide spread with dice range 100", func(t *testing.T) {
		got := generateRandomHitsWithinSpread(05, 17, 100, RollRangeBasic)

		wantLowest := 0
		wantHighest := 84
		if len(got) != 100 {
			t.Errorf("got %d results, want 100", len(got))
		}

		for _, hit := range got {
			if hit < wantLowest {
				t.Errorf("got %d, want lowest %d", hit, wantLowest)
			}
			if hit > wantHighest {
				t.Errorf("got %d, want highest %d", hit, wantHighest)
			}
		}
	})

	t.Run("it gets 100 locations respecting maximum location 99 for high digit location with wide spread with dice range 100", func(t *testing.T) {
		got := generateRandomHitsWithinSpread(95, 17, 100, RollRangeBasic)

		wantLowest := 16
		wantHighest := 99
		if len(got) != 100 {
			t.Errorf("got %d results, want 100", len(got))
		}

		for _, hit := range got {
			if hit < wantLowest {
				t.Errorf("got %d, want lowest %d", hit, wantLowest)
			}
			if hit > wantHighest {
				t.Errorf("got %d, want highest %d", hit, wantHighest)
			}
		}
	})

	t.Run("it gets 300 locations for mid location with tight spread with dice range 1000", func(t *testing.T) {
		got := generateRandomHitsWithinSpread(500, -10, 300, RollRangeAdvancedOpen)

		wantLowest := 480
		wantHighest := 520
		if len(got) != 300 {
			t.Errorf("got %d results, want 300", len(got))
		}

		for _, hit := range got {
			if hit < wantLowest {
				t.Errorf("got %d, want lowest %d", hit, wantLowest)
			}
			if hit > wantHighest {
				t.Errorf("got %d, want highest %d", hit, wantHighest)
			}
		}
	})

	t.Run("it gets 500 locations respecting minimum location 0 for low digit location with wide spread with dice range 1000", func(t *testing.T) {
		got := generateRandomHitsWithinSpread(05, 17, 500, RollRangeAdvancedOpen)

		wantLowest := 0
		wantHighest := 795
		if len(got) != 500 {
			t.Errorf("got %d results, want 100", len(got))
		}

		for _, hit := range got {
			if hit < wantLowest {
				t.Errorf("got %d, want lowest %d", hit, wantLowest)
			}
			if hit > wantHighest {
				t.Errorf("got %d, want highest %d", hit, wantHighest)
			}
		}
	})

	t.Run("it gets 500 locations respecting maximum location 999 for high digit location with wide spread with dice range 1000", func(t *testing.T) {
		got := generateRandomHitsWithinSpread(995, 17, 500, RollRangeAdvancedOpen)

		wantLowest := 205
		wantHighest := 999
		if len(got) != 500 {
			t.Errorf("got %d results, want 500", len(got))
		}

		for _, hit := range got {
			if hit < wantLowest {
				t.Errorf("got %d, want lowest %d", hit, wantLowest)
			}
			if hit > wantHighest {
				t.Errorf("got %d, want highest %d", hit, wantHighest)
			}
		}
	})
}
