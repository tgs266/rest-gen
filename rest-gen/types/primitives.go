package types

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

type Primitive struct {
	Name       string
	ImportPath string
	IsString   bool
}

var PRIMITIVES = map[string]TypeInterface{
	"string":  makePrimitive("string", true),
	"byte":    makePrimitive("byte", false),
	"binary":  makePrimitive("[]byte", false),
	"short":   makePrimitive("int16", false),
	"int":     makePrimitive("int", false),
	"long":    makePrimitive("int64", false),
	"float":   makePrimitive("float32", false),
	"double":  makePrimitive("float64", false),
	"boolean": makePrimitive("bool", false),

	"datetime": makeImportPrimitive("time", "Time", false),
}

var ALLOWED_PRIMITIVE_STRING = fmt.Sprintf("[%s]", getPrimitiveString())

func makePrimitive(name string, conv bool) TypeInterface {
	return Primitive{
		Name:     name,
		IsString: conv,
	}
}

func makeImportPrimitive(importPath string, name string, conv bool) TypeInterface {
	return Primitive{
		Name:       name,
		ImportPath: importPath,
		IsString:   conv,
	}
}

func getPrimitiveString() string {
	str := []string{}
	for p := range PRIMITIVES {
		str = append(str, p)
	}
	return strings.Join(str, ", ")
}

func writeIntConversion(arg string, inputValueName string, outputValueName string) jen.Code {
	return jen.List(jen.Id(outputValueName), jen.Id("err")).
		Op(":=").
		Qual("strconv", "Atoi").
		Call(jen.Id(inputValueName)).
		Line().
		Add(utils.WriteErrorCheck("err", fmt.Sprintf("failed to parse \"%s\" as %s", arg, "int")))
}

func writeTypeConversion(
	arg string,
	inType string,
	inputValueName string,
	outputValueName string,
) jen.Code {
	to := "To" + strcase.ToCamel(inType) + "E"
	return jen.List(jen.Id(outputValueName), jen.Id("err")).
		Op(":=").
		Qual("github.com/spf13/cast", to).
		Call(jen.Id(inputValueName)).
		Line().
		Add(utils.WriteErrorCheck("err", fmt.Sprintf("failed to parse \"%s\" as %s", arg, inType)))
}

func (t Primitive) WriteStringConverter(arg string, inputName string, outputName string) jen.Code {
	if t.Name == "string" {
		return nil
	}
	if t.Name == "byte" {
		panic("cannot convert string to byte")
	}
	if t.Name == "time" {
		panic("cannot convert string to time (for now)")
	}
	return writeTypeConversion(arg, t.Name, inputName, outputName)
}

func (t Primitive) Write() *jen.Statement {
	if t.ImportPath == "" {
		return jen.Id(t.Name)
	}
	return jen.Qual(t.ImportPath, t.Name)
}

func (t Primitive) GetBaseType() BaseType {
	return TYPE_PRIMITIVE
}

func (t Primitive) GetWrappedType() string {
	panic("type is not wrapped")
}

func GetPrimitive(prim string) TypeInterface {
	if val, exists := PRIMITIVES[prim]; exists {
		return val
	}
	panic(
		fmt.Errorf(
			"primitive type %s is unknown. primitives must be one of %s",
			prim,
			ALLOWED_PRIMITIVE_STRING,
		),
	)
}

func IsPrimitive(str string) bool {
	match, _ := regexp.MatchString("^[a-z]+[a-zA-Z0-9]*$", str)
	return match
}
