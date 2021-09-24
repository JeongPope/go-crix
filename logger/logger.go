package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota + 1
	INFO
	WARN
	ERROR
	FATAL
	PANIC
)

type Logger struct {
	*log.Logger
	level Level
}

func init() {
	logLevel := os.Getenv("GOCRIX_LOG_LEVEL")
	var l Level
	switch logLevel {
	case "debug", "DEBUG":
		l = DEBUG
	case "info", "INFO":
		l = INFO
	case "warn", "WARN":
		l = WARN
	case "error", "ERROR":
		l = ERROR
	case "fatal", "FATAL":
		l = FATAL
	case "panic", "PANIC":
		l = PANIC
	default:
		l = ERROR
	}
	SetLevel(l)

	logFileName := os.Getenv("GOCRIX_LOG_FILE")
	if logFileName != "" {
		f, err := os.OpenFile(logFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
		if err == nil {
			SetOut(f)
		} else {
			Warn("log file not open ??? ")
			Error(err.Error())
		}
	}
}

var Log = NewLogger()

func SetOut(out io.Writer) {
	Log.SetOut(out)
}

func SetLevel(level Level) {
	Log.SetLevel(level)
}

func Debug(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Println(s)
	Log.output(DEBUG, "[DEBUG]", fmt.Sprint(args...))
}

func Debugf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println(s)
	Log.output(DEBUG, "[DEBUG]", fmt.Sprintf(format, args...))
}

func Info(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Println(s)
	Log.output(INFO, "[INFO]", fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println(s)
	Log.output(INFO, "[INFO]", fmt.Sprintf(format, args...))
}

func Warn(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Println(s)
	Log.output(WARN, "[WARNING]", fmt.Sprint(args...))
}

func Warnf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println(s)
	Log.output(WARN, "[WARNING]", fmt.Sprintf(format, args...))
}

func Error(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Println(s)
	Log.output(ERROR, "[ERROR]", fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println(s)
	Log.output(ERROR, "[ERROR]", fmt.Sprintf(format, args...))
}

func Fatal(args ...interface{}) {
	if Log.level <= FATAL {
		s := fmt.Sprint(args...)
		log.Println(s)
		Log.output(FATAL, "[FATAL]", fmt.Sprint(args...))
		os.Exit(1)
	}
}

func Fatalf(format string, args ...interface{}) {
	if Log.level <= FATAL {
		s := fmt.Sprintf(format, args...)
		log.Println(s)
		Log.output(FATAL, "[FATAL]", fmt.Sprintf(format, args...))

		os.Exit(1)
	}
}

func Panic(args ...interface{}) {
	if Log.level <= PANIC {
		s := fmt.Sprint(args...)
		log.Println(s)
		Log.output(PANIC, "[PANIC]", fmt.Sprint(args...))
		panic("")
	}
}

func Panicf(format string, args ...interface{}) {
	if Log.level <= PANIC {
		s := fmt.Sprintf(format, args...)
		log.Println(s)
		Log.output(PANIC, "[PANIC]", fmt.Sprintf(format, args...))

		panic("")
	}
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
		level:  INFO,
	}
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) SetOut(out io.Writer) {
	l.Logger.SetOutput(out)
}

func (l *Logger) output(le Level, prefix string, log string) {
	if l.level <= le {
		l.Output(3, fmt.Sprintf("%s %s", prefix, log))
	}
}

func (l *Logger) Debug(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Println(s)
	l.output(DEBUG, "[DEBUG]", fmt.Sprint(args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println(s)
	l.output(DEBUG, "[DEBUG]", fmt.Sprintf(format, args...))
}

func (l *Logger) Info(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Println(s)
	l.output(INFO, "[INFO]", fmt.Sprint(args...))
}

func (l *Logger) Infof(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println(s)
	l.output(INFO, "[INFO]", fmt.Sprintf(format, args...))
}

func (l *Logger) Warn(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Println(s)
	l.output(WARN, "[WARNING]", fmt.Sprint(args...))
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println(s)
	l.output(WARN, "[WARNING]", fmt.Sprintf(format, args...))
}

func (l *Logger) Error(args ...interface{}) {
	s := fmt.Sprint(args...)
	log.Println(s)
	l.output(ERROR, "[ERROR]", fmt.Sprint(args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	log.Println(s)
	l.output(ERROR, "[ERROR]", fmt.Sprintf(format, args...))
}

func (l *Logger) Fatal(args ...interface{}) {
	if l.level <= FATAL {
		s := fmt.Sprint(args...)
		log.Println(s)
		l.output(FATAL, "[FATAL]", fmt.Sprint(args...))
		os.Exit(1)
	}
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	if l.level <= FATAL {
		s := fmt.Sprintf(format, args...)
		log.Println(s)
		l.output(FATAL, "[FATAL]", fmt.Sprintf(format, args...))
		os.Exit(1)
	}
}

func (l *Logger) Panic(args ...interface{}) {
	if l.level <= PANIC {
		s := fmt.Sprint(args...)
		log.Println(s)
		l.output(PANIC, "[PANIC]", s)
		panic(s)
	}
}

func (l *Logger) Panicf(format string, args ...interface{}) {
	if l.level <= PANIC {
		s := fmt.Sprintf(format, args...)

		log.Println(s)
		l.output(PANIC, "[PANIC]", s)
		panic(s)
	}
}
