package utils

/*
 * Gonys - A Notification Service for SMS
 *
 * String Utilities
 *
 * @author A. A. Sumitro <hello@aasumitro.id>
 * https://aasumitro.id
 */

import (
	"regexp"
	"strings"
)

// Transpose the string
func Transpose(data string) string {
	output := strings.Replace(data, "\r\n", "\\r\\n", -1)
	return strings.Replace(output, "\r", "\\r", -1)
}

func addWordBoundariesToNumbers(s string) string {
	var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
	var numberReplacement = []byte(`$1 $2 $3`)

	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)

	return string(b)
}

// ToSnake converts a string to snake_case
func ToSnake(str string, delimiter uint8) string {
	str = addWordBoundariesToNumbers(str)
	str = strings.Trim(str, " ")
	out := ""

	for i, v := range str {
		nextCaseIsChanged := false

		if i + 1 < len(str) {
			next := str[i+1]
			vIsCap := v >= 'A' && v <= 'Z'
			vIsLow := v >= 'a' && v <= 'z'
			nextIsCap := next >= 'A' && next <= 'Z'
			nextIsLow := next >= 'a' && next <= 'z'
			if (vIsCap && nextIsLow) || (vIsLow && nextIsCap) {
				nextCaseIsChanged = true
			}
		}

		if i > 0 && out[len(out) - 1] != delimiter && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if v >= 'A' && v <= 'Z' {
				out += string(delimiter) + string(v)
			} else if v >= 'a' && v <= 'z' {
				out += string(v) + string(delimiter)
			}
		} else if v == ' ' || v == '_' || v == '-' {
			out += string(delimiter)
		} else {
			out = out + string(v)
		}
	}

	return strings.ToLower(out)
}