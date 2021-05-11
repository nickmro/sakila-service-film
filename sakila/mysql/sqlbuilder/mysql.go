package sqlbuilder

import "strings"

// Placeholders returns a string array of MySQL query param placeholders.
func Placeholders(len int) string {
	if len < 1 {
		return ""
	}

	return "? " + strings.Repeat(", ?", len-1) + ""
}
