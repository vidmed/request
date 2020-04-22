package request

import (
	"github.com/sirupsen/logrus"
)

var logInstance *logrus.Logger

func GetLogger() *logrus.Logger {
	if logInstance == nil {
		return logrus.New()
	}
	return logInstance
}

func InitLogger(level int) {
	logInstance = logrus.New()
	logInstance.Level = logrus.AllLevels[level]
}
