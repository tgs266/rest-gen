package types

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

type OptionalWrapper struct {
	Type TypeInterface
}

type ListWrapper struct {
	Type TypeInterface
}

type MapWrapper struct {
	Type1 TypeInterface
	Type2 TypeInterface
}

type WrapperType = string

const (
	LIST_WRAPPER     WrapperType = "LIST_WRAPPER"
	MAP_WRAPPER      WrapperType = "MAP_WRAPPER"
	OPTIONAL_WRAPPER WrapperType = "OPTIONAL_WRAPPER"
)

type Wrapper struct {
	Types       []TypeInterface
	WrapperType WrapperType
}

func (w Wrapper) IsAllPrimitive() bool {
	for _, t := range w.Types {
		if t.GetBaseType() != TYPE_PRIMITIVE {
			return false
		}
	}
	return true
}

func (w Wrapper) WriteOptionalPrimitiveStringConverter(varName string, stringName string) jen.Code {
	primType := w.Types[0].(Primitive)
	return jen.Add(primType.WriteStringConverter(varName, stringName, varName+"Val")).Line().
		Id(varName).Op("=").Op("&").Id(varName + "Val")
}

var LIST_WRAPPER_REGEX = "^list<[a-zA-Z0-9<>,\\.]*>$"
var OPTIONAL_WRAPPER_REGEX = "^optional<[a-zA-Z0-9<>,\\.]*>$"
var MAP_WRAPPER_REGEX = "^map<[a-zA-Z0-9<>,\\.]*\\,[a-zA-Z0-9<>,\\.]*>$"

func IsWrapped(ty string) bool {
	match, _ := regexp.MatchString("^[a-z]+[a-zA-Z0-9]*<[a-zA-Z0-9,\\.<>]*>$", ty)
	return match
}

func GetWrapper(typeStr string, imports map[string]Import) TypeInterface {
	if utils.UnsafeMatchString(LIST_WRAPPER_REGEX, typeStr) {
		innerType := getSingleInnerType("list", typeStr, imports)
		return Wrapper{
			Types:       []TypeInterface{innerType},
			WrapperType: LIST_WRAPPER,
		}
	}
	if utils.UnsafeMatchString(OPTIONAL_WRAPPER_REGEX, typeStr) {
		innerType := getSingleInnerType("optional", typeStr, imports)
		return Wrapper{
			Types:       []TypeInterface{innerType},
			WrapperType: OPTIONAL_WRAPPER,
		}
	}
	if utils.UnsafeMatchString(MAP_WRAPPER_REGEX, typeStr) {
		t1, t2 := getDoubleInnerType("map", typeStr, imports)
		return Wrapper{
			Types:       []TypeInterface{t1, t2},
			WrapperType: MAP_WRAPPER,
		}
	}
	panic(fmt.Errorf("could not parse wrapped type %s", typeStr))
}

func getSingleInnerType(name, typeStr string, imports map[string]Import) TypeInterface {
	newTypeStr := strings.Replace(typeStr, name, "", 1)
	return ParseType(strings.TrimPrefix(strings.TrimSuffix(newTypeStr, ">"), "<"), imports)
}

func getDoubleInnerType(
	name, typeStr string,
	imports map[string]Import,
) (TypeInterface, TypeInterface) {
	newTypeStr := strings.Replace(typeStr, name, "", 1)
	newTypeStr = strings.TrimPrefix(strings.TrimSuffix(newTypeStr, ">"), "<")
	count := 0
	step := 0
	for _, c := range newTypeStr {
		if c == '<' {
			count += 1
		} else if c == '>' {
			count -= 1
		}
		if count == 0 && c == ',' {
			break
		}
		step += 1
	}

	firstType := newTypeStr[:step]
	secondType := newTypeStr[step+1:] // +1 to skip comma

	return ParseType(firstType, imports), ParseType(secondType, imports)
}

func (w Wrapper) Write() *jen.Statement {
	switch w.WrapperType {
	case MAP_WRAPPER:
		return jen.Map(w.Types[0].Write()).Add(w.Types[1].Write())
	case LIST_WRAPPER:
		return jen.Index().Add(w.Types[0].Write())
	case OPTIONAL_WRAPPER:
		return jen.Op("*").Add(w.Types[0].Write())
	}
	return jen.Empty()
}

func (w Wrapper) GetBaseType() BaseType {
	return TYPE_WRAPPER
}

func (w Wrapper) GetWrappedType() WrapperType {
	return w.WrapperType
}
