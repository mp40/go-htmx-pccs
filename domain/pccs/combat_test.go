package pccs

import (
	"testing"
)

func TestCalculateBaseSpeed(t *testing.T) {
	t.Run("it should return 4.5 for str 21 and lbs 10", func(t *testing.T) {
		got := calculateBaseSpeed(21, 10)
		want := float32(4.5)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 1.5 for str 1 and lbs 10", func(t *testing.T) {
		got := calculateBaseSpeed(1, 10)
		want := float32(1.5)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 2 for str 21 and lbs 200", func(t *testing.T) {
		got := calculateBaseSpeed(21, 200)
		want := float32(2)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 0 for str 1 and lbs 200", func(t *testing.T) {
		got := calculateBaseSpeed(1, 200)
		want := float32(0)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should round up 15.1 lbs to 20 and return 2 when str is 10", func(t *testing.T) {
		got := calculateBaseSpeed(10, 15.1)
		want := float32(2)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should round up 100.1 lbs to 125 and return 0 when str is 12", func(t *testing.T) {
		got := calculateBaseSpeed(12, 100.1)
		want := float32(0)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 1 for str 10 and lbs 100", func(t *testing.T) {
		got := calculateBaseSpeed(10, 100)
		want := float32(1)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 2 for str 15 and lbs 50", func(t *testing.T) {
		got := calculateBaseSpeed(15, 50)
		want := float32(2)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 2 for str 13 and lbs 24.5", func(t *testing.T) {
		got := calculateBaseSpeed(13, 24.5)
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
			got := calculateBaseSpeed(tc.str, tc.enc)

			if got != tc.want {
				t.Errorf("when str is %v and lbs is %v want %v, got %v", tc.str, tc.enc, tc.want, got)
			}

		}
	})
}

func TestCalculateMaxSpeed(t *testing.T) {
	t.Run("it should return 0 if base speed is 0", func(t *testing.T) {
		got := calculateMaxSpeed(21, 0)
		want := 0

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 2 if agi is 21 and base speed is 1", func(t *testing.T) {
		got := calculateMaxSpeed(21, 1)
		want := 2

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 5 if agi is 10 and base speed is 2.5", func(t *testing.T) {
		got := calculateMaxSpeed(10, 2.5)
		want := 5

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it should return 3 if agi is 1 and base speed is 4.5", func(t *testing.T) {
		got := calculateMaxSpeed(1, 4.5)
		want := 3

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestCalculateActions(t *testing.T) {
	t.Run("it should return 0 actions if max speed is 0", func(t *testing.T) {
		//
	})

	t.Run("it should return 1 for max speed 1 and skill factor 7", func(t *testing.T) {
		//
	})

	t.Run("it should return 24 for max speed 13 and skill factor 39", func(t *testing.T) {
		//
	})

	t.Run("it should return 8 for max speed 7 and skill factor 21", func(t *testing.T) {
		//
	})

	t.Run("it should round down skill factor 10 to 9 and return 3 for max speed 6", func(t *testing.T) {
		//
	})

	t.Run("it ", func(t *testing.T) {
		//
	})
}
