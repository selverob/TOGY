package util

import (
	"os"
	"strings"
	"time"
)

//Returns true if the file on path a was modified later than the file on path b.
//If an error is encountered, returns false and the error.
func IsNewer(a, b string) (bool, error) {
	fia, err := os.Stat(a)
	if err != nil {
		return false, err
	}

	fib, err := os.Stat(b)
	if err != nil {
		return false, err
	}

	return fia.ModTime().After(fib.ModTime()), nil
}

func NormalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func NormalizeTime(t time.Time) time.Time {
	return time.Date(1, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}

func GetFileType(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
	}

	return parts[len(parts)-1]
}
