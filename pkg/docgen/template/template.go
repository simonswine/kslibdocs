package template

import (
	"strings"
)

// ClassFor converts a slice of strings to a class
func ClassFor(items []string) string {
	classes := strings.Join(items, " ")
	return classes
}
