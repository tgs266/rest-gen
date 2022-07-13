package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/tgs266/rest-gen/rest-gen/spec"
	"github.com/tgs266/rest-gen/rest-gen/types"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

type ServerGeneratorInterface interface {
	GetContextParameter() jen.Code
	WriteRegisterRoutes(name string, service *spec.ServiceSpec) jen.Code
	WriteHandlerFunctionStub(handleType string, endpointName string, endpoint *spec.Endpoint) jen.Code
	WriteErrReturn(code int, errName string) jen.Code
	WriteErrReturnWithJenCode(code int, jc jen.Code) jen.Code
	WriteJsonReturn(valueName string) jen.Code
	WriteStatusCodeReturn() jen.Code
	WriteCookieReader(varName string, cookieName string) jen.Code
	WriteHeaderReader(varName string, headerName string) jen.Code

	WritePathParamReader(varName string, argName string) jen.Code
	WriteQueryParamReader(varName string, argName string) jen.Code
	WriteQueryParamArrayReader(varName string, ty types.TypeInterface) jen.Code
	WriteBodyReader(varName string, ty types.TypeInterface) jen.Code
}

type ServerGenerator struct {
	generator ServerGeneratorInterface
	auth      *spec.Auth
	context   bool
}

func (g *Generator) writeServer(
	name string,
	service *spec.ServiceSpec,
	generatorImpl ServerGeneratorInterface,
) {
	sg := &ServerGenerator{
		generator: generatorImpl,
		auth:      service.ParsedAuth,
		context:   service.Context,
	}

	file := g.Files[FILETYPE_SERVER]

	sg.writeRegisterRoutes(file, name, service)
}

func (sg *ServerGenerator) writeRegisterRoutes(
	file *jen.File,
	name string,
	service *spec.ServiceSpec,
) {
	writeServerHandlerStruct(file, name)
	file.Add(sg.generator.WriteRegisterRoutes(name, service)).Line()

	for _, endpointName := range utils.GetSortedKeys(service.Endpoints) {
		endpoint := service.Endpoints[endpointName]
		code := sg.writeServerHandlerFunction(name+"Handler", strcase.ToCamel(endpointName), endpoint)
		file.Add(code).Line()
	}
}

func writeServerHandlerStruct(file *jen.File, name string) {
	file.Type().Id(name + "Handler").Struct(
		jen.Id("Handler").Id(name + "Interface"),
	)
}

func (sg *ServerGenerator) writeServerHandlerFunction(handleType string, endpointName string, endpoint *spec.Endpoint) jen.Code {
	code := sg.generator.WriteHandlerFunctionStub(handleType, endpointName, endpoint)
	lines := []jen.Code{}
	params := []jen.Code{}
	if sg.auth != nil {
		if sg.auth.Type == spec.AUTH_COOKIE {
			lines = append(lines, sg.generator.WriteCookieReader("authToken", sg.auth.Name))
		} else if sg.auth.Type == spec.AUTH_HEADER {
			lines = append(lines, sg.generator.WriteHeaderReader("authToken", sg.auth.Name))
		}
	}
	if sg.context {
		params = append(params, jen.Id("ctx"))
	}
	for _, argName := range utils.GetSortedKeys(endpoint.ParsedArgs) {
		arg := endpoint.ParsedArgs[argName]
		switch arg.ArgLocation {
		case spec.PATH:
			lines = append(lines, writePathParamReader(sg.generator, argName, arg.Type))
		case spec.QUERY:
			lines = append(lines, writeQueryParamReader(sg.generator, argName, arg.Type))
		case spec.BODY:
			lines = append(lines, writeBodyParamReader(sg.generator, argName, arg.Type))
		}
		params = append(params, jen.Id(argName))
	}
	if sg.auth != nil {
		params = append(params, jen.Qual("github.com/tgs266/rest-gen/runtime/authentication", "Token").Call(jen.Id("authToken")))
	}

	resultName := endpointName + "Result"
	fcnCall := jen.List(endpoint.WriteReturnValue(jen.Id(resultName), jen.Id("err"))).Op(":=").Id("handler").Dot("Handler").Dot(endpointName).Call(params...)
	fcnHandle := jen.If(jen.Id("err").Op("!=").Nil()).Block(sg.generator.WriteErrReturn(500, "err"))
	var fcnReturn jen.Code
	if endpoint.HasValueReturn() {
		fcnReturn = sg.generator.WriteJsonReturn(resultName)
	} else {
		fcnReturn = sg.generator.WriteStatusCodeReturn()
	}
	lines = append(lines, fcnCall)
	lines = append(lines, fcnHandle)
	lines = append(lines, fcnReturn)
	return jen.Add(code).Block(lines...)
}

func writePathParamReader(gen ServerGeneratorInterface, argName string, ty types.TypeInterface) jen.Code {
	if ty.GetBaseType() != types.TYPE_PRIMITIVE {
		panic("path args must be primitives")
	}
	primType := ty.(types.Primitive)
	if primType.IsString {
		return gen.WritePathParamReader(argName, argName)
	}
	newArgName := argName + "Arg"
	code := gen.WritePathParamReader(newArgName, argName)
	return jen.Add(code).Line().
		Add(primType.WriteStringConverter(argName, newArgName, argName))
}

func writeQueryParamReader(gen ServerGeneratorInterface, argName string, ty types.TypeInterface) jen.Code {
	if ty.GetBaseType() != types.TYPE_PRIMITIVE && ty.GetBaseType() != types.TYPE_WRAPPER {
		panic("query args must be primitives or wrapped primitives")
	}

	if ty.GetBaseType() == types.TYPE_PRIMITIVE {
		primType := ty.(types.Primitive)
		if primType.IsString {
			return gen.WriteQueryParamReader(argName, argName)
		}
		pArgName := argName + "Arg"
		code := gen.WriteQueryParamReader(pArgName, argName)
		return jen.Add(code).Line().
			Add(primType.WriteStringConverter(argName, pArgName, argName))
	}

	code := jen.Empty()
	wrappedType := ty.(types.Wrapper)
	if !wrappedType.IsAllPrimitive() {
		panic("wrapped query args must wrap a primitive")
	}
	if wrappedType.WrapperType == types.OPTIONAL_WRAPPER {
		code.Var().Id(argName).Add(wrappedType.Write()).Line()
		code.If(jen.Add(gen.WriteQueryParamReader(argName+"Str", argName)), jen.Id(argName+"Str").Op("!=").Lit("")).Block(
			wrappedType.WriteOptionalPrimitiveStringConverter(argName, argName+"Str"),
		)
	} else if wrappedType.WrapperType == types.LIST_WRAPPER {
		code.Add(gen.WriteQueryParamArrayReader(argName, ty))
	} else {
		panic(fmt.Errorf("cannot write %s", argName))
	}
	return code
}

func writeBodyParamReader(gen ServerGeneratorInterface, argName string, ty types.TypeInterface) jen.Code {
	code := jen.Empty()
	code.Add(gen.WriteBodyReader(argName, ty))
	if ty.GetBaseType() == types.TYPE_USER {
		code.Line().If(jen.Err().Op(":=").Id(argName).Dot("Validate").Call(), jen.Err().Op("!=").Nil()).Block(
			gen.WriteErrReturnWithJenCode(400, jen.Qual("github.com/tgs266/rest-gen/runtime/errors", "NewInvalidArgumentError").Call(jen.Err())),
		)
	}
	return code
}
