package generator

import (
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

func (g *Generator) writeImports(file *jen.File) {
	imports := g.Spec.Imports
	for _, importData := range imports {
		writeImportString(importData, file, g.BaseImportPath)
	}
}

func writeImportString(str string, file *jen.File, baseImportPath string) {
	splitPkg := strings.Split(str, ".")
	path := strings.ReplaceAll(str, ".", "/")
	fullpath := utils.ImportPathJoin(baseImportPath, path)
	pkgName := splitPkg[len(splitPkg)-1]
	file.ImportName(fullpath, pkgName)
}
