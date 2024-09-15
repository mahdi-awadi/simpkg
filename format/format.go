package format

import (
	"fmt"
	"strings"
)

// String returns string representation of interface
func String(i any) string {
	return Format("%v", i)
}

// Format returns string representation of interface
func Format(format string, i ...any) string {
	return fmt.Sprintf(format, i...)
}

// Error returns error representation of interface
func Error(format string, i ...any) error {
	return fmt.Errorf(format, i...)
}

// Replace replaces all keys in string with paired values
func Replace(str string, t map[string]any) string {
	if t != nil {
		for key, value := range t {
			prefix := "{"
			suffix := "}"
			switch value.(type) {
			case int:
			case bool:
				prefix = `"` + prefix
				suffix = suffix + `"`
				break
			}

			str = strings.ReplaceAll(str, prefix+key+suffix, Format("%v", value))
		}
	}

	return str
}
