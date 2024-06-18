package lib

import (
	"regexp"
	"strings"
	"unicode"
)

func Unique[T comparable](a []T) []T {
	used := make(map[T]struct{}, 0)
	elements := make([]T, 0)
	for _, element := range a {
		if _, ok := used[element]; !ok {
			used[element] = struct{}{}
			elements = append(elements, element)
		}
	}
	return elements
}

func AsAny[T any](array []T) []any {
	a := make([]any, 0, len(array))
	for _, x := range array {
		a = append(a, x)
	}
	return a
}

func PackageName(module string) string {
	modulePath := strings.Split(module, "/")
	return modulePath[len(modulePath)-1]
}

func CapitalizeName(s string) string {
	var matchIgnoredSymbols = regexp.MustCompile("[^a-zA-Z0-9]")
	var matchUnderscores = regexp.MustCompile("_+")

	s = matchIgnoredSymbols.ReplaceAllString(s, "_")
	s = matchUnderscores.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")

	var result strings.Builder
	makeCapital := false
	for _, c := range s {
		if c == '_' {
			makeCapital = true
			continue
		}
		if unicode.IsTitle(c) {
			makeCapital = true
		}
		if makeCapital {
			result.WriteRune(unicode.ToTitle(c))
		} else {
			result.WriteRune(unicode.ToLower(c))
		}
		makeCapital = false
	}
	return result.String()
}
