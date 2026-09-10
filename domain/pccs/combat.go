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

// Historically have rounded down, unless less than 7 due grey area in Rules As Written
// Starting to re-think that round up is better
func calculateActions(maxSpeed int, skillFactor int) int {
	if maxSpeed == 0 {
		return 0
	}

	actionsRow := combatActionsTable[maxSpeed]
	if skillFactor < 7 {
		return combatActionsTable[maxSpeed][0]
	}

	var thresholdIndex int
	for i, t := range skillThresholds {
		if skillFactor == t {
			thresholdIndex = i
			break
		}
		if skillFactor < t {
			thresholdIndex = i - 1
			break
		}
	}

	return actionsRow[thresholdIndex]
}

// Historically have rounded down, unless less than 7 due grey area in Rules As Written
// Starting to re-think that round up is better
// Legacy app also returned 0.5 on ms 0
func calculateDamageBonus(maxSpeed int, skillFactor int) float32 {
	if maxSpeed == 0 {
		return 0
	}

	dbRow := handToHandDamageBonusTable[maxSpeed]
	var thresholdIndex int
	for i, t := range skillThresholds {
		if skillFactor == t {
			thresholdIndex = i
			break
		}
		if skillFactor < t {
			thresholdIndex = i - 1
			break
		}
	}

	return dbRow[thresholdIndex]
}
