package gen

import (
	"path/filepath"
	"strings"
)

func getPackageDetails(outputDir string, pkg string) (string, string) {
	pkgSlice := strings.Split(pkg, ".")
	pkgName := pkgSlice[len(pkgSlice)-1]
	return pkgName, filepath.Join(append([]string{outputDir}, pkgSlice...)...)
}
