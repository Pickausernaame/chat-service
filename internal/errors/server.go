package errors

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ServerError is used to return custom error codes to client.
type ServerError struct {
	Code    int
	Message string
	cause   error
}

func NewServerError(code int, msg string, err error) *ServerError {
	return &ServerError{
		Code:    code,
		Message: msg,
		cause:   err,
	}
}

func (s *ServerError) Is(target error) bool {
	return errors.Is(s.cause, target)
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("%s: %v", s.Message, s.cause)
}

func GetServerErrorCode(err error) int {
	code, _, _ := ProcessServerError(err)
	return code
}

// ProcessServerError tries to retrieve from given error it's code, message and some details.
// For example, that fields can be used to build error response for client.
func ProcessServerError(err error) (code int, msg string, details string) {
	sErr := &ServerError{}
	if errors.As(err, &sErr) {
		return sErr.Code, sErr.Message, sErr.Error()
	}

	echoErr := &echo.HTTPError{}
	if errors.As(err, &echoErr) {
		return echoErr.Code, echoErr.Message.(string), echoErr.Error()
	}

	return http.StatusInternalServerError, "something went wrong", err.Error()
}
