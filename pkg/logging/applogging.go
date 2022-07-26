// based on standard logger from log package
package logging

import (
	"fmt"
	standard_log "log"
	"os"
)

const (
	c_info  string = "info"
	c_debug string = "debug"
	c_warn  string = "warn"
	c_error string = "error"
	c_fatal string = "fatal"
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

func NewTaskLogger() *TaskLogger { //infoOut io.Writer, errorOut io.Writer, debugOut io.Writer) TaskLogger {
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
	tl.msg(c_debug, message, args...)
}

// Info -.
func (tl *TaskLogger) Info(message string, args ...interface{}) {
	tl.msg(c_info, message, args...)
}

// Warn -.
func (tl *TaskLogger) Warn(message interface{}, args ...interface{}) {
	tl.msg(c_warn, message, args...)
}

// Error -.
func (tl *TaskLogger) Error(message interface{}, args ...interface{}) {
	tl.msg(c_error, message, args...)
}

// Fatal -.
func (tl *TaskLogger) Fatal(message interface{}, args ...interface{}) {
	tl.msg(c_fatal, message, args...)
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

func (tl *TaskLogger) log(level string, msg string, args ...interface{}) {
	switch level {
	case c_info:
		tl.infoLog.Printf(msg, args...)
	case c_debug:
		tl.debugLog.Printf(msg, args...)
	case c_warn:
		tl.warnLog.Printf(msg, args...)
	case c_error:
		tl.ErrorLog.Printf(msg, args...)
	case c_fatal:
		tl.fatalLog.Printf(msg, args...)
	}
}
