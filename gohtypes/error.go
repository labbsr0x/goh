package gohtypes

import (
	"fmt"
	"os"
	"strings"

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
	mode := strings.Trim(os.Getenv("RUNTIME_MODE"), " ")
	if mode == "DEBUG" {
		return fmt.Sprintf("ERROR (%v): %s; \n Inner error: %s", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("ERROR (%v): %s", e.Code, e.Message)
}

// PanicIfError is just a wrapper to a panic call that propagates a custom Error when the err property != nil
func PanicIfError(message string, code int, err error) {
	e := Error{Message: message, Code: code, Err: err}
	if e.Err != nil {
		logrus.Errorf(e.Error())
		panic(e)
	}
}

// Panic wraps a panic call propagating the given error parameter
func Panic(message string, code int) {
	e := Error{Message: message, Code: code}
	logrus.Errorf(e.Error())
	panic(e)
}
