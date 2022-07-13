package generator

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/tgs266/rest-gen/rest-gen/spec"
	"github.com/tgs266/rest-gen/rest-gen/types"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

func (g *Generator) writeError(name string, object *spec.ErrorSpec) {
	file := g.Files[FILETYPE_ERROR]
	statement := jen.Empty()

	objectCamelName := strcase.ToCamel(name)
	lowerObjectCamelName := strcase.ToLowerCamel(name)

	errStruct := &spec.Object{
		Docs:         object.Docs,
		ParsedFields: object.ParsedArgs,
		ObjectType:   types.TYPE_USER,
	}
	errStruct.WriteDocs(statement)
	statement.Add(errStruct.WriteDef(lowerObjectCamelName)).Line().Line()
	statement.Add(g.writeErrorWrapType(objectCamelName, object)).Line().Line()
	statement.Add(g.writeErrorNewFunction(objectCamelName, object)).Line().Line()
	statement.Add(g.writeIsErrorFunction(objectCamelName, object)).Line().Line()
	statement.Add(g.writeGetCause(objectCamelName, object)).Line().Line()
	statement.Add(g.writeGetCode(objectCamelName, object)).Line().Line()
	statement.Add(g.writeGetName(objectCamelName, object)).Line().Line()
	statement.Add(g.writeGetErrorId(objectCamelName, object)).Line().Line()
	statement.Add(g.writeErrorStringFunction(objectCamelName, object)).Line().Line()
	statement.Add(g.writeMarshalFcn(objectCamelName, object)).Line().Line()

	file.Add(statement)
}

func (g *Generator) writeErrorWrapType(name string, object *spec.ErrorSpec) jen.Code {
	return jen.Type().Id(name).Struct(
		jen.Id(strcase.ToLowerCamel(name)),
		jen.Id("cause").Error(),
		jen.Id("errorId").String(),
	)
}

func (g *Generator) writeErrorNewFunction(name string, object *spec.ErrorSpec) jen.Code {
	params := []jen.Code{jen.Id("err").Error()}
	fields := []jen.Code{}
	lowerName := strcase.ToLowerCamel(name)
	for _, argName := range utils.GetSortedKeys(object.ParsedArgs) {
		arg := object.ParsedArgs[argName]
		params = append(params, jen.Id(argName).Add(arg.Type.Write()))
		fields = append(fields, jen.Id(strcase.ToCamel(argName)).Op(":").Id(argName).Op(","))
	}
	return jen.Func().Id("New"+name).Params(params...).Op("*").Id(name).Block(
		jen.Id("e").Op(":=").Id(lowerName).Block(fields...),
		jen.Return(jen.Op("&").Id(name).Block(
			jen.Id("cause").Op(":").Err().Op(","),
			jen.Id("errorId").Op(":").Qual("github.com/google/uuid", "New").Call().Dot("String").Call().Op(","),
			jen.Id(lowerName).Op(":").Id("e").Op(","),
		)),
	)
}

func (g *Generator) writeIsErrorFunction(name string, object *spec.ErrorSpec) jen.Code {
	return jen.Func().Id("Is"+name).Params(jen.Err().Error()).Bool().Block(
		jen.If(jen.Err().Op("==").Nil()).Block(jen.Return(jen.False())),
		jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").Err().Assert(jen.Id(name)),
		jen.Return(jen.Id("ok")),
	)
}

func (g *Generator) writeGetCause(name string, object *spec.ErrorSpec) jen.Code {
	return jen.Func().Parens(jen.Id(strcase.ToLowerCamel(name)).Id(name)).Id("Cause").Params().Error().Block(
		jen.Return(jen.Id(strcase.ToLowerCamel(name)).Dot("cause")),
	)
}

func (g *Generator) writeGetCode(name string, object *spec.ErrorSpec) jen.Code {
	return jen.Func().Parens(jen.Id(strcase.ToLowerCamel(name)).Id(name)).Id("Code").Params().Int().Block(
		jen.Return(jen.Qual("github.com/tgs266/rest-gen/runtime/errors", object.ErrorType).Dot("Code").Call()),
	)
}

func (g *Generator) writeGetName(name string, object *spec.ErrorSpec) jen.Code {
	return jen.Func().Parens(jen.Id(strcase.ToLowerCamel(name)).Id(name)).Id("Name").Params().String().Block(
		jen.Return(jen.Lit(fmt.Sprintf("%s:%s", name, object.ErrorType))),
	)
}

func (g *Generator) writeGetErrorId(name string, object *spec.ErrorSpec) jen.Code {
	return jen.Func().Parens(jen.Id(strcase.ToLowerCamel(name)).Id(name)).Id("ErrorId").Params().String().Block(
		jen.Return(jen.Id(strcase.ToLowerCamel(name)).Dot("errorId")),
	)
}

func (g *Generator) writeErrorStringFunction(name string, object *spec.ErrorSpec) jen.Code {
	return jen.Func().Parens(jen.Id(strcase.ToLowerCamel(name)).Id(name)).Id("Error").Params().String().Block(
		jen.Return(jen.Qual("fmt", "Sprintf").Call(jen.Lit(
			fmt.Sprintf("%s:%s: %%s", name, object.ErrorType),
		), jen.Id(strcase.ToLowerCamel(name)).Dot("errorId"))),
	)
}

func (g *Generator) writeMarshalFcn(name string, object *spec.ErrorSpec) jen.Code {
	obj := strcase.ToLowerCamel(name)
	paramDict := jen.Dict{}
	for _, arg := range utils.GetSortedKeys(object.SafeArgs) {
		paramDict[jen.Lit(arg)] = jen.Id("e").Dot(obj).Dot(strcase.ToCamel(arg))
	}
	return jen.Func().Parens(jen.Id("e").Id(name)).Id("MarshalJSON").Params().Parens(jen.List(jen.Index().Byte(), jen.Error())).Block(
		jen.Id("params").Op(":=").Map(jen.String()).Interface().Values(paramDict),
		jen.Id("m").Op(":=").Qual("github.com/tgs266/rest-gen/runtime/errors", "SerializableError").Values(jen.Dict{
			jen.Id("ErrorName"):  jen.Id("e").Dot("Name").Call(),
			jen.Id("ErrorId"):    jen.Id("e").Dot("ErrorId").Call(),
			jen.Id("ErrorCode"):  jen.Id("e").Dot("Code").Call(),
			jen.Id("Parameters"): jen.Id("params"),
		}),
		jen.Return(jen.Qual("encoding/json", "Marshal").Call(jen.Id("m"))),
	)
}
