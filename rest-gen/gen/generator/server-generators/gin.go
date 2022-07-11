package servergenerator

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/tgs266/rest-gen/rest-gen/spec"
	"github.com/tgs266/rest-gen/rest-gen/types"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

var ginImport = "github.com/gin-gonic/gin"

type GinServerGenerator struct {
}

func (gsg GinServerGenerator) WriteRegisterRoutes(name string, service *spec.ServiceSpec) jen.Code {
	statements := []jen.Code{}
	for _, endpointName := range utils.GetSortedKeys(service.Endpoints) {
		endpoint := service.Endpoints[endpointName]
		statements = append(statements, gsg.writeRegisterRoutesRoute(endpointName, endpoint, "handler"))
	}
	return jen.Func().Id("Register"+name+"Routes").Params(
		jen.Id("router").Op("*").Qual(ginImport, "Engine"),
		jen.Id("handler").Id(name+"Handler"),
	).Block(
		statements...,
	)
}

// name + "Handler"
func (gsg GinServerGenerator) writeRegisterRoutesRoute(
	name string,
	endpoint *spec.Endpoint,
	serviceName string,
) jen.Code {
	fcnCall := ""
	switch endpoint.ParsedHTTP.Method {
	case spec.GET:
		fcnCall = "GET"
	case spec.PUT:
		fcnCall = "PUT"
	case spec.POST:
		fcnCall = "POST"
	case spec.DELETE:
		fcnCall = "DELETE"
	}
	return jen.Id("router").Dot(fcnCall).Call(jen.Lit(endpoint.ParsedHTTP.Path), jen.Id(serviceName).Dot("Handle"+strcase.ToCamel(name)))
}

func (gsg GinServerGenerator) WriteHandlerFunctionStub(
	handleType string,
	endpointName string,
	endpoint *spec.Endpoint,
) jen.Code {
	return jen.Func().
		Parens(jen.Id("handler").Id(handleType)).
		Id("Handle" + endpointName).
		Params(jen.Id("ctx").Op("*").Qual(ginImport, "Context"))
}

func (gsg GinServerGenerator) WritePathParamReader(varName, argName string) jen.Code {
	return jen.Id(varName).Op(":=").Id("ctx").Dot("Param").Call(jen.Lit(argName))
}

func (gsg GinServerGenerator) WriteQueryParamReader(varName, argName string) jen.Code {
	return jen.Id(varName).Op(":=").Id("ctx").Dot("Query").Call(jen.Lit(argName))
}

func (gsg GinServerGenerator) WriteQueryParamArrayReader(varName string, ty types.TypeInterface) jen.Code {
	return jen.Var().Id(varName).Add(ty.Write()).Line().
		If(jen.Id("err").Op(":=").Id("ctx").Dot("ShouldBindQuery").Call(jen.Op("&").Id(varName)), jen.Id("err").Op("!=").Nil()).Block(
		jen.Id(varName).Op("=").Make(ty.Write(), jen.Lit(0)),
	)
}

// allow optional
func (gsg GinServerGenerator) WriteBodyReader(varName string, ty types.TypeInterface) jen.Code {
	return jen.Var().Id(varName).Add(ty.Write()).Line().
		If(jen.Id("err").Op(":=").Id("ctx").Dot("ShouldBindJSON").Call(jen.Op("&").Id(varName)), jen.Id("err").Op("!=").Nil()).Block(
		jen.Panic(jen.Lit("needs filler")),
	)
}

func (gsg GinServerGenerator) WriteErrReturn(errName string) jen.Code {
	return gsg.WriteErrReturnWithCode(500, errName)
}

func (gsg GinServerGenerator) WriteErrReturnWithCode(code int, errName string) jen.Code {
	return jen.Id("ctx").Dot("AbortWithError").Call(jen.Lit(code), jen.Id(errName)).Line().Return()
}

func (gsg GinServerGenerator) WriteJsonReturn(value string) jen.Code {
	return jen.Id("ctx").Dot("JSON").Call(jen.Qual("net/http", "StatusOK"), jen.Id(value)).Line().Return()
}

func (gsg GinServerGenerator) WriteStatusCodeReturn() jen.Code {
	return jen.Id("ctx").Dot("Status").Call(jen.Qual("net/http", "StatusOK")).Line().Return()
}
