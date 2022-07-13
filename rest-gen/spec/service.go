package spec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/tgs266/rest-gen/rest-gen/utils"
)

func (e *Endpoint) WriteParams(auth bool) *jen.Statement {
	params := []jen.Code{}
	for _, argName := range utils.GetSortedKeys(e.ParsedArgs) {
		arg := e.ParsedArgs[argName]
		params = append(params, jen.Id(argName).Add(arg.Type.Write()))
	}
	if auth {
		params = append(params, jen.Id("authToken").Qual("github.com/tgs266/rest-gen/runtime/token", "Token"))
	}
	return jen.Params(params...)
}

func (e *Endpoint) HasValueReturn() bool {
	return e.Returns != ""
}

func (e *Endpoint) WriteReturn() *jen.Statement {
	if !e.HasValueReturn() {
		return jen.Error()
	}
	return jen.Parens(jen.List(e.ParsedReturns.Type.Write(), jen.Error()))
}

func (e *Endpoint) WriteReturnValue(returnValue, returnErr jen.Code) jen.Code {
	if !e.HasValueReturn() {
		return returnErr
	}
	return jen.List(returnValue, returnErr)
}

func (e *Endpoint) WriteActualReturnValue(returnValue, returnErr jen.Code) jen.Code {
	if !e.HasValueReturn() {
		return jen.Return().Add(returnErr)
	}
	return jen.Return().List(returnValue, returnErr)
}

func (o *Endpoint) WriteDocs(code *jen.Statement) *jen.Statement {
	if o.Docs != "" {
		code.Comment(o.Docs).Line()
	}
	return code
}

func (s *ServiceSpec) Parse(spec *Spec) error {
	if s.Auth != "" {
		if strings.HasPrefix(s.Auth, "cookie") {
			name := strings.Split(s.Auth, ":")[1]
			s.ParsedAuth = &Auth{
				Name: name,
				Type: AUTH_COOKIE,
			}
		} else if strings.HasPrefix(s.Auth, "header") {
			s.ParsedAuth = &Auth{
				Name: "Authorization",
				Type: AUTH_HEADER,
			}
		} else {
			return fmt.Errorf("auth type %s not supported", s.Auth)
		}
	}
	httpValidate := map[string]string{}
	for eName, endpoint := range s.Endpoints {
		splitHttp := strings.Split(endpoint.HTTP, " ")
		if len(splitHttp) != 2 {
			return fmt.Errorf("HTTP must contain a method and path, seperated by a space")
		}
		splitHttp[1] = utils.CleanUrlPath(splitHttp[1])
		properMethod := GET
		switch strings.ToLower(splitHttp[0]) {
		case "get":
			properMethod = GET
		case "post":
			properMethod = POST
		case "put":
			properMethod = PUT
		case "delete":
			properMethod = DELETE
		}
		parsedHttp := HTTP{
			Method: properMethod,
			Path:   utils.UrlPathJoin(s.BasePath, splitHttp[1]),
		}
		if _, exists := httpValidate[parsedHttp.String()]; exists {
			panic("cannot have multiple endpoints with the same method and path")
		}
		s.Endpoints[eName].ParsedHTTP = parsedHttp
		if err := endpoint.Parse(spec); err != nil {
			return fmt.Errorf("%s: %s", eName, err)
		}

	}
	return nil
}

func (e *Endpoint) Parse(spec *Spec) error {
	if err := e.buildInternal(spec); err != nil {
		return err
	}
	if err := e.validateParams(spec); err != nil {
		return err
	}
	return nil
}

func mergeInto(m1, mOther map[string]*ParsedField) (map[string]*ParsedField, error) {
	for k, v := range mOther {
		if _, exists := m1[k]; exists {
			return nil, fmt.Errorf("used the field name \"%s\" twice", k)
		}
		m1[k] = v
	}
	return m1, nil
}

func wrapWithLocation(args map[string]*ParsedField, location ArgLocation) map[string]*ParsedField {
	for k, v := range args {
		v.ArgLocation = location
		args[k] = v
	}
	return args
}

func (e *Endpoint) buildInternal(spec *Spec) error {
	args := map[string]*ParsedField{}
	if e.Args.Body != nil {
		bodyArgs, err := buildInternalFieldsFromInterface(spec, map[string]interface{}{"body": e.Args.Body}, false)
		if err != nil {
			return err
		}
		args = wrapWithLocation(bodyArgs, BODY)
	}
	pathArgs, err := buildInternalFieldsFromInterface(spec, e.Args.Path, false)
	if err != nil {
		return err
	}
	queryArgs, err := buildInternalFieldsFromInterface(spec, e.Args.Query, false)
	if err != nil {
		return err
	}
	pathArgs = wrapWithLocation(pathArgs, PATH)
	queryArgs = wrapWithLocation(queryArgs, QUERY)

	args, err = mergeInto(args, pathArgs)
	if err != nil {
		return err
	}
	args, err = mergeInto(args, queryArgs)
	if err != nil {
		return err
	}
	e.ParsedArgs = args
	if e.Returns != "" {
		e.ParsedReturns = convertStringToField(spec, e.Returns)
	}
	return nil
}

func (e *Endpoint) validateParams(spec *Spec) error {
	paramRegex := regexp.MustCompile(`\{.*?\}`)
	params := paramRegex.FindAllStringSubmatch(e.HTTP, -1)
	if len(e.Args.Path) != len(params) {
		return fmt.Errorf("provided path args must match those defined in the url path")

	}
	for _, p := range params {
		pString := strings.TrimPrefix(strings.TrimSuffix(p[0], "}"), "{")
		if _, exists := e.Args.Path[pString]; !exists {
			return fmt.Errorf("path parameter {%s} must be defined", pString)
		}
	}
	return nil
}

// func buildInternalFieldsFromStruct(spec *Spec, providedFields map[string]interface{}) (map[string]ParsedField, error) {
// 	fields := map[string]ParsedField{}
// 	for fieldName, fieldData := range providedFields {
// 		fieldString, isString := fieldData.(string)
// 		fieldMap, isMap := fieldData.(map[interface{}]interface{})
// 		if isString {
// 			fields[strcase.ToCamel(fieldName)] = convertStringToField(spec, fieldString)
// 		} else if isMap {
// 			fields[strcase.ToCamel(fieldName)] = convertMapToField(spec, fieldMap)
// 		} else {
// 			return nil, fmt.Errorf("field \"%s\" does not satisfy the constraints of a struct", fieldName)
// 		}
// 	}
// 	return fields, nil
// }

// func buildParsedField(spec *Spec, field Field) ParsedField {
// 	ty := types.ParseType(field.Type, spec.Types.ParsedImports)
// 	return ParsedField{
// 		Field: field,
// 		Type:  ty,
// 	}
// }

// func convertMapToField(spec *Spec, fieldData map[interface{}]interface{}) ParsedField {
// 	docs := fieldData["docs"].(string)
// 	ty := fieldData["type"].(string)
// 	f := Field{
// 		Docs: docs,
// 		Type: ty,
// 	}
// 	return buildParsedField(spec, f)
// }

// func convertStringToField(spec *Spec, ty string) ParsedField {
// 	f := Field{
// 		Type: ty,
// 	}
// 	return buildParsedField(spec, f)
// }
