package types

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

type BaseType = string

const (
	TYPE_PRIMITIVE = "primitive"
	TYPE_USER      = "user"
	TYPE_WRAPPER   = "wrapper"
)

type TypeInterface interface {
	Write() *jen.Statement
	GetBaseType() BaseType
	GetWrappedType() string
}

type Type struct {
	ImportPath string
	Name       string
}

func (t Type) Write() *jen.Statement {
	if t.ImportPath == "" {
		return jen.Id(t.Name)
	}
	return jen.Qual(t.ImportPath, t.Name)
}

func (t Type) GetBaseType() BaseType {
	return TYPE_USER
}

func (t Type) GetWrappedType() string {
	panic("type is not wrapped")
}

func ParseType(typeName string, imports map[string]Import) TypeInterface {
	typeName = strings.ReplaceAll(typeName, " ", "")
	isPrimitive := IsPrimitive(typeName)
	isWrapped := IsWrapped(typeName)

	if isPrimitive {
		return GetPrimitive(typeName)
	}

	if isWrapped {
		return GetWrapper(typeName, imports)
	}
	splitType := strings.Split(typeName, ".")
	if len(splitType) == 2 {
		pkgName := splitType[0]
		if _, exists := imports[pkgName]; exists {
			imp := imports[pkgName]
			return Type{
				ImportPath: imp.Path,
				Name:       splitType[1],
			}
		}
	} else if len(splitType) > 2 {
		panic("types must not contain more than 1 period")
	}
	return Type{
		Name: typeName,
	}
}
