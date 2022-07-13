package spec

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/tgs266/rest-gen/rest-gen/types"
	"github.com/tgs266/rest-gen/runtime/errors"
)

type AuthType string

const (
	AUTH_COOKIE AuthType = "COOKIE"
	AUTH_HEADER AuthType = "HEADER"
)

type Auth struct {
	Type AuthType
	Name string
}

type ObjectType = string

const (
	ALIAS  ObjectType = "ALIAS"
	STRUCT ObjectType = "STRUCT"
)

type Method = string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	DELETE Method = "DELETE"
)

type ArgLocation = string

const (
	PATH  ArgLocation = "PATH"
	BODY  ArgLocation = "BODY"
	QUERY ArgLocation = "QUERY"
)

type Spec struct {
	Package       string   `yaml:"package"`
	Imports       []string `yaml:"imports"`
	ParsedImports map[string]types.Import
	Types         TypeSpec                `yaml:"types"`
	Services      map[string]*ServiceSpec `yaml:"services"`
	Errors        map[string]*ErrorSpec   `yaml:"errors"`
}

type TypeSpec struct {
	Objects map[string]*Object `yaml:"objects"`
}

type Object struct {
	BSON         bool                   `yaml:"bson"`
	Builder      bool                   `yaml:"builder"`
	Docs         string                 `yaml:"docs"`
	Fields       map[string]interface{} `yaml:"fields"`
	Alias        *string                `yaml:"alias"`
	ParsedFields map[string]*ParsedField
	ParsedAlias  types.TypeInterface
	ObjectType   ObjectType
}

type Field struct {
	Validation string `yaml:"validation"`
	Type       string `yaml:"type"`
	Docs       string `yaml:"docs"`
}

type ParsedField struct {
	Field       Field
	Type        types.TypeInterface
	ArgLocation ArgLocation
}

type ServiceSpec struct {
	Package    string               `yaml:"package"`
	BasePath   string               `yaml:"base-path"`
	Endpoints  map[string]*Endpoint `yaml:"endpoints"`
	Auth       string               `yaml:"auth"`
	ParsedAuth *Auth
}

type HTTP struct {
	Method Method
	Path   string
}

func (h HTTP) String() string {
	return h.Method + " " + h.Path
}

type Endpoint struct {
	Docs          string       `yaml:"docs"`
	HTTP          string       `yaml:"http"`
	Args          EndpointArgs `yaml:"args"`
	Returns       string       `yaml:"returns"`
	ParsedArgs    map[string]*ParsedField
	ParsedReturns *ParsedField
	ParsedHTTP    HTTP
}

type EndpointArgs struct {
	Path  map[string]interface{} `yaml:"path"`
	Query map[string]interface{} `yaml:"query"`
	Body  interface{}            `yaml:"body"`
}

type ErrorSpec struct {
	ErrorCode       string `yaml:"errorCode"`
	ParsedErrorCode errors.ErrorCode
	Docs            string                 `yaml:"docs"`
	SafeArgs        map[string]interface{} `yaml:"safeArgs"`
	UnsafeArgs      map[string]interface{} `yaml:"unsafeArgs"`
	ParsedArgs      map[string]*ParsedField
}

func (s *Spec) Parse(baseImportPath string) {
	parsedImports := map[string]types.Import{}
	for _, importSpec := range s.Imports {
		pi := types.GenerateParsedImport(importSpec, baseImportPath)
		if _, exists := parsedImports[pi.PkgName]; exists {
			panic("spec contains multiple imports with the same package name")
		}
		parsedImports[pi.PkgName] = pi
	}
	s.ParsedImports = parsedImports
}

func buildInternalFieldsFromInterface(
	spec *Spec,
	providedFields map[string]interface{},
	capitalize bool,
) (map[string]*ParsedField, error) {
	fields := map[string]*ParsedField{}
	for fieldName, fieldData := range providedFields {
		if fieldData == nil {
			continue
		}
		fieldString, isString := fieldData.(string)
		fieldMap, isMap := fieldData.(map[interface{}]interface{})
		useFieldName := fieldName
		if capitalize {
			useFieldName = strcase.ToCamel(fieldName)
		}
		if isString {
			fields[useFieldName] = convertStringToField(spec, fieldString)
		} else if isMap {
			fields[useFieldName] = convertMapToField(spec, fieldMap)
		} else {
			return nil, fmt.Errorf("field \"%s\" does not satisfy the constraints to be a field", fieldName)
		}
	}
	return fields, nil
}
