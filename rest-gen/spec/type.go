package spec

import (
	"github.com/dave/jennifer/jen"
	"github.com/tgs266/rest-gen/rest-gen/types"
)

func (o *Object) WriteDocs(code *jen.Statement) {
	if o.Docs != "" {
		code.Comment("o.Docs").Line()
	}
}

func (o *Object) Parse(spec *Spec) error {
	o.ObjectType = o.detectType()
	if err := o.buildInternal(spec); err != nil {
		return err
	}
	return nil
}

// Detects type of object
// will fail if more than one type is detected
func (o *Object) detectType() ObjectType {
	isStruct := len(o.Fields) != 0
	isAlias := o.Alias != nil

	if isStruct && !isAlias {
		return STRUCT
	}
	if !isStruct && isAlias {
		return ALIAS
	}
	panic("cannot declare an object as more than one type")
}

func (o *Object) buildInternal(spec *Spec) error {
	if o.ObjectType == STRUCT {
		fields, err := buildInternalFieldsFromInterface(spec, o.Fields, true)
		if err != nil {
			return err
		}
		o.ParsedFields = fields
	}
	if o.ObjectType == ALIAS {
		aliasType := types.ParseType(*o.Alias, spec.ParsedImports)
		o.ParsedAlias = aliasType
	}
	return nil
}

func buildParsedField(spec *Spec, field Field) *ParsedField {
	ty := types.ParseType(field.Type, spec.ParsedImports)
	return &ParsedField{
		Field: field,
		Type:  ty,
	}
}

func convertMapToField(spec *Spec, fieldData map[interface{}]interface{}) *ParsedField {
	docs := fieldData["docs"].(string)
	ty := fieldData["type"].(string)
	f := Field{
		Docs: docs,
		Type: ty,
	}
	return buildParsedField(spec, f)
}

func convertStringToField(spec *Spec, ty string) *ParsedField {
	f := Field{
		Type: ty,
	}
	return buildParsedField(spec, f)
}
