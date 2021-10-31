package metrics

import "time"

// DayLabels returns a label map with the current date on the format '2006-01-02'
func DayLabels() map[string]string {
	return map[string]string{
		"day": time.Now().UTC().Format("2006-01-02"),
	}
}