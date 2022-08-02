package errors

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

const (
	NoType = ErrorType(iota)
	BadRequest
	NotFound
	MethodNotAllowed
)

type ErrorType uint
type taskErrorContext struct {
	Field   string
	Message string
}

type customTaskError struct {
	errorType     ErrorType
	originalError error
	StackTrace    []byte
	contextInfo   taskErrorContext
}

// Error returns the mssage of a customError
func (e customTaskError) Error() string {
	return e.originalError.Error()
}

func (e customTaskError) PrintWithStack() string {
	return fmt.Sprintf("error: %v, \n context: %v \n Stack trace %v \n", e.Error(), e.contextInfo, string(e.StackTrace))
}

//-------------------------------------------------------------------

// New creates a new customError
func (errorType ErrorType) New(msg string) customTaskError {
	return customTaskError{errorType: errorType, originalError: errors.New(msg)}
}

// New creates a new customError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) customTaskError {
	stack := getStack()
	return customTaskError{errorType: errorType,
		originalError: fmt.Errorf(msg, args...),
		StackTrace:    stack}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) customTaskError {
	return errorType.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) customTaskError {
	return customTaskError{errorType: errorType, originalError: errors.Wrapf(err, msg, args...)}
}

//-------------------------------------------------------------------

// New creates a new customError
func New(msg string) customTaskError {
	return customTaskError{errorType: NoType, originalError: errors.New(msg)}
}

// New creates a new customError with formatted message
func Newf(msg string, args ...interface{}) customTaskError {
	err := fmt.Errorf(msg, args...)

	return customTaskError{errorType: NoType, originalError: err}
}

// Wrap creates a new wrapped error
func Wrap(err error, msg string) customTaskError {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf wraps an error with format string
func Wrapf(err error, msg string, args ...interface{}) customTaskError {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customTaskError); ok {
		return customTaskError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			contextInfo:   customErr.contextInfo,
		}
	}

	return customTaskError{errorType: NoType, originalError: wrappedError}
}

//-------------------------------------------------------

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := taskErrorContext{Field: field, Message: message}
	if customErr, ok := err.(customTaskError); ok {
		return customTaskError{errorType: customErr.errorType, originalError: customErr.originalError, contextInfo: context}
	}

	return customTaskError{errorType: NoType, originalError: err, contextInfo: context}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := taskErrorContext{}
	customErr, ok := err.(customTaskError)
	if ok {
		if customErr.contextInfo != emptyContext {
			return map[string]string{"field": customErr.contextInfo.Field, "message": customErr.contextInfo.Message}
		}
	}
	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customTaskError); ok {
		return customErr.errorType
	}

	return NoType
}

func CheckErrorType(err error) error {
	if customErr, ok := err.(customTaskError); ok {
		return customErr
	}

	return err
}

func getStack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	return buf
}
