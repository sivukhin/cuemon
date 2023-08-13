package lib

import (
	"strings"
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
