package str

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"regexp"
	"strings"
)

// replacement map
var replacement = map[string]string{
	"ي": "ی",
	"ى": "ی",
	"ك": "ک",
}

// regex patterns
var stripRegexp = regexp.MustCompile(`\W|_+`)
var whiteSpacesRegexp = regexp.MustCompile(`\s+`)
var snakeCaseRegexp = regexp.MustCompile(`([A-Z])|(\d+)`)

// Match determine is matches strings all exists in text
func Match(matches []string, text string, normalizeValues ...bool) bool {
	if len(matches) == 0 {
		return true
	}

	normalize := len(normalizeValues) > 0 && normalizeValues[0]
	if normalize {
		text = strings.ToLower(ArabicToPersian(text))
	}

	matchedItems := 0
	for _, match := range matches {
		if normalize {
			match = strings.ToLower(ArabicToPersian(match))
		}

		// if string len is greater than 3
		if len(match) > 3 {

			// check for expression
			expr := match[:3]
			if !strings.Contains(expr, "@") {
				if strings.Contains(text, match) {
					matchedItems++
				}
				continue
			}

			// remove expression from match string
			match = match[3:]

			// ni@:	Not Include
			// rm@:	Regex Match
			// rn@:	Regex Not match
			switch expr {
			case "ni@":
				if !strings.Contains(text, match) {
					matchedItems++
				}
				continue
				break
			case "rn@":
			case "rm@":
				re := regexp.MustCompile(match)
				rm := re.MatchString(text)
				if expr == "rn@" && !rm {
					matchedItems++
				}
				if expr == "rm@" && rm {
					matchedItems++
				}
				continue
				break
			}
		} else {
			if strings.Contains(text, match) {
				matchedItems++
			}
		}
	}

	return matchedItems == len(matches)
}

// Strip remove all non-alphanumeric chars
func Strip(str string, d ...string) string {
	replaceWith := ""
	if d != nil && len(d) > 0 {
		replaceWith = d[0]
	}

	return stripRegexp.ReplaceAllString(str, replaceWith)
}

// SnakeCase convert string to snake case
func SnakeCase(str string) string {
	str = Strip(str, " ")
	str = snakeCaseRegexp.ReplaceAllString(str, " $1$2")
	str = strings.ToLower(whiteSpacesRegexp.ReplaceAllString(str, "_"))
	return strings.Trim(str, "_")
}

// ArabicToPersian normalize arabic chars to persian chars
func ArabicToPersian(text string) string {
	for key, value := range replacement {
		text = strings.Replace(text, key, value, -1)
	}

	return text
}

// UniqueSlice remove duplicate items
func UniqueSlice(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

// Md5 make md5 hash
func Md5(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

// Checksum make crc32 checksum
func Checksum(v string, upper ...bool) string {
	table := crc32.MakeTable(crc32.IEEE)
	checksum := crc32.Checksum([]byte(v), table)
	format := "%x"
	if len(upper) > 0 && upper[0] {
		format = "%X"
	}

	return fmt.Sprintf(format, checksum)
}
