package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger interface {
	Debug(message string)
	Info(message string)
	Warning(message string)
	Error(message string)
	Fatal(message string)
	Format(message string, args ...interface{})
}

type logger struct {
	level LogLevel
	out   io.Writer
}

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

func NewLogger(level LogLevel, out io.Writer) Logger {
	return &logger{level: level, out: out}
}

func (l *logger) Debug(message string) {
	if l.level <= DebugLevel {
		l.log("DEBUG", message)
	}
}

func (l *logger) Info(message string) {
	if l.level <= InfoLevel {
		l.log("INFO", message)
	}
}

func (l *logger) Warning(message string) {
	if l.level <= WarningLevel {
		l.log("WARNING", message)
	}
}

func (l *logger) Error(message string) {
	if l.level <= ErrorLevel {
		pc, _, _, ok := runtime.Caller(1)
		if !ok {

		}
		fn := runtime.FuncForPC(pc)
		if fn == nil {

		}
		name := fn.Name()
		idx := strings.LastIndex(name, ".")
		if idx >= 0 {
			name = name[idx+1:]
		}
		l.log("ERROR", message)
	}
}

func (l *logger) Fatal(message string) {
	l.log("FATAL", message)
	os.Exit(1)
}

func (l *logger) Format(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.log("FORMAT", message)
}

func getFunctionName() string {
	// Get the function name from the runtime call stack
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	name := fn.Name()
	idx := strings.LastIndex(name, ".")
	if idx >= 0 {
		name = name[idx+1:]
	}
	return name
}

func (l *logger) log(level string, message string) {
	timestamp := time.Now().Format("2006-01-02T15:04:05Z07:00")
	fmt.Fprintf(l.out, "%s %s %s\n", timestamp, level, message)
	logfile, err := os.OpenFile("bot.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		fmt.Println(err)
	}
	defer logfile.Close()
	mw := io.MultiWriter(os.Stdout, logfile)

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(mw)
}
