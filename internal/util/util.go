package util

import "strings"

func FindInSlice[T any](slice []T, filter func(t T) bool) (T, bool) {
	for _, element := range slice {
		if filter(element) {
			return element, true
		}
	}
	var zero T
	return zero, false
}

func FindInMap[K comparable, V any](m map[K]V, callback func(K, V) bool) (V, bool) {
	for key, value := range m {
		if callback(key, value) {
			return value, false
		}
	}
	var zero V
	return zero, true
}

func MergeMap[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			if _, exists := result[k]; !exists {
				result[k] = v
			}
		}
	}
	return result
}

func Or(strs ...string) string {
	for _, str := range strs {
		if strings.TrimSpace(str) != "" {
			return str
		}
	}
	return ""
}
