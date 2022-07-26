package logging

import (
	logging "log"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger
var lock = sync.Mutex{}

type TaskServiceLogger struct {
	Log *logrus.Logger
}

func NewLogging() *TaskServiceLogger {
	log = logrus.New()
	//log.Formatter = &logrus.JSONFormatter{}
	log.SetReportCaller(false)
	log.SetOutput(os.Stdout)

	return &TaskServiceLogger{Log: log}
}

// Info ...
func Info(format string, v ...interface{}) {
	lock.Lock()
	log.Infof(format, v...)
	lock.Unlock()
}

// Warn ...
func Warn(format string, v ...interface{}) {
	lock.Lock()
	log.Warnf(format, v...)
	lock.Unlock()
}

// Error ...
func Error(format string, v ...interface{}) {
	lock.Lock()
	log.Errorf(format, v...)
	lock.Unlock()
}

// Fatal ...
func Fatal(format string, v ...interface{}) {
	lock.Lock()
	logging.Fatalf(format, v...)
	lock.Unlock()
}
