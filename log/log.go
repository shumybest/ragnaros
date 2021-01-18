package log

import (
    "github.com/sirupsen/logrus"
    "sync"
)

type Logger struct {
    *logrus.Logger
}

var loggerInstance *Logger
var Once sync.Once

func GetLoggerInstance() *Logger {
    Once.Do(func() {
        loggerInstance = &Logger{}
        loggerInstance.Logger = logrus.New()
    })
    return loggerInstance
}

func (l *Logger) InitConfig() {
    l.SetReportCaller(false)

    // TODO logging level from configuration
    l.SetLevel(logrus.DebugLevel)

    // TODO initialize logstash here
}