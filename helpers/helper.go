package helpers

import (
	"os"
	"os/exec"
	"runtime"
	"time"
)

// IsExists check if directory exists
func IsExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// Includes returns true if the string is in the slice
func Includes[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// Diff returns the difference between two slices
func Diff[T comparable](slice1 []T, slice2 []T) []T {
	var diff []T
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := Includes(slice2, s1)
			if !found {
				diff = append(diff, s1)
			}
		}

		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

// RemoveItem removes an item from a slice
func RemoveItem[T comparable](s []T, el T) []T {
	return RemoveByIndex(s, IndexOf(s, el))
}

// RemoveByIndex removes an item by its index
func RemoveByIndex[T comparable](s []T, i int) []T {
	if i < 0 || i >= len(s) {
		return s
	}

	ret := make([]T, 0)
	ret = append(ret, s[:i]...)
	return append(ret, s[i+1:]...)
}

// IndexOf returns the index of the first instance of el in s, or -1 if el is not present in s.
func IndexOf[T comparable](s []T, el T) int {
	for i, x := range s {
		if x == el {
			return i
		}
	}
	return -1
}

// TimeInRange check if time is in range
func TimeInRange(start, end time.Time, include bool, base ...time.Time) bool {
	baseTime := time.Now()
	if len(base) > 0 && !base[0].IsZero() {
		baseTime = base[0]
	}

	if !include {
		return baseTime.After(start) && baseTime.Before(end)
	}

	return (baseTime.Equal(start) || baseTime.After(start)) && (baseTime.Equal(end) || baseTime.Before(end))
}

// OpenUrlInBrowser open url in default browser
// @see https://stackoverflow.com/a/39324149/1705598
func OpenUrlInBrowser(fileOrURL string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, fileOrURL)
	return exec.Command(cmd, args...).Start()
}
