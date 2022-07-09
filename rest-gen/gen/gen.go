package gen

import (
	"fmt"
	"os"

	"github.com/dave/jennifer/jen"
	"github.com/tgs266/rest-gen/rest-gen/gen/generator"
	serverGenerators "github.com/tgs266/rest-gen/rest-gen/gen/generator/server-generators"
	"github.com/tgs266/rest-gen/rest-gen/spec"
	"github.com/tgs266/rest-gen/rest-gen/utils"
	"golang.org/x/mod/modfile"
)

func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GenerateFromSpec(path string, outputDir string) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		panic("cannot find go.mod file")
	}
	mod, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		panic(fmt.Errorf("could not parse go.mod: %s", err))
	}

	fileSpec := spec.Read(path)
	pkgName, pkgPath := getPackageDetails(outputDir, fileSpec.Package)
	serverGeneratorType := serverGenerators.GinServerGenerator{}

	if pathExists(pkgPath) {
		os.RemoveAll(pkgPath)
	}
	os.MkdirAll(pkgPath, os.ModePerm)

	baseImportPath := utils.ImportPathJoin(mod.Module.Mod.Path, outputDir)

	generator := &generator.Generator{
		Spec:            fileSpec,
		PkgPath:         pkgPath,
		PkgName:         pkgName,
		SrcModPath:      mod.Module.Mod.Path,
		BaseImportPath:  baseImportPath,
		ServerGenerator: serverGeneratorType,

		Files: map[string]*jen.File{},
	}

	generator.Run()

}
