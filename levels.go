package logging

import (
	"encoding/json"
)

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

type Level int8

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

func (level Level) MarshalJSON() ([]byte, error) {
	return json.Marshal(level.String())
}
