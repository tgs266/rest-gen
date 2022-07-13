package spec

import (
	"fmt"

	"github.com/tgs266/rest-gen/runtime/errors"
)

func ErrorCodeFromString(errorCodeName string, spec *Spec) (errors.ErrorType, error) {
	if v, exists := errors.KnownErrorCode[errorCodeName]; exists {
		return v, nil
	}
	return errors.ErrorType{}, fmt.Errorf("error code %s unknown", errorCodeName)
}

func (s *ErrorSpec) IsSafe(argName string) bool {
	_, exists := s.SafeArgs[argName]
	return exists
}

func (s *ErrorSpec) WriteErrorStruct(argName string) bool {
	_, exists := s.SafeArgs[argName]
	return exists
}

func (s *ErrorSpec) Parse(spec *Spec) error {
	ec, err := ErrorCodeFromString(s.ErrorType, spec)
	if err != nil {
		return err
	}
	s.ParsedErrorType = ec
	return s.buildInternal(spec)
}

func (e *ErrorSpec) buildInternal(spec *Spec) error {
	newSafeArgs, err := buildInternalFieldsFromInterface(spec, e.SafeArgs, false)
	if err != nil {
		return err
	}
	newUnsafeArgs, err := buildInternalFieldsFromInterface(spec, e.UnsafeArgs, false)
	if err != nil {
		return err
	}
	args := map[string]*ParsedField{}

	args, err = mergeInto(args, newSafeArgs)
	if err != nil {
		return err
	}
	args, err = mergeInto(args, newUnsafeArgs)
	if err != nil {
		return err
	}
	e.ParsedArgs = args
	return nil
}
