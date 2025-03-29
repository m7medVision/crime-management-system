package util

import (
	"github.com/m7medVision/crime-management-system/internal/model"
)

func IsClearnceLevelHigherOrEqual(level1 model.ClearanceLevel, level2 model.ClearanceLevel) bool {
	level1Value := ClearanceLevelToInt(level1)
	level2Value := ClearanceLevelToInt(level2)
	return level1Value >= level2Value
}

func ClearanceLevelToInt(level model.ClearanceLevel) int {
	switch level {
	case model.ClearanceLow:
		return 1
	case model.ClearanceMedium:
		return 2
	case model.ClearanceHigh:
		return 3
	case model.ClearanceCritical:
		return 4
	default:
		return 0
	}
}
