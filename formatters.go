package logging

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	StandardTimeFormat = "2006-01-02 15:04:05.000Z07:00"
	tolerance          = 25 * time.Millisecond
)

type Formatter interface {
	Format(record *Record) string
}

type StandardFormatter struct {
	TimeFormat string
}

func (formatter *StandardFormatter) Format(record *Record) string {
	var message string
	now := time.Now()
	if now.Sub(record.Time) <= tolerance {
		message = record.Message
	} else {
		message = fmt.Sprintf("[%v] %v", record.Time.Format(formatter.TimeFormat), record.Message)
	}
	return fmt.Sprintf("[%v] %v %v", now.Format(formatter.TimeFormat), record.Level, message)
}

type JsonFormatter struct{}

func (formatter *JsonFormatter) Format(record *Record) string {
	data, _ := json.Marshal(record)
	return string(data)
}
