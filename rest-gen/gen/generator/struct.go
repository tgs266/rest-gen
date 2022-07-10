package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/tgs266/rest-gen/rest-gen/spec"
	"github.com/tgs266/rest-gen/rest-gen/types"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

func (g *Generator) writeStruct(name string, object *spec.Object) {
	file := g.Files[FILETYPE_STRUCT]
	statement := jen.Empty()

	object.WriteDocs(statement)
	structFields := []jen.Code{}
	for _, fieldName := range utils.GetSortedKeys(object.ParsedFields) {
		fieldData := object.ParsedFields[fieldName]
		code := jen.Empty()

		if fieldData.Field.Docs != "" {
			code.Comment(fieldData.Field.Docs).Line()
		}
		code.Id(strcase.ToCamel(fieldName)).Add(fieldData.Type.Write())
		code.Tag(map[string]string{
			"json": strcase.ToLowerCamel(fieldName),
			"yaml": strcase.ToLowerCamel(fieldName),
		})
		structFields = append(structFields, code)
	}

	file.Add(statement).Type().Id(name).Struct(structFields...).Line()

	if object.Builder {
		file.Add(writeStructBuilderType(name, object)).Line()
		file.Add(writeNewStructBuilderFunction(name, object)).Line()
		file.Add(writeStrctBuilderFieldFunctions(name, object)).Line()
		file.Add(writeStructBuilderBuildFunction(name, object)).Line()
	}

	// writeStructMarshal(file, "encoding/json", "MarshalJSON", name, object).Line()
	// writeStructUnmarshal(file, "encoding/json", "UnmarshalJSON", name, object).Line()
	// writeStructMarshal(file, "gopkg.in/yaml.v3", "MarshalYAML", name, object).Line()
	// writeStructUnmarshal(file, "gopkg.in/yaml.v3", "UnmarshalYAML", name, object).Line()
}

func writeStructMarshal(
	file *jen.File,
	pkgName, fcnName, name string,
	object *spec.Object,
) *jen.File {
	lowerName := strcase.ToLowerCamel(name)
	file.Id("func").
		Parens(jen.Id(lowerName).Id(name)).
		Id(fcnName).
		Params().
		Parens(jen.List(jen.Index().Byte(), jen.Error())).
		Block(
			writeStructMarshalInner(pkgName, lowerName, object),
		)
	return file
}

func writeStructMarshalInner(pkgName string, lowerName string, object *spec.Object) *jen.Statement {
	code := jen.Empty()
	for fieldName, field := range object.ParsedFields {
		if field.Type.GetBaseType() == types.TYPE_WRAPPER {
			maker := jen.Make(field.Type.Write(), jen.Lit(0))
			if v, ok := field.Type.(types.Wrapper); ok {
				if v.WrapperType == types.OPTIONAL_WRAPPER {
					maker = jen.New(v.Types[0].Write())
				}
			}
			code.If(jen.Id(lowerName).Dot(fieldName).Op("==").Nil()).Block(
				jen.Id(lowerName).Dot(fieldName).Op("=").Add(maker),
			).Line()
		}
	}
	code.Return(jen.Qual(pkgName, "Marshal").Call(jen.Id(lowerName)))
	return code
}

func writeStructUnmarshal(
	file *jen.File,
	pkgName, fcnName, name string,
	object *spec.Object,
) *jen.File {
	lowerName := strcase.ToLowerCamel(name)
	file.Id("func").
		Parens(jen.Id(lowerName).Op("*").Id(name)).
		Id(fcnName).
		Params(jen.Id("bytes").Index().Byte()).
		Error().
		Block(
			writeStructUnmarshalInner(pkgName, lowerName, object),
		)
	return file
}

func writeStructUnmarshalInner(
	pkgName string,
	lowerName string,
	object *spec.Object,
) *jen.Statement {
	return jen.Return(
		jen.Qual(pkgName, "Unmarshal").Call(jen.Id("bytes"), jen.Op("&").Id(lowerName)),
	)
}

func writeStructBuilderType(
	name string,
	object *spec.Object,
) jen.Code {
	structFields := []jen.Code{}
	for _, fieldName := range utils.GetSortedKeys(object.ParsedFields) {
		fieldData := object.ParsedFields[fieldName]
		code := jen.Empty()
		code.Id(strcase.ToCamel(fieldName)).Add(fieldData.Type.Write())
		structFields = append(structFields, code)
	}
	return jen.Type().Id(strcase.ToLowerCamel(name) + "Builder").Struct(structFields...)
}

func writeNewStructBuilderFunction(
	name string,
	object *spec.Object,
) jen.Code {
	return jen.Func().Id("New" + name + "Builder").Params().Op("*").Id(strcase.ToLowerCamel(name) + "Builder").Block(
		jen.Return().Op("&").Id(strcase.ToLowerCamel(name) + "Builder").Block(),
	)
}

func writeStrctBuilderFieldFunctions(
	name string,
	object *spec.Object,
) jen.Code {
	fcns := []jen.Code{}
	for _, fieldName := range utils.GetSortedKeys(object.ParsedFields) {
		capFName := strcase.ToCamel(fieldName)
		fieldData := object.ParsedFields[fieldName]
		code := jen.Func().Parens(jen.Id("builder").Id(strcase.ToLowerCamel(name)+"Builder")).Id("Set"+strcase.ToCamel(fieldName)).Params(jen.Id(fieldName).Add(fieldData.Type.Write())).Id(strcase.ToLowerCamel(name)+"Builder").Block(
			jen.Id("builder").Dot(capFName).Op("=").Id(fieldName),
			jen.Return().Id("builder"),
		).Line()
		fcns = append(fcns, code)
	}
	return jen.Empty().Add(fcns...)
}

func writeStructBuilderBuildFunction(
	name string,
	object *spec.Object,
) jen.Code {
	setters := []jen.Code{}
	for _, fieldName := range utils.GetSortedKeys(object.ParsedFields) {
		code := jen.Id(strcase.ToCamel(fieldName)).Op(":").Id("builder").Dot(strcase.ToCamel(fieldName)).Op(",")
		setters = append(setters, code)
	}
	code := jen.Func().Parens(jen.Id("builder").Id(strcase.ToLowerCamel(name) + "Builder")).Id("Build").Params().Id(strcase.ToCamel(name)).Block(
		jen.Return().Id(strcase.ToCamel(name)).Block(
			setters...,
		),
	).Line()
	return code
}
