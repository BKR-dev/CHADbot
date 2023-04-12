package logging

import (
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
	l.logf(format, err, a...)
}

// logger error function
func (l *Logger) Error(err error, a ...any) {
	l.logf("", err, a...)
}

// logger info formatting function
func (l *Logger) Infof(format string, a ...any) {
	l.logf(format, nil, a...)
}

// TODO: add logg level for logger funcs
// TODO: add log rotate
// TODO: add formatting (json?)
// logger function to log to stdout and to a file
func (l *Logger) logf(format string, err error, a ...any) {

	// Add log message to the logger as it would be printf
	msg := fmt.Sprintf(format, a...)

	// Log to stdout
	funcName, file, line := getCaller()
	timestamp := time.Now().Format("15:04:05.000")

	// function does not contain nil
	if err != nil {
		// log message without error
		logMsg := fmt.Sprintf("%s: \nMessage: [%s] \n%s \n%s in line: %d\n", timestamp, msg, funcName, file, line)
		// to stdout
		l.log.Println(logMsg)

		// Log to file
		if l.file != nil {
			fmt.Fprintln(l.file, logMsg)
		}
	}

	// log message with error
	logMsg := fmt.Sprintf("%s: \nMessage: [%s] \n%s \n%s \n%s in line: %d\n", timestamp, err, msg, funcName, file, line)
	// to stdout
	l.log.Println(logMsg)

	// Log to file
	if l.file != nil {
		fmt.Fprintln(l.file, logMsg)
	}
}

// Logger Close
func (l *Logger) Close() error {
	if l.file != nil {
		if err := l.file.Close(); err != nil {
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
