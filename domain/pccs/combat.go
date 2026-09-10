package pccs

func calculateBaseSpeed(str int, enc float32) float32 {
	strRow := speedTable[str]

	low := 0
	high := len(encumbranceThresholds)

	for low < high {
		mid := (low + high) / 2
		if encumbranceThresholds[mid] < enc {
			low = mid + 1
		} else {
			high = mid
		}
	}

	if low == len(encumbranceThresholds) {
		return 0
	}

	return strRow[low]
}

func calculateMaxSpeed(agi int, baseSpeed float32) int {
	if baseSpeed == 0 {
		return 0
	}

	var thresholdIndex int
	for i, t := range baseSpeedThresholds {
		if t == baseSpeed {
			thresholdIndex = i
			break
		}
	}

	agiRow := maxSpeedTable[agi]
	return agiRow[thresholdIndex]
}

func calculateActions(maxSpeed int, skillFactor int) int {
	if maxSpeed == 0 {
		return 0
	}

	return 1
}
