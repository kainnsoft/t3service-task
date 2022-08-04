// based on standard logger from log package
package logging

import (
	"fmt"
	standard_log "log"
	"os"
)

const (
	cInfo  string = "info"
	cDebug string = "debug"
	cWarn  string = "warn"
	cError string = "error"
	cFatal string = "fatal"
)

// Interface -.
type Interface22 interface {
	Debug(message interface{}, args ...interface{})
	Info(message string, args ...interface{})
	Warn(message interface{}, args ...interface{})
	Error(message interface{}, args ...interface{})
	Fatal(message interface{}, args ...interface{})
}

type TaskLogger struct {
	infoLog  *standard_log.Logger
	debugLog *standard_log.Logger
	warnLog  *standard_log.Logger
	ErrorLog *standard_log.Logger
	fatalLog *standard_log.Logger
}

var _ Interface22 = (*TaskLogger)(nil)

func NewTaskLogger() *TaskLogger { // infoOut io.Writer, errorOut io.Writer, debugOut io.Writer) TaskLogger {
	infoLog := standard_log.New(os.Stdout, "INFO\t", standard_log.Ldate|standard_log.Ltime)
	debugLog := standard_log.New(os.Stdout, "DEBUG\t", standard_log.Ldate|standard_log.Ltime|standard_log.Lshortfile)
	warnLog := standard_log.New(os.Stderr, "WARN\t", standard_log.Ldate|standard_log.Ltime|standard_log.Lshortfile)
	errorLog := standard_log.New(os.Stderr, "ERROR\t", standard_log.Ldate|standard_log.Ltime|standard_log.Lshortfile)
	fatalLog := standard_log.New(os.Stderr, "FATAL\t", standard_log.Ldate|standard_log.Ltime|standard_log.Lshortfile)
	taskLogger := &TaskLogger{
		infoLog:  infoLog,
		debugLog: debugLog,
		warnLog:  warnLog,
		ErrorLog: errorLog,
		fatalLog: fatalLog,
	}

	return taskLogger
}

// Debug -.
func (tl *TaskLogger) Debug(message interface{}, args ...interface{}) {
	tl.msg(cDebug, message, args...)
}

// Info -.
func (tl *TaskLogger) Info(message string, args ...interface{}) {
	tl.msg(cInfo, message, args...)
}

// Warn -.
func (tl *TaskLogger) Warn(message interface{}, args ...interface{}) {
	tl.msg(cWarn, message, args...)
}

// Error -.
func (tl *TaskLogger) Error(message interface{}, args ...interface{}) {
	tl.msg(cError, message, args...)
}

// Fatal -.
func (tl *TaskLogger) Fatal(message interface{}, args ...interface{}) {
	tl.msg(cFatal, message, args...)
	os.Exit(1)
}

func (tl *TaskLogger) msg(level string, message interface{}, args ...interface{}) {
	switch msg := message.(type) {
	case error:
		tl.log(level, msg.Error(), args...)
	case string:
		tl.log(level, msg, args...)
	default:
		tl.debugLog.Printf(fmt.Sprintf("%s message %v has unknown type %v", level, message, msg), args...)
	}
}

func (tl *TaskLogger) log(level, msg string, args ...interface{}) {
	switch level {
	case cInfo:
		tl.infoLog.Printf(msg, args...)
	case cDebug:
		tl.debugLog.Printf(msg, args...)
	case cWarn:
		tl.warnLog.Printf(msg, args...)
	case cError:
		tl.ErrorLog.Printf(msg, args...)
	case cFatal:
		tl.fatalLog.Printf(msg, args...)
	}
}
