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

type CustomTaskError struct {
	errorType     ErrorType
	originalError error
	StackTrace    []byte
	contextInfo   taskErrorContext
}

// Error returns the mssage of a customError
func (e CustomTaskError) Error() string {
	return e.originalError.Error()
}

func (e *CustomTaskError) PrintWithStack() string {
	return fmt.Sprintf("error: %v, \n context: %v \n Stack trace %v \n", e.Error(), e.contextInfo, string(e.StackTrace))
}

//-------------------------------------------------------------------

// New creates a new customError
func (errorType ErrorType) New(msg string) CustomTaskError {
	return CustomTaskError{errorType: errorType, originalError: errors.New(msg)}
}

// New creates a new customError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) CustomTaskError {
	stack := getStack()
	return CustomTaskError{errorType: errorType,
		originalError: fmt.Errorf(msg, args...),
		StackTrace:    stack}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) CustomTaskError {
	return errorType.Wrapf(err, msg)
}

// Wrap creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) CustomTaskError {
	return CustomTaskError{errorType: errorType, originalError: errors.Wrapf(err, msg, args...)}
}

//-------------------------------------------------------------------

// New creates a new customError
func New(msg string) CustomTaskError {
	return CustomTaskError{errorType: NoType, originalError: errors.New(msg)}
}

// New creates a new customError with formatted message
func Newf(msg string, args ...interface{}) CustomTaskError {
	err := fmt.Errorf(msg, args...)

	return CustomTaskError{errorType: NoType, originalError: err}
}

// Wrap creates a new wrapped error
func Wrap(err error, msg string) CustomTaskError {
	return Wrapf(err, msg)
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// Wrapf wraps an error with format string
func Wrapf(err error, msg string, args ...interface{}) CustomTaskError {
	wrappedError := errors.Wrapf(err, msg, args...)

	if customErr, ok := err.(CustomTaskError); ok {
		return CustomTaskError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			contextInfo:   customErr.contextInfo,
		}
	}

	return CustomTaskError{errorType: NoType, originalError: wrappedError}
}

//-------------------------------------------------------

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := taskErrorContext{Field: field, Message: message}

	if customErr, ok := err.(CustomTaskError); ok {
		return CustomTaskError{errorType: customErr.errorType, originalError: customErr.originalError, contextInfo: context}
	}

	return CustomTaskError{errorType: NoType, originalError: err, contextInfo: context}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := taskErrorContext{}
	customErr, ok := err.(CustomTaskError)

	if ok {
		if customErr.contextInfo != emptyContext {
			return map[string]string{"field": customErr.contextInfo.Field, "message": customErr.contextInfo.Message}
		}
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(CustomTaskError); ok {
		return customErr.errorType
	}

	return NoType
}

func CheckErrorType(err error) error {
	if customErr, ok := err.(CustomTaskError); ok {
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
