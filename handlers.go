package logging

import (
	"io"
	"os"
	"sync"
)

type Handler interface {
	Handle(record *Record) error
	HandleError(error) error
	GetLevel() Level
	Close() error
}

type NullHandler struct{}

func (handler *NullHandler) Handle(record *Record) error {
	return nil
}

func (handler *NullHandler) HandleError(err error) error {
	os.Stderr.WriteString(err.Error())
	return nil
}

func (handler *NullHandler) GetLevel() Level {
	return DEBUG
}

func (handler *NullHandler) Close() error {
	return nil
}

type WriteHandler struct {
	NullHandler

	Level     Level
	Formatter Formatter
	Writer    io.Writer
	lock      sync.RWMutex
}

func (handler *WriteHandler) Handle(record *Record) error {
	handler.lock.Lock()
	defer handler.lock.Unlock()

	_, err := handler.Writer.Write([]byte(handler.Formatter.Format(record) + "\n"))
	return err
}

func (handler *WriteHandler) GetLevel() Level {
	return handler.Level
}

type MemoryHandler struct {
	NullHandler

	Level     Level
	Formatter Formatter
	Messages  []string
	lock      sync.RWMutex
}

func (handler *MemoryHandler) Handle(record *Record) error {
	handler.lock.Lock()
	defer handler.lock.Unlock()

	handler.Messages = append(handler.Messages, handler.Formatter.Format(record))
	return nil
}

func (handler *MemoryHandler) GetLevel() Level {
	return handler.Level
}
