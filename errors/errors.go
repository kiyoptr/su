package errors

import (
	"errors"
	gerr "errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Error is a lightweight error struct with context
type Error struct {
	Source  string
	Message error
	Inner   error
}

func (e *Error) Error() string {
	var err error
	err = e

	builder := new(strings.Builder)

	first := true
	addPrefix := func() {
		if first {
			builder.WriteString("Error: ")
			first = false
		} else {
			builder.WriteString("  - ")
		}
	}

	for err != nil {
		var cast *Error
		if As(err, &cast) {
			addPrefix()
			builder.WriteString(fmt.Sprintf("%v: %v", cast.Source, cast.Message))
			err = cast.Unwrap()
		} else {
			addPrefix()

			errT := reflect.TypeOf(err)
			if errT.Kind() == reflect.Ptr {
				errT = errT.Elem()
			}

			n := errT.Name()
			builder.WriteString(fmt.Sprintf("[%v] %s", n, err.Error()))
			err = nil
		}

		if err != nil {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func (e *Error) Unwrap() error { return e.Inner }

// getCallerInfo returns the file and line that called any of New functions as string
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "<no source>"
	}

	return fmt.Sprintf("%v:%v", file, line)
}

// New constructs a new Error
func New(msg string) error {
	return &Error{
		Source:  getCallerInfo(),
		Message: errors.New(msg),
	}
}

// Newi attaches a new Error to an existing error to give it context
func Newi(inner error, msg string) error {
	return &Error{
		Source:  getCallerInfo(),
		Message: errors.New(msg),
		Inner:   inner,
	}
}

// Newf is same as New but it formats the message
func Newf(format string, params ...interface{}) error {
	return &Error{
		Source:  getCallerInfo(),
		Message: errors.New(fmt.Sprintf(format, params...)),
	}
}

// Newfi is same as Newi but with formatted message
func Newfi(inner error, format string, params ...interface{}) error {
	return &Error{
		Source:  getCallerInfo(),
		Message: errors.New(fmt.Sprintf(format, params...)),
		Inner:   inner,
	}
}

// As is a wrapper around go's standard errors.As
func As(err error, target interface{}) bool { return gerr.As(err, target) }

// Is is a wrapper around go's standard errors.Is
func Is(err, target error) bool { return gerr.Is(err, target) }

// Unwrap is a wrapper around go's standard errors.Unwrap
func Unwrap(err error) error { return gerr.Unwrap(err) }
