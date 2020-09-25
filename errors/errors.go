package errors

import (
	"errors"
	goerrors "errors"
	"fmt"
	"runtime"
)

// Error is a lightweight error struct with context
type Error struct {
	Source  string
	Message error
	Inner   error
}

func NewError(source string, message, inner error) *Error {
	return &Error{
		Source:  source,
		Message: message,
		Inner:   inner,
	}
}

// Each iterates all inner errors as long as they're Error, starting from itself
func (e *Error) Each(it func(err error) bool) {
	if it == nil {
		return
	}

	var current error = e
	for current != nil {
		var cast *Error
		if As(current, &cast) {
			current = cast.Unwrap()
		} else {
			current = nil
		}

		if !it(current) {
			break
		}
	}
}

// StackTrace builds the stack trace of all inner errors of Error
func (e *Error) StackTrace() (list []string) {
	list = make([]string, 0, 5)

	e.Each(func(err error) bool {
		list = append(list, err.Error())
		return true
	})

	return
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v: %v", e.Source, e.Message)
}

func (e *Error) Unwrap() error { return e.Inner }

// getCallerInfo returns the file and line that called any of New functions as string
// skipFrames parameter defines how many functions to skip
func getCallerInfo(skipFrames int) string {
	_, file, line, ok := runtime.Caller(2 + skipFrames)
	if !ok {
		return "<no source>"
	}

	return fmt.Sprintf("%v:%v", file, line)
}

// New constructs a new Error
func New(msg string) error {
	return NewError(getCallerInfo(0), errors.New(msg), nil)
}

// Newi attaches an existing error to a new error
// This is used to provide an easier way for wrapping errors and stack trace
func Newi(inner error, msg string) error {
	return NewError(getCallerInfo(0), errors.New(msg), inner)
}

// Newf constructs a formatted error
func Newf(format string, params ...interface{}) error {
	return NewError(getCallerInfo(0), errors.New(fmt.Sprintf(format, params...)), nil)
}

// Newif constructs a new formatted error with an attached inner error
func Newif(inner error, format string, params ...interface{}) error {
	return NewError(getCallerInfo(0), errors.New(fmt.Sprintf(format, params...)), inner)
}

// News constructs a new Error and skips given frames for getting stack info.
func News(skip int, msg string) error {
	return NewError(getCallerInfo(skip), errors.New(msg), nil)
}

func Newsi(skip int, inner error, msg string) error {
	return NewError(getCallerInfo(skip), errors.New(msg), inner)
}

func Newsf(skip int, format string, params ...interface{}) error {
	return NewError(getCallerInfo(skip), errors.New(fmt.Sprintf(format, params...)), nil)
}

func Newsif(skip int, inner error, format string, params ...interface{}) error {
	return NewError(getCallerInfo(skip), errors.New(fmt.Sprintf(format, params...)), inner)
}

// As is a wrapper around go's standard errors.As
func As(err error, target interface{}) bool { return goerrors.As(err, target) }

// Is is a wrapper around go's standard errors.Is
func Is(err, target error) bool { return goerrors.Is(err, target) }

// Unwrap is a wrapper around go's standard errors.Unwrap
func Unwrap(err error) error { return goerrors.Unwrap(err) }
