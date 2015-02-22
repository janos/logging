package logging

import (
	"sync/atomic"
	"time"
)

type Record struct {
	Time    time.Time `json:"time"`
	Level   Level     `json:"level"`
	Message string    `json:"message"`
	logger  *Logger
}

func (record *Record) process() {
	logger := record.logger
	logger.lock.RLock()
	defer func() {
		logger.lock.RUnlock()
		atomic.AddUint64(&logger.countOut, 1)
	}()

	if record.Level <= logger.Level {
		for _, handler := range logger.Handlers {
			if record.Level <= handler.GetLevel() {
				if err := handler.Handle(record); err != nil {
					handler.HandleError(err)
				}
			}
		}
	}
}
