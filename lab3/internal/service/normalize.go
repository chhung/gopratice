package service

import (
	"strings"
	"unicode"
)

func NormalizeInput(input string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}

		return r
	}, input)
}
