package tinytime

import (
	"fmt"
	"strconv"
	"strings"
)

// parseTime is a shared helper function for parsing time strings ("HH:MM" or "HH:MM:SS").
func parseTime(timeStr string) (int16, error) {
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid time format: %s", timeStr)
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours < 0 || hours > 23 {
		return 0, fmt.Errorf("invalid hours: %s", parts[0])
	}
	minutes, err := strconv.Atoi(parts[1])
	if err != nil || minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("invalid minutes: %s", parts[1])
	}
	return int16(hours*60 + minutes), nil
}

// daysBetween is a shared helper function for calculating the number of full days between two timestamps.
func daysBetween(nano1, nano2 int64) int {
	// 86400000000000 nanoseconds in a day
	const nanosInDay = 86400000000000
	return int((nano2 - nano1) / nanosInDay)
}
