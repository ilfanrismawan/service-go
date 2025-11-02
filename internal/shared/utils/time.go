package utils

import (
	"time"
)

// GetCurrentTimestamp returns current timestamp in RFC3339 format
func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}

// GetCurrentTime returns current time
func GetCurrentTime() time.Time {
	return time.Now()
}

// FormatTime formats time to string
func FormatTime(t time.Time, format string) string {
	if format == "" {
		format = time.RFC3339
	}
	return t.Format(format)
}

// ParseTime parses time string
func ParseTime(timeStr, format string) (time.Time, error) {
	if format == "" {
		format = time.RFC3339
	}
	return time.Parse(format, timeStr)
}
