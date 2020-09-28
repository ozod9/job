package logger

import (
	"os"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger  *logrus.Logger
	host    string
	guid    string
	app     string
	version string
}

func (logger *Logger) Info(message string, source string) {
	logger.logger.WithFields(logrus.Fields{
		"GUID":    logger.guid,
		"Version": logger.version,
		"App":     logger.app,
		"Message": message,
		"Source":  source,
	}).Info(message)
}

func (logger *Logger) Error(message string, source string) {
	logger.logger.WithFields(logrus.Fields{
		"GUID":    logger.guid,
		"Version": logger.version,
		"App":     logger.app,
		"Message": message,
		"Source":  source,
	}).Error(message)
}

func (logger *Logger) Warning(message string, source string) {
	logger.logger.WithFields(logrus.Fields{
		"GUID":    logger.guid,
		"Version": logger.version,
		"App":     logger.app,
		"Message": message,
		"Source":  source,
	}).Warning(message)
}

func (logger *Logger) Fatal(message string, source string) {
	logger.logger.WithFields(logrus.Fields{
		"GUID":    logger.guid,
		"Version": logger.version,
		"App":     logger.app,
		"Message": message,
		"Source":  source,
	}).Fatal(message)
}

func (logger *Logger) SetHost(host string) *Logger {
	logger.host = host
	return logger
}

func (logger *Logger) SetVersion(version string) *Logger {
	logger.version = version
	return logger
}

func (logger *Logger) SetApp(app string) *Logger {
	logger.app = app
	return logger
}

func NewLogger() *Logger {
	logger := new(Logger)
	return logger
}

func (logger *Logger) CreateLogger() (*Logger, error) {
	logg := logrus.New()
	logg.SetOutput(os.Stdout)
	logger.guid = ksuid.New().String()
	logger.logger = logg
	return logger, nil
}
