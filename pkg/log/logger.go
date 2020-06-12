package log

import (
	"io"
	"log"
	"os"
)

var logger = New(os.Stdout)

func GetLogger() *Logger {
	return logger
}

type Logger struct {
	out    io.Writer
	logger *log.Logger
}

func New(w io.Writer) *Logger {
	return &Logger{
		out:    w,
		logger: log.New(w, "", log.Llongfile|log.Lmicroseconds),
	}
}

func NewWithFlags(w io.Writer, flag int) *Logger {
	return &Logger{
		out:    w,
		logger: log.New(w, "", flag),
	}
}

func (l *Logger) Debug(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Println(v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Println(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logger.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.logger.Println(v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logger.Println(v...)
}
