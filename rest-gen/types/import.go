package types

import (
	"strings"

	"github.com/tgs266/rest-gen/rest-gen/utils"
)

type Import struct {
	Path    string
	PkgName string
}

func GenerateParsedImport(specImport string, baseImportPath string) Import {
	splitPkg := strings.Split(specImport, ".")
	path := strings.ReplaceAll(specImport, ".", "/")
	fullpath := utils.ImportPathJoin(baseImportPath, path)
	pkgName := splitPkg[len(splitPkg)-1]
	return Import{
		Path:    fullpath,
		PkgName: pkgName,
	}
}
