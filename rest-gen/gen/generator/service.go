package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/tgs266/rest-gen/rest-gen/spec"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

func (g *Generator) writeService(name string, service *spec.ServiceSpec) {
	file := g.Files[FILETYPE_SERVICE]
	g.writeServiceInterface(file, name, service)
}

func (g *Generator) writeServiceInterface(file *jen.File, name string, service *spec.ServiceSpec) {
	statements := []jen.Code{}
	for _, endpointName := range utils.GetSortedKeys(service.Endpoints) {
		endpoint := service.Endpoints[endpointName]
		statements = append(statements, g.writeServiceInterfaceStub(endpointName, endpoint, service.ParsedAuth != nil, service.Context))
	}
	file.Type().Id(name + "Interface").Interface(
		statements...,
	)
}

func (g *Generator) writeServiceInterfaceStub(endpointName string, endpoint *spec.Endpoint, auth bool, context bool) jen.Code {
	code := jen.Empty()
	endpointCamelName := strcase.ToCamel(endpointName)
	endpoint.WriteDocs(code).
		Id(endpointCamelName).
		Add(endpoint.WriteParams(auth, context, g.ServerGenerator.GetContextParameter())).
		Add(endpoint.WriteReturn())
	return code
}
