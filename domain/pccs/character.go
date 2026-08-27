package pccs

var LevelToLearningPoints = map[int]float32{
	0:  0,
	1:  2,
	2:  4,
	3:  8,
	4:  16,
	5:  32,
	6:  56,
	7:  88,
	8:  126,
	9:  170,
	10: 218,
	11: 274,
	12: 346,
	13: 434,
	14: 542,
	15: 674,
	16: 834,
	17: 1026,
	18: 1254,
	19: 1552,
	20: 1834,
}

func ConvertLearningPointsToLevel(learningPoints float32) int {
	level := 0
	for l := 1; l <= 20; l++ {
		if learningPoints >= LevelToLearningPoints[l] {
			level = l
		}
	}
	return level
}
