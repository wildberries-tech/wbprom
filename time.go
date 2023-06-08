package wbprom

import "time"

func MillisecondsFromStart(start time.Time) float64 {
	return float64(time.Since(start).Milliseconds())
}

func SecondsFromStart(start time.Time) float64 {
	return float64(time.Since(start).Milliseconds()) / 1000
}
