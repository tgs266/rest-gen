package spec

func (s *ErrorSpec) IsSafe(argName string) bool {
	_, exists := s.SafeArgs[argName]
	return exists
}

func (s *ErrorSpec) WriteErrorStruct(argName string) bool {
	_, exists := s.SafeArgs[argName]
	return exists
}

func (s *ErrorSpec) Parse(spec *Spec) error {
	if s.StatusCode <= 0 {
		s.StatusCode = 500
	}
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
