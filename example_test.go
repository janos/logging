package logging

import "fmt"

// Utility stuff that is not very useful to show as part of example tests, but
// that are needed in order to execute examples successful

// custom formatter that does not modifiy provided message
type PassThroughFormatter struct{}

// Formatter interface implementation for PassThroughFormatter
func (formatter *PassThroughFormatter) Format(record *Record) string {
	return record.Message
}

// Shows basic example of using library without any setup.
func Example_basic() {
	// most basic way of using logging is by just calling logging functions
	// without any setup. All these messages will be written to stdout.
	Debug("Debug message")
	Debugf("Debug message %s", "with formatting")
	Info("Info message")
	Infof("Info message %s", "with formatting")

}

// This example show how to create logger with specific name and custom
// handlers.
func Example_custom() {

	// handler that writes all log messages to memory
	// PassThroughFormatter only returns messages provided, without any modification
	memoryHandler := &MemoryHandler{
		Level:     WARNING,
		Formatter: &PassThroughFormatter{},
	}

	if _, err := NewLogger("myLogger", WARNING, []Handler{memoryHandler}, 0); err != nil {
		panic(err)
	}

	// this will probably be called somewhere else in the code
	if logger, err := GetLogger("myLogger"); err != nil {
		panic(err)
	} else {
		logger.Infof("Should not be logged")
		logger.Error("Should be logged")

	}

	// wait for all messages to be processed, this is blocking
	WaitForAllUnprocessedRecords()

	fmt.Printf("%v\n", memoryHandler.Messages)

	// Output: [Should be logged]
}
