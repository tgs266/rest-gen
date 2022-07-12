package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/tgs266/rest-gen/rest-gen/spec"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

func (g *Generator) writeService(name string, service *spec.ServiceSpec) {
	file := g.Files[FILETYPE_SERVICE]
	writeServiceInterface(file, name, service)
}

func writeServiceInterface(file *jen.File, name string, service *spec.ServiceSpec) {
	statements := []jen.Code{}
	for _, endpointName := range utils.GetSortedKeys(service.Endpoints) {
		endpoint := service.Endpoints[endpointName]
		statements = append(statements, writeServiceInterfaceStub(endpointName, endpoint, service.ParsedAuth != nil))
	}
	file.Type().Id(name + "Interface").Interface(
		statements...,
	)
}

func writeServiceInterfaceStub(endpointName string, endpoint *spec.Endpoint, auth bool) jen.Code {
	code := jen.Empty()
	endpointCamelName := strcase.ToCamel(endpointName)
	endpoint.WriteDocs(code).
		Id(endpointCamelName).
		Add(endpoint.WriteParams(auth)).
		Add(endpoint.WriteReturn())
	return code
}
