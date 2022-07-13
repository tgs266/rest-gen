package errors

import (
	"fmt"

	"github.com/google/uuid"
)

type Code = string
type StatusCode = int

type ErrorCode struct {
	Code
	StatusCode
}

var (
	INVALID_ARGUMENT = ErrorCode{"INVALID_ARGUMENT", 400}
	NOT_FOUND        = ErrorCode{"NOT_FOUND", 404}
	INTERNAL         = ErrorCode{"INTERNAL", 500}
	UNAUTHORIZED     = ErrorCode{"UNAUTHORIZED", 500}
)

var KnownErrorCode = map[string]ErrorCode{
	"INVALID_ARGUMENT": INVALID_ARGUMENT,
	"NOT_FOUND":        NOT_FOUND,
	"INTERNAL":         INTERNAL,
	"UNAUTHORIZED":     UNAUTHORIZED,
}

type StandardErrorInterface interface {
	Error() string
	Cause() string
}

type StandardError struct {
	ErrorId    string                 `json:"errorId"`
	Code       string                 `json:"code"`
	StatusCode int                    `json:"-"`
	Name       string                 `json:"name"`
	Params     map[string]interface{} `json:"params"`
	cause      error
}

type name = string
type value = interface{}

type Param struct {
	name
	value
}

func Wrap(err error, name string, errorCode ErrorCode, params ...Param) StandardError {
	joinedParams := map[string]interface{}{}
	for _, p := range params {
		joinedParams[p.name] = p.value
	}
	return StandardError{
		ErrorId:    uuid.New().String(),
		Code:       errorCode.Code,
		Name:       name,
		StatusCode: errorCode.StatusCode,
		Params:     joinedParams,
		cause:      err,
	}
}

func (se StandardError) Error() string {
	return fmt.Sprintf("%s: (%s)", se.Code, se.ErrorId)
}

func (se StandardError) Cause() error {
	return se.cause
}

func IsStandardError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(StandardErrorInterface)
	return ok
}

func NewInvalidArgumentError(err error) StandardError {
	return Wrap(err, "InvalidArgument", INVALID_ARGUMENT)
}
