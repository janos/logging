package logging // import "resenje.org/loggign"

import (
	"container/ring"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	loggers = make(map[string]*Logger)
	lock    = &sync.RWMutex{}
)

type Logger struct {
	Name          string
	Level         Level
	Handlers      []Handler
	buffer        *ring.Ring
	stateChannel  chan uint8
	recordChannel chan *Record
	waiter        sync.WaitGroup
	lock          sync.RWMutex
	countIn       uint64
	countOut      uint64
}

func NewLogger(name string, level Level, handlers []Handler, bufferLength int) (*Logger, error) {
	lock.Lock()
	defer lock.Unlock()
	if _, ok := loggers[name]; ok {
		return nil, errors.New("Logger with that name already exists")
	}
	logger := &Logger{
		Name:          name,
		Level:         level,
		Handlers:      handlers,
		buffer:        ring.New(bufferLength),
		stateChannel:  make(chan uint8, 0),
		recordChannel: make(chan *Record, 2048),
		waiter:        sync.WaitGroup{},
		lock:          sync.RWMutex{},
		countOut:      0,
		countIn:       0,
	}
	go logger.run()
	loggers[name] = logger
	return logger, nil
}

func GetLogger(name string) (*Logger, error) {
	logger, ok := loggers[name]
	if !ok {
		return nil, fmt.Errorf("Unknown logger %s", name)
	}
	return logger, nil
}

func RemoveLogger(name string) {
	lock.Lock()
	defer lock.Unlock()
	delete(loggers, name)
}

func RemoveLoggers() {
	lock.Lock()
	defer lock.Unlock()
	loggers = make(map[string]*Logger)
}

func WaitForAllUnprocessedRecords() {
	var wg sync.WaitGroup
	for _, logger := range loggers {
		wg.Add(1)
		go func(logger *Logger) {
			defer wg.Done()
			logger.WaitForUnprocessedRecords()
		}(logger)
	}
	wg.Wait()
}

func (logger *Logger) String() string {
	return logger.Name
}

func (logger *Logger) run() {
	defer func() {
		logger.WaitForUnprocessedRecords()
		logger.closeHandlers()
	}()
recordLoop:
	for {
		select {
		case record := <-logger.recordChannel:
			record.process()
		case state := <-logger.stateChannel:
			switch state {
			case stopped:
				logger.waiter.Done()
				break recordLoop
			case paused:
			stateLoop:
				for {
					select {
					case state := <-logger.stateChannel:
						switch state {
						case stopped:
							logger.waiter.Done()
							break recordLoop
						case running:
							break stateLoop
						default:
							continue
						}
					}
				}
			}
		}
	}
}

func (logger *Logger) WaitForUnprocessedRecords() {
	runtime.Gosched()
	logger.Unpause()
	var (
		diff     uint64
		diffPrev uint64
		i        uint8
	)
	for {
		diff = atomic.LoadUint64(&logger.countIn) - atomic.LoadUint64(&logger.countOut)
		if diff == diffPrev {
			i++
		}
		if i >= 100 {
			return
		}
		if diff > 0 {
			diffPrev = diff
			time.Sleep(10 * time.Millisecond)
		} else {
			return
		}
	}
}

func (logger *Logger) closeHandlers() {
	for _, handler := range logger.Handlers {
		handler.Close()
	}
}

func (logger *Logger) Pause() {
	logger.stateChannel <- paused
}

func (logger *Logger) Unpause() {
	logger.stateChannel <- running
}

func (logger *Logger) Stop() {
	logger.stateChannel <- stopped
	logger.waiter.Wait()
}

func (logger *Logger) SetBufferLength(length int) {
	logger.lock.Lock()
	defer logger.lock.Unlock()

	if length == 0 {
		logger.buffer = nil
	} else if length != logger.buffer.Len() {
		logger.buffer = ring.New(length)
	}
}

func (logger *Logger) AddHandler(handler Handler) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	logger.Handlers = append(logger.Handlers, handler)
	logger.flushBuffer()
}

func (logger *Logger) ClearHandlers() {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	logger.closeHandlers()
	logger.Handlers = make([]Handler, 0)
	logger.flushBuffer()
}

func (logger *Logger) SetLevel(level Level) {
	logger.lock.Lock()
	logger.Level = level
	logger.flushBuffer()
}

func (logger *Logger) flushBuffer() {
	if logger.buffer != nil {
		oldBuffer := logger.buffer
		logger.buffer = ring.New(oldBuffer.Len())

		go func() {
			oldBuffer.Do(func(x interface{}) {

				if x == nil {
					return
				}

				record := x.(*Record)

				atomic.AddUint64(&logger.countIn, 1)
				logger.recordChannel <- record
			})
		}()
	}
}

func (logger *Logger) log(level Level, format string, a ...interface{}) {
	var message string
	if format == "" {
		message = fmt.Sprint(a...)
	} else {
		message = fmt.Sprintf(format, a...)
	}

	record := &Record{
		Level:   level,
		Message: message,
		Time:    time.Now(),
		logger:  logger,
	}
	atomic.AddUint64(&logger.countIn, 1)
	logger.recordChannel <- record
}

func (logger *Logger) Logf(level Level, format string, a ...interface{}) {
	logger.log(level, format, a...)
}

func (logger *Logger) Log(level Level, a ...interface{}) {
	logger.log(level, "", a...)
}

func (logger *Logger) Emergencyf(format string, a ...interface{}) {
	logger.log(EMERGENCY, format, a...)
}

func (logger *Logger) Emergency(a ...interface{}) {
	logger.log(EMERGENCY, "", a...)
}

func (logger *Logger) Alertf(format string, a ...interface{}) {
	logger.log(ALERT, format, a...)
}

func (logger *Logger) Alert(a ...interface{}) {
	logger.log(ALERT, "", a...)
}

func (logger *Logger) Criticalf(format string, a ...interface{}) {
	logger.log(CRITICAL, format, a...)
}

func (logger *Logger) Critical(a ...interface{}) {
	logger.log(CRITICAL, "", a...)
}

func (logger *Logger) Errorf(format string, a ...interface{}) {
	logger.log(ERROR, format, a...)
}

func (logger *Logger) Error(a ...interface{}) {
	logger.log(ERROR, "", a...)
}

func (logger *Logger) Warningf(format string, a ...interface{}) {
	logger.log(WARNING, format, a...)
}

func (logger *Logger) Warning(a ...interface{}) {
	logger.log(WARNING, "", a...)
}

func (logger *Logger) Noticef(format string, a ...interface{}) {
	logger.log(NOTICE, format, a...)
}

func (logger *Logger) Notice(a ...interface{}) {
	logger.log(NOTICE, "", a...)
}

func (logger *Logger) Infof(format string, a ...interface{}) {
	logger.log(INFO, format, a...)
}

func (logger *Logger) Info(a ...interface{}) {
	logger.log(INFO, "", a...)
}

func (logger *Logger) Debugf(format string, a ...interface{}) {
	logger.log(DEBUG, format, a...)
}

func (logger *Logger) Debug(a ...interface{}) {
	logger.log(DEBUG, "", a...)
}
