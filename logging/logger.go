package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

// Logger object
type Logger struct {
	Level LogLevel
	out   io.Writer
	file  *os.File
	log   *log.Logger
}

type logMessage struct {
	LogLevel     string `json:"loglevel"`
	Timestamp    string `json:"timestamp"`
	Message      string `json:"message"`
	Error        error  `json:"error,omitempty"`
	FunctionName string `json:"functionname"`
	FileName     string `json:"filename"`
	FileLine     int    `json:"fileline"`
}

// NewLogger Constructor
func NewLogger() (*Logger, error) {
	var file *os.File
	var err error
	var out io.Writer

	file, err = os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	out = os.Stdout

	l := &Logger{
		Level: DEBUG,
		out:   out,
		file:  file,
	}

	l.log = log.New(io.MultiWriter(out, file), "", 1)
	return l, nil
}

// logger error formatting function
func (l *Logger) Errorf(format string, err error, a ...any) {
	l.logf(ERROR, format, err, a...)
}

// logger error function
func (l *Logger) Error(err error, a ...any) {
	l.logf(ERROR, "", err, a...)
}

// logger warning formatting function
func (l *Logger) Warningf(format string, err error, a ...any) {
	l.logf(WARNING, format, err, a...)
}

// logger warning function
func (l *Logger) Warning(format string, a ...any) {
	l.logf(WARNING, format, nil, a...)
}

// logger info formatting function
func (l *Logger) Infof(format string, a ...any) {
	l.logf(INFO, format, nil, a...)
}

// logger function to log to stdout and to a file
func (l *Logger) logf(level LogLevel, format string, err error, a ...any) {

	// checks if logLevel is not below specified LogLevel
	if level < l.Level {
		return
	}

	// Add log message to the logger as it would be printf
	msg := fmt.Sprintf(format, a...)

	// Log to stdout
	funcName, file, line := getCaller()
	timestamp := time.Now().Format("15:04:05.000")

	// Create logMessage struct
	logMsg := &logMessage{
		LogLevel:     logLevelToString(level),
		Timestamp:    timestamp,
		Message:      msg,
		Error:        err,
		FunctionName: funcName,
		FileName:     file,
		FileLine:     line,
	}

	// Marshal logMessage struct to JSON
	jsonLogMsg, err := json.Marshal(logMsg)
	if err != nil {
		fmt.Println("Error marshaling log message to JSON:", err)
		return
	}

	// to stdout
	l.log.Println(string(jsonLogMsg))

	// Log to file
	if l.file != nil {
		if _, err := fmt.Fprintln(l.file, logMsg); err != nil {
			fmt.Println("Error writing to log file:", err)
		}
	}
}

// return loglevel string from leel
func logLevelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger Close
func (l *Logger) Close() error {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
			fmt.Println("Error closing log file:", err)
			return err
		}
	}
	return nil
}

// returns function invocations information
func getCaller() (string, string, int) {

	pc, fileName, lineNumber, ok := runtime.Caller(3)
	if !ok {
		return "unknown", "", 0
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown", "", 0
	}
	return fn.Name(), fileName, lineNumber
}
