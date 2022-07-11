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
	ErrorId    string                 `json:"errorId"`
	Code       string                 `json:"code"`
	StatusCode int                    `json:""`
	Name       string                 `json:"name"`
	Params     map[string]interface{} `json:"params"`
	cause      error
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
	_, ok := err.(StandardError)
	return ok
}

func NewInvalidArgumentError(err error) StandardError {
	return StandardError{
		Code:       INVALID_ARGUMENT,
		StatusCode: 400,
		Name:       "InvalidArgument",
		ErrorId:    uuid.New().String(),
		cause:      err,
	}
}
