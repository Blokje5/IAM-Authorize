package log

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.Llongfile|log.Lmicroseconds)

func Debug(v ...interface{}) {
	logger.Println(v...)
}

func Debugf(format string, v ...interface{}) {
	logger.Println(v...)
}

func Info(v ...interface{}) {
	logger.Println(v...)
}

func Infof(format string, v ...interface{}) {
	logger.Println(v...)
}

func Warn(v ...interface{}) {
	logger.Println(v...)
}

func Warnf(format string, v ...interface{}) {
	logger.Println(v...)
}

func Error(v ...interface{}) {
	logger.Println(v...)
}

func Errorf(format string, v ...interface{}) {
	logger.Println(v...)
}
