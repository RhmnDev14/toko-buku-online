package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log  *logrus.Logger
	file *os.File
}

func NewLogger() *Logger {
	log := logrus.New()

	file, err := os.OpenFile("TOKO_BUKU_ONLINE.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("failed to open log file: ", err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	log.SetOutput(mw)

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger := &Logger{log: log, file: file}

	return logger
}

func (l *Logger) Info(message string, data any) {
	if l == nil || l.log == nil {
		logrus.WithFields(logrus.Fields{
			"message": message,
			"data":    data,
		}).Warn("Logger not initialized, fallback to default logrus")
		return
	}
	l.log.WithFields(logrus.Fields{
		"data": data,
	}).Info(message)
}

func (l *Logger) Error(message string, data any) {
	if l == nil || l.log == nil {
		logrus.WithFields(logrus.Fields{
			"message": message,
			"data":    data,
		}).Warn("Logger not initialized, fallback to default logrus")
		return
	}
	l.log.WithFields(logrus.Fields{
		"data": data,
	}).Error(message)
}
