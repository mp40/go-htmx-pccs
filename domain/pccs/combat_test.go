package pccs

import (
	"testing"
)

func TestCalculateBaseSpeed(t *testing.T) {
	t.Run("it should return 4.5 for str 21 and lbs 10", func(t *testing.T) {
		got := CalculateBaseSpeed(21, 10)
		want := float32(4.5)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 1.5 for str 1 and lbs 10", func(t *testing.T) {
		got := CalculateBaseSpeed(1, 10)
		want := float32(1.5)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 2 for str 21 and lbs 200", func(t *testing.T) {
		got := CalculateBaseSpeed(21, 200)
		want := float32(2)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 0 for str 1 and lbs 200", func(t *testing.T) {
		got := CalculateBaseSpeed(1, 200)
		want := float32(0)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should round up 15.1 lbs to 20 and return 2 when str is 10", func(t *testing.T) {
		got := CalculateBaseSpeed(10, 15.1)
		want := float32(2)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should round up 100.1 lbs to 125 and return 0 when str is 12", func(t *testing.T) {
		got := CalculateBaseSpeed(12, 100.1)
		want := float32(0)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 1 for str 10 and lbs 100", func(t *testing.T) {
		got := CalculateBaseSpeed(10, 100)
		want := float32(1)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 2 for str 15 and lbs 50", func(t *testing.T) {
		got := CalculateBaseSpeed(15, 50)
		want := float32(2)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 2 for str 13 and lbs 24.5", func(t *testing.T) {
		got := CalculateBaseSpeed(13, 24.5)
		want := float32(2)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it returns correct value", func(t *testing.T) {
		cases := []struct {
			str  int
			enc  float32
			want float32
		}{
			{str: 21, enc: float32(200), want: float32(2)},
			{str: 20, enc: float32(200), want: float32(2)},
			{str: 19, enc: float32(200), want: float32(1.5)},
			{str: 18, enc: float32(200), want: float32(1.5)},
			{str: 17, enc: float32(200), want: float32(1)},
			{str: 16, enc: float32(200), want: float32(1)},
			{str: 15, enc: float32(150), want: float32(1)},
			{str: 15, enc: float32(150.1), want: float32(0)},
			{str: 14, enc: float32(150), want: float32(1)},
			{str: 14, enc: float32(150.1), want: float32(0)},
			{str: 14, enc: float32(39), want: float32(2)},
			{str: 14, enc: float32(41), want: float32(1.5)},
			{str: 13, enc: float32(20), want: float32(2.5)},
			{str: 13, enc: float32(70), want: float32(1.5)},
			{str: 13, enc: float32(71), want: float32(1)},
			{str: 12, enc: float32(0), want: float32(3)},
			{str: 12, enc: float32(10.01), want: float32(2.5)},
			{str: 12, enc: float32(35), want: float32(2)},
			{str: 12, enc: float32(60), want: float32(1.5)},
			{str: 11, enc: float32(125), want: float32(0)},
			{str: 11, enc: float32(100), want: float32(1)},
			{str: 11, enc: float32(40), want: float32(1.5)},
			{str: 11, enc: float32(35), want: float32(2)},
			{str: 8, enc: float32(125), want: float32(0)},
			{str: 8, enc: float32(70), want: float32(1)},
			{str: 8, enc: float32(35), want: float32(1.5)},
			{str: 8, enc: float32(15), want: float32(2.5)},
			{str: 5, enc: float32(100), want: float32(0)},
			{str: 5, enc: float32(90), want: float32(1)},
			{str: 5, enc: float32(55), want: float32(1)},
			{str: 5, enc: float32(20), want: float32(2)},
			{str: 5, enc: float32(15), want: float32(2.5)},
			{str: 5, enc: float32(10), want: float32(2.5)},
			{str: 4, enc: float32(40), want: float32(1.5)},
			{str: 4, enc: float32(60.1), want: float32(1)},
			{str: 4, enc: float32(70.1), want: float32(0)},
			{str: 3, enc: float32(0), want: float32(2.5)},
			{str: 3, enc: float32(14.1), want: float32(2)},
			{str: 3, enc: float32(19.02), want: float32(1.5)},
			{str: 3, enc: float32(60), want: float32(1)},
			{str: 3, enc: float32(60.1), want: float32(0)},
		}

		for _, tc := range cases {
			got := CalculateBaseSpeed(tc.str, tc.enc)

			if got != tc.want {
				t.Errorf("when str is %v and lbs is %v want %v, got %v", tc.str, tc.enc, tc.want, got)
			}

		}
	})
}

