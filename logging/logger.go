package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"time"
)

type Logger struct {
	out  io.Writer
	file *os.File
	log  *log.Logger
}

func NewLogger() (Logger, error) {
	var file *os.File
	var err error
	var out io.Writer
	file, err = os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	out = os.Stdout

	l := &Logger{
		out:  out,
		file: file,
	}
	l.log = log.New(io.MultiWriter(out, file), "", 1)
	return *l, nil
}

func (l *Logger) Errorf(format string, err error, v ...interface{}) {
	l.logf(format, err, v...)
}

func (l *Logger) Error(err error, v ...interface{}) {
	l.logf("", err, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(format, nil, v...)
}

func (l *Logger) logf(format string, err error, v ...interface{}) {

	if err != nil {
		// Add log message to the logger as it would be printf
		msg := fmt.Sprintf(format, v...)

		// Log to stdout
		funcName, file, line := getCaller()
		timestamp := time.Now().Format("15:04:05.000")
		logMsg := fmt.Sprintf("%s: \nMessage: [%s] \n\t%s \n\t%s \n\t%s in line: %d", timestamp, err, msg, funcName, file, line)
		l.log.Println(logMsg)

		// Log to file
		if l.file != nil {
			fmt.Fprintln(l.file, logMsg)
		}
	}
	// Add log message to the logger as it would be printf
	msg := fmt.Sprintf(format, v...)

	// Log to stdout
	funcName, file, line := getCaller()
	timestamp := time.Now().Format("15:04:05.000")
	logMsg := fmt.Sprintf("%s: \nMessage: [%s]  \n\t%s \n\t\t%s in line: %d", timestamp, msg, funcName, file, line)
	l.log.Println(logMsg)

	// Log to file
	if l.file != nil {
		fmt.Fprintln(l.file, logMsg)
	}

}

func (l *Logger) Close() error {

	if l.file != nil {
		if err := l.file.Close(); err != nil {
			return err
		}
	}

	return nil
}

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
