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
	switch {
	case learningPoints < 2:
		return 0
	case learningPoints < 4:
		return 1
	case learningPoints < 8:
		return 2
	case learningPoints < 16:
		return 3
	case learningPoints < 32:
		return 4
	case learningPoints < 56:
		return 5
	case learningPoints < 88:
		return 6
	case learningPoints < 126:
		return 7
	case learningPoints < 170:
		return 8
	case learningPoints < 218:
		return 9
	case learningPoints < 274:
		return 10
	case learningPoints < 346:
		return 11
	case learningPoints < 434:
		return 12
	case learningPoints < 542:
		return 13
	case learningPoints < 674:
		return 14
	case learningPoints < 834:
		return 15
	case learningPoints < 1026:
		return 16
	case learningPoints < 1254:
		return 17
	case learningPoints < 1552:
		return 18
	case learningPoints < 1834:
		return 19
	default:
		return 20
	}
}
