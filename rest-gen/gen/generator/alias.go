package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/tgs266/rest-gen/rest-gen/spec"
)

func (g *Generator) writeAlias(name string, object *spec.Object) {
	file := g.Files[FILETYPE_ALIAS]
	code := jen.Empty()
	object.WriteDocs(code)
	file.Add(code).Type().Id(name).Add(object.ParsedAlias.Write()).Line()
}
