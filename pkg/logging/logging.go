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
	log.SetReportCaller(false)
	log.SetOutput(os.Stdout)
	// log.Formatter = &logrus.JSONFormatter{}

	return &TaskServiceLogger{Log: log}
}

// Infof
func Infof(format string, v ...interface{}) {
	lock.Lock()
	log.Infof(format, v...)
	lock.Unlock()
}

// Warnf
func Warnf(format string, v ...interface{}) {
	lock.Lock()
	log.Warnf(format, v...)
	lock.Unlock()
}

// Errorf
func Errorf(format string, v ...interface{}) {
	lock.Lock()
	log.Errorf(format, v...)
	lock.Unlock()
}

// Fatalf
func Fatalf(format string, v ...interface{}) {
	lock.Lock()
	logging.Fatalf(format, v...)
	lock.Unlock()
}
