package utils

import "sort"

func GetSortedKeys[T any](data map[string]T) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
