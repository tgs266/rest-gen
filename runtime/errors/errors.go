package errors

import (
	"fmt"

	"github.com/google/uuid"
)

type Code = string

const (
	INVALID_ARGUMENT      Code = "INVALID_ARGUMENT"
	INTERNAL_SERVER_ERROR Code = "INTERNAL_SERVER_ERROR"
)

type StandardError struct {
	ErrorId string                 `json:"errorId"`
	Code    string                 `json:"code"`
	Name    string                 `json:"name"`
	Params  map[string]interface{} `json:"params"`
	cause   error
}

func (se StandardError) Error() string {
	return fmt.Sprintf("%s: (%s)", se.Code, se.ErrorId)
}

func (se StandardError) Cause() error {
	return se.cause
}

func NewInvalidArgumentError(err error) StandardError {
	return StandardError{
		Code:    INVALID_ARGUMENT,
		Name:    "InvalidArgument",
		ErrorId: uuid.New().String(),
		cause:   err,
	}
}
