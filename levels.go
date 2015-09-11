package logging

import (
	"encoding/json"
)

// Levels of logging supported by library.
// They are ordered in descending order or imporance.
const (
	EMERGENCY Level = iota
	ALERT
	CRITICAL
	ERROR
	WARNING
	NOTICE
	INFO
	DEBUG
)

const (
	stopped uint8 = iota
	paused
	running
)

// Level represents log level for log message.
type Level int8

// String returns stirng representation of log level.
func (level Level) String() string {
	switch level {
	case EMERGENCY:
		return "EMERGENCY"
	case ALERT:
		return "ALERT"
	case CRITICAL:
		return "CRITICAL"
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case NOTICE:
		return "NOTICE"
	case INFO:
		return "INFO"
	default:
		return "DEBUG"
	}
}

// MarshalJSON is implementation of json.Marshaler interface, will be used when
// log level is serialized to json.
func (level Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(level.String())
}
