package service

import (
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

var specialCharMap = map[rune]rune{
	'ç': 'c',
	'Ç': 'C',
	'ñ': 'n',
	'Ñ': 'N',
	' ': '+',
}

var transformer = transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)

func NormalizeString(s string) (string, error) {
	s = replaceSpecialChars(s)
	s, err := removeAccents(s)
	return s, err
}

func replaceSpecialChars(input string) string {
	var sb strings.Builder
	for _, r := range input {
		if replacement, exists := specialCharMap[r]; exists {
			sb.WriteRune(replacement)
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func removeAccents(input string) (string, error) {
	result, _, err := transform.String(transformer, input)
	if err != nil {
		return "", err
	}
	return result, nil
}
