package v1

import (
	"time"
	"fmt"
)
func parseWeekday(v string) (time.Weekday, error) {
	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday '%s'", v)
}

func generateCorrectPercent(correct int, incorrect int) string {
	correctF := float64(correct)
	incorrectF := float64(incorrect)

	percent := correctF / (correctF + incorrectF) * 100

	return fmt.Sprintf("%.2f", percent)
}