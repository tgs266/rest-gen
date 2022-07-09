package utils

import (
	"path/filepath"
	"strings"
)

// creates proper import path no matter the os
func ImportPathJoin(paths ...string) string {
	return strings.ReplaceAll(filepath.Join(paths...), "\\", "/")
}
