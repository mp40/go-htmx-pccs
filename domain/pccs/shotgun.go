package pccs

import (
	"math/rand/v2"
)

type RollRange int

const (
	RollRangeBasic        RollRange = 100
	RollRangeAdvancedOpen RollRange = 1000
)

func convertSALMToHitLocationSpacing(salm int) int {
	switch {
	case salm < -12:
		return 1
	case salm <= -7:
		return 2
	case salm <= -5:
		return 3
	case salm <= -3:
		return 4
	case salm <= -1:
		return 6
	case salm <= 1:
		return 8
	case salm <= 3:
		return 11
	case salm <= 5:
		return 14
	case salm <= 7:
		return 19
	case salm <= 9:
		return 25
	case salm <= 11:
		return 34
	case salm <= 13:
		return 45
	case salm <= 15:
		return 60
	case salm <= 17:
		return 79
	default:
		return 100
	}
}

func GenerateRandomHitsWithinSpread(initialLocation int, salm int, hitCount int, rollRange RollRange) []int {
	hits := []int{}

	spacingPercent := convertSALMToHitLocationSpacing(salm)
	spacing := spacingPercent * int(rollRange) / 100

	maxRoll := int(rollRange) - 1

	min := initialLocation - spacing
	if min < 0 {
		min = 0
	}
	max := initialLocation + spacing
	if max > maxRoll {
		max = maxRoll
	}

	for i := 0; i < hitCount; i++ {
		random := rand.IntN(max-min+1) + min
		hits = append(hits, random)
	}

	return hits
}