func TestCalculateMaxSpeed(t *testing.T) {
	t.Run("it should return 0 if base speed is 0", func(t *testing.T) {
		got := CalculateMaxSpeed(21, 0)
		want := 0

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 2 if agi is 21 and base speed is 1", func(t *testing.T) {
		got := CalculateMaxSpeed(21, 1)
		want := 2

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 5 if agi is 10 and base speed is 2.5", func(t *testing.T) {
		got := CalculateMaxSpeed(10, 2.5)
		want := 5

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 3 if agi is 1 and base speed is 4.5", func(t *testing.T) {
		got := CalculateMaxSpeed(1, 4.5)
		want := 3

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestCalculateActions(t *testing.T) {
	t.Run("it should return 0 actions if max speed is 0", func(t *testing.T) {
		got := calculateActions(0, 50)
		want := 0

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("it should return 1 for max speed 1 and skill factor 7", func(t *testing.T) {
		got := calculateActions(1, 7)
		want := 1

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("it should return 24 for max speed 13 and skill factor 39", func(t *testing.T) {
		got := calculateActions(13, 39)
		want := 24

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("it should return 8 for max speed 7 and skill factor 21", func(t *testing.T) {
		got := calculateActions(7, 21)
		want := 8

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("it should round down skill factor 10 to 9 and return 3 for max speed 6", func(t *testing.T) {
		got := calculateActions(6, 10)
		want := 3

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("it should round down skill factor 14 to 13 and return 2 for max speed 2", func(t *testing.T) {
		got := calculateActions(2, 14)
		want := 2

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("it should treat skill factors less than 7 as 7", func(t *testing.T) {
		got := calculateActions(1, 6)
		want := 1

		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("it returns correct action count", func(t *testing.T) {
		cases := []struct {
			ms   int
			isf  int
			want int
		}{
			{ms: 1, isf: 29, want: 1},
			{ms: 1, isf: 31, want: 2},
			{ms: 3, isf: 13, want: 2},
			{ms: 3, isf: 15, want: 3},
			{ms: 5, isf: 23, want: 6},
			{ms: 5, isf: 25, want: 7},
			{ms: 10, isf: 29, want: 15},
			{ms: 10, isf: 31, want: 16},
			{ms: 10, isf: 35, want: 17},
			{ms: 10, isf: 37, want: 18},
			{ms: 11, isf: 33, want: 18},
			{ms: 11, isf: 37, want: 19},
		}

		for _, tc := range cases {
			got := calculateActions(tc.ms, tc.isf)

			if got != tc.want {
				t.Errorf("when ms is %v and isf is %v want %v, got %v", tc.ms, tc.isf, tc.want, got)
			}

		}
	})
}

func TestCalculateDamageBonus(t *testing.T) {
	t.Run("it should return 0 for max speed 0", func(t *testing.T) {
		got := CalculateDamageBonus(0, 27)
		want := float32(0)

		if got != want {
			t.Errorf("got %g, want %g", got, want)
		}
	})

	t.Run("it should return 0.5 for max speed 1 and skill factor 7", func(t *testing.T) {
		got := CalculateDamageBonus(1, 7)
		want := float32(0.5)

		if got != want {
			t.Errorf("got %g, want %g", got, want)
		}
	})
	t.Run("it should return 12 for max speed 11 and skill factor 39", func(t *testing.T) {
		got := CalculateDamageBonus(11, 39)
		want := float32(12)

		if got != want {
			t.Errorf("got %g, want %g", got, want)
		}
	})
	t.Run("it should return 2.5 for max speed 7 and skill factor 21", func(t *testing.T) {
		got := CalculateDamageBonus(7, 21)
		want := float32(2.5)

		if got != want {
			t.Errorf("got %g, want %g", got, want)
		}
	})
	t.Run("it should round down skill factor 12 to 11 and return 1 for max speed 6", func(t *testing.T) {
		got := CalculateDamageBonus(6, 11)
		want := float32(1)

		if got != want {
			t.Errorf("got %g, want %g", got, want)
		}
	})
}
