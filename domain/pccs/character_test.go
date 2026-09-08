package pccs

import (
	"fmt"
	"testing"
)

func TestLearningPointsToLevel(t *testing.T) {
	tc := []struct {
		level      int
		lowerLimit float32
		upperLimit float32
	}{
		{level: 0, lowerLimit: 0, upperLimit: 1.9999},
		{level: 1, lowerLimit: 2, upperLimit: 3.9999},
		{level: 2, lowerLimit: 4, upperLimit: 7.9999},
		{level: 3, lowerLimit: 8, upperLimit: 15.9999},
		{level: 4, lowerLimit: 16, upperLimit: 31.9999},
		{level: 5, lowerLimit: 32, upperLimit: 55.9999},
		{level: 6, lowerLimit: 56, upperLimit: 87.9999},
		{level: 7, lowerLimit: 88, upperLimit: 125.9999},
		{level: 8, lowerLimit: 126, upperLimit: 169.9999},
		{level: 9, lowerLimit: 170, upperLimit: 217.9999},
		{level: 10, lowerLimit: 218, upperLimit: 273.9999},
		{level: 11, lowerLimit: 274, upperLimit: 345.9999},
		{level: 12, lowerLimit: 346, upperLimit: 433.9999},
		{level: 13, lowerLimit: 434, upperLimit: 541.9999},
		{level: 14, lowerLimit: 542, upperLimit: 673.9999},
		{level: 15, lowerLimit: 674, upperLimit: 833.9999},
		{level: 16, lowerLimit: 834, upperLimit: 1025.9999},
		{level: 17, lowerLimit: 1026, upperLimit: 1253.9999},
		{level: 18, lowerLimit: 1254, upperLimit: 1551.9999},
		{level: 19, lowerLimit: 1552, upperLimit: 1833.9999},
		{level: 20, lowerLimit: 1834, upperLimit: 9999},
	}

	for _, test := range tc {
		t.Run(fmt.Sprintf("it converts lower limit threshold learning points %f to level %d", test.lowerLimit, test.level), func(t *testing.T) {
			t.Parallel()

			got := ConvertLearningPointsToLevel(test.lowerLimit)
			if got != test.level {
				t.Errorf("got %v, want %v", got, test.level)
			}
		})
		t.Run(fmt.Sprintf("it converts upper limit threshold learning points %f to level %d", test.upperLimit, test.level), func(t *testing.T) {
			t.Parallel()

			got := ConvertLearningPointsToLevel(test.upperLimit)
			if got != test.level {
				t.Errorf("got %v, want %v", got, test.level)
			}
		})
	}
}
