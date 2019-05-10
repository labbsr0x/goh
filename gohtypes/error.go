package gohtypes

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Error groups together information that defines an error. Should always be used to
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Err     error  `json:"err"`
}

// Error() gives a string representing the error; also, forces the Error type to comply with the error interface
// If the RUNTIME_MODE variable is defined and has the value DEBUG, more information is emitted to the final error message
func (e *Error) Error() string {
	return fmt.Sprintf("ERROR (%v): %s\n", e.Code, e.Message)
}

// InnerError gives a string with the details of an inner error, if it exists
func (e *Error) InnerError() string {
	iErr := ""
	if e.Err != nil {
		iErr = fmt.Sprintf("Inner Error: '%v'", e.Err)
	}
	return fmt.Sprintf("ERROR (%v): %s;\n%s", e.Code, e.Message, iErr)
}

// PanicIfError is just a wrapper to a panic call that propagates a custom Error when the err property != nil
func PanicIfError(message string, code int, err error) {
	e := Error{Message: message, Code: code, Err: err}
	if e.Err != nil {
		logrus.Errorf(e.InnerError())
		panic(e)
	}
}

// Panic wraps a panic call propagating the given error parameter
func Panic(message string, code int) {
	e := Error{Message: message, Code: code}
	logrus.Errorf(e.InnerError())
	panic(e)
}
