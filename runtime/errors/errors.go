package errors

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

var (
	INVALID_ARGUMENT = ErrorType{"INVALID_ARGUMENT", 400}
	NOT_FOUND        = ErrorType{"NOT_FOUND", 404}
	INTERNAL         = ErrorType{"INTERNAL", 500}
	UNAUTHORIZED     = ErrorType{"UNAUTHORIZED", 401}
	FORBIDDEN        = ErrorType{"FORBIDDEN", 403}
)

var KnownErrorCode = map[string]ErrorType{
	"INVALID_ARGUMENT": INVALID_ARGUMENT,
	"NOT_FOUND":        NOT_FOUND,
	"INTERNAL":         INTERNAL,
	"UNAUTHORIZED":     UNAUTHORIZED,
}

type Error interface {
	Cause() error
	Code() int
	Name() string
	ErrorId() string
	Error() string
}

type ErrorType struct {
	name string
	code int
}

func (e ErrorType) Code() int {
	return e.code
}

func (e ErrorType) Name() string {
	return e.name
}

type standardError struct {
	errorId   string
	errorType ErrorType
	cause     error
	params    map[string]interface{}
}

type SerializableError struct {
	ErrorId    string                 `json:"errorId"`
	ErrorName  string                 `json:"errorName"`
	ErrorCode  int                    `json:"errorCode"`
	Parameters map[string]interface{} `json:"parameters"`
}

type name = string
type value = interface{}

type Param struct {
	name
	value
}

func Wrap(err error, name string, errorType ErrorType, params ...Param) standardError {
	joinedParams := map[string]interface{}{}
	for _, p := range params {
		joinedParams[p.name] = p.value
	}
	return standardError{
		errorId:   uuid.New().String(),
		errorType: errorType,
		params:    joinedParams,
		cause:     err,
	}
}

func (se standardError) Error() string {
	return fmt.Sprintf("%s: (%s)", se.errorType.name, se.errorId)
}

func (se standardError) Cause() error {
	return se.cause
}

func (se standardError) Code() int {
	return se.errorType.code
}

func (se standardError) Name() string {
	return se.errorType.name
}

func (se standardError) ErrorId() string {
	return se.errorId
}

func (se standardError) MarshalJSON() ([]byte, error) {
	return json.Marshal(SerializableError{
		ErrorCode:  se.errorType.code,
		ErrorName:  se.errorType.name,
		ErrorId:    se.errorId,
		Parameters: se.params,
	})
}

func (e *standardError) UnmarshalJSON(data []byte) (err error) {
	var se SerializableError
	if err := json.Unmarshal(data, &se); err != nil {
		return err
	}
	e.errorType = ErrorType{name: se.ErrorName, code: se.ErrorCode}
	e.errorId = se.ErrorId
	e.params = se.Parameters
	return nil
}

func IsErrorOfType(err error, errorType ErrorType) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(Error); ok {
		return e.Code() == errorType.code && e.Name() == errorType.name
	}
	return false
}

func GetError(err error) Error {
	if v, ok := err.(Error); ok {
		return v
	}
	return standardError{
		errorId:   uuid.NewString(),
		errorType: INTERNAL,
		cause:     err,
	}
}

func NewInvalidArgumentError(err error) error {
	return Wrap(err, "InvalidArgument", INVALID_ARGUMENT)
}
func NewInternalError(err error) error {
	return Wrap(err, "Internal", INTERNAL)
}
func NewNotFound(err error) error {
	return Wrap(err, "NotFound", NOT_FOUND)
}
func NewForbidden(err error) error {
	return Wrap(err, "Forbidden", FORBIDDEN)
}
func NewUnauthorized(err error) error {
	return Wrap(err, "Unauthorized", UNAUTHORIZED)
}
