package metrics

import "time"

func DayLabels() map[string]string {
	return map[string]string{
		"day": time.Now().UTC().Format("2006-01-02"),
	}
}