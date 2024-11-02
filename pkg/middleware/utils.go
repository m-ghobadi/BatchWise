package middleware

import (
	"math"
)

func clamp(value, min, max float64) float64 {
	return math.Max(min, math.Min(value, max))
}

func getSystemLoad() float64 {
	// Simulate system load between 1 and 10
	// return rand.Float64()*9 + 1
	return 1
}

func getTrafficPatterns() float64 {
	// Simulate traffic patterns between 1 and 10
	// return rand.Float64()*9 + 1
	return 1
}

func getAveragePriority() float64 {
	// Simulate average priority between 1 and 10
	// return rand.Float64()*9 + 1
	return 1
}
