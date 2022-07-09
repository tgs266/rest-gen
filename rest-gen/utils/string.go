package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

func StringOrDefault(str *string, def string) string {
	if str == nil {
		return def
	}
	return *str
}

func UnsafeMatchString(re, str string) bool {
	match, _ := regexp.MatchString(re, str)
	return match
}

func UrlPathJoin(paths ...string) string {
	return strings.ReplaceAll(filepath.Join(paths...), "\\", "/")
}

func CleanUrlPath(path string) string {
	if path[0] != '/' {
		path = "/" + path
	}
	return strings.TrimSuffix(path, "/")
}
