package logging

import (
	"os"
)

func init() {
	InitDefaultLogger()
}

func InitDefaultLogger() {
	var err error
	_, err = NewLogger("default", DEBUG, []Handler{
		&WriteHandler{
			Level:     DEBUG,
			Formatter: &StandardFormatter{TimeFormat: StandardTimeFormat},
			Writer:    os.Stdout,
		},
	}, 0)
	if err != nil {
		panic(err)
	}
}

func Pause() {
	if logger, ok := loggers["default"]; ok {
		logger.Pause()
	}
}

func Unpause() {
	if logger, ok := loggers["default"]; ok {
		logger.Unpause()
	}
}

func Stop() {
	if logger, ok := loggers["default"]; ok {
		logger.Stop()
	}
}

func SetLevel(level Level) {
	if logger, ok := loggers["default"]; ok {
		logger.SetLevel(level)
	}
}

func SetBufferLength(length int) {
	if logger, ok := loggers["default"]; ok {
		logger.SetBufferLength(length)
	}
}

func AddHandler(handler Handler) {
	if logger, ok := loggers["default"]; ok {
		logger.AddHandler(handler)
	}
}

func ClearHandlers() {
	if logger, ok := loggers["default"]; ok {
		logger.ClearHandlers()
	}
}

func Logf(level Level, format string, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(level, format, a...)
	}
}

func Log(level Level, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(level, "", a...)
	}
}

func Emergencyf(format string, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(EMERGENCY, format, a...)
	}
}

func Emergency(a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(EMERGENCY, "", a...)
	}
}

func Alertf(format string, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(ALERT, format, a...)
	}
}

func Alert(a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(ALERT, "", a...)
	}
}

func Criticalf(format string, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(CRITICAL, format, a...)
	}
}

func Critical(a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(CRITICAL, "", a...)
	}
}

func Errorf(format string, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(ERROR, format, a...)
	}
}

func Error(a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(ERROR, "", a...)
	}
}

func Warningf(format string, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(WARNING, format, a...)
	}
}

func Warning(a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(WARNING, "", a...)
	}
}

func Noticef(format string, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(NOTICE, format, a...)
	}
}

func Notice(a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(NOTICE, "", a...)
	}
}

func Infof(format string, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(INFO, format, a...)
	}
}

func Info(a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(INFO, "", a...)
	}
}

func Debugf(format string, a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(DEBUG, format, a...)
	}
}

func Debug(a ...interface{}) {
	if logger, ok := loggers["default"]; ok {
		logger.log(DEBUG, "", a...)
	}
}
