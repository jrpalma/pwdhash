package logs

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

// Destination Represents the log file destination. Valid values are: STDOUT, STDERR, FILE.
type Destination string

const (
	// STDOUT Instructs the logger to print standard output.
	STDOUT Destination = "STDOUT"
	// STDERR Instructs the logger to print standard error.
	STDERR Destination = "STDERR"
	// FILE Instructs the logger to print to a file.
	FILE Destination = "FILE"
)

// LogLevel Represents the log level for the logs. Valid values are: ERROR, WARNING, INFO, and DEBUG.
type LogLevel string

const (
	// ERROR Prints errors messages only.
	ERROR LogLevel = "ERROR"
	// WARN Prints prints errors and warnings messages only.
	WARN LogLevel = "WARN"
	// INFO Prints prints errors, warnings, and info messages only.
	INFO LogLevel = "INFO"
	// DEBUG Prints errors, warnings, info, and debug messages.
	DEBUG LogLevel = "DEBUG"
)

// Logger Represents a logger that prints to a log stream.
type Logger interface {
	// Errorf Prints a formated error message to the log stream much like fmt.Printf.
	Errorf(format string, args ...interface{})
	// Warnf Prints a formated warning message to the log stream much like fmt.Printf.
	Warnf(format string, args ...interface{})
	// Infof Prints a formated info message to the log stream much like fmt.Printf.
	Infof(format string, args ...interface{})
	// Debugf Prints a formated debug message to the log stream much like fmt.Printf.
	Debugf(format string, args ...interface{})
}

// NewStreamLogger Creates a STDOUT or STDERR logger. This function returns an error
// if destination or the log level is invalid.
func NewStreamLogger(destination Destination, level LogLevel) (Logger, error) {

	if destination != STDOUT && destination != STDERR {
		return nil, fmt.Errorf("logs: Invalid log destination")
	}
	if level != ERROR && level != WARN && level != INFO && level != DEBUG {
		return nil, fmt.Errorf("logs: Invalid log level")
	}

	log := &logger{}
	log.level = logLevelToInt(level)

	if destination == STDOUT {
		log.file = os.Stdout
	} else {
		log.file = os.Stderr
	}

	return log, nil
}

// NewFileLogger Creates a FILE logger. This function returns an error
// if filePath cannot be opened or the log level is invalid.
func NewFileLogger(filePath string, level LogLevel) (Logger, error) {

	if level != ERROR && level != WARN && level != INFO && level != DEBUG {
		return nil, fmt.Errorf("logs: Invalid log level")
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_SYNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("logs: could not open file: %v", err)
	}

	log := &logger{file: file}
	log.level = logLevelToInt(level)

	return log, nil
}

type logger struct {
	file  io.Writer
	level int
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if errorLevel < l.level || l.file == nil {
		return
	}
	fmt.Fprint(l.file, getMessage(format, args...))
}
func (l *logger) Warnf(format string, args ...interface{}) {
	if warnLevel < l.level || l.file == nil {
		return
	}
	fmt.Fprint(l.file, getMessage(format, args...))
}
func (l *logger) Infof(format string, args ...interface{}) {
	if infoLevel < l.level || l.file == nil {
		return
	}
	fmt.Fprint(l.file, getMessage(format, args...))
}
func (l *logger) Debugf(format string, args ...interface{}) {
	if debugLevel < l.level || l.file == nil {
		return
	}
	fmt.Fprint(l.file, getMessage(format, args...))
}

const (
	errorLevel int = 4
	warnLevel  int = 3
	infoLevel  int = 2
	debugLevel int = 1
)

func logLevelToInt(level LogLevel) int {
	var value int
	switch level {
	case "ERROR":
		value = errorLevel
	case "WARN":
		value = warnLevel
	case "INFO":
		value = infoLevel
	case "DEBUG":
		value = debugLevel
	default:
		value = warnLevel
	}
	return value
}

func getMessage(format string, args ...interface{}) string {
	msg := time.Now().Format("2017-09-07 17:06:06 ")

	_, fileName, fileLine, ok := runtime.Caller(2)
	if ok {
		msg += fmt.Sprintf("%s:%d", fileName, fileLine)
	}

	msg += fmt.Sprintf(format+"\n", args...)

	return msg
}
