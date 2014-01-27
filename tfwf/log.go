package tfwf

import (
	"io"
	"log"
	"os"
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
)

var Level = []string{"DEBUG", "INFO", "WARNING", "ERROR"}

type Log struct {
	log   *log.Logger
	level int
}

var Logger *Log

func init() {
	Logger = NewLogger(os.Stdout, "", log.LstdFlags, INFO)
}

func NewLogger(w io.Writer, prefix string, i int, level int) *Log {
	return &Log{
		log.New(w, prefix, i),
		level,
	}
}

func SetLogLevel(level int) {
	Logger.level = level
}

func (l *Log) printLog(level int, msg string) {
	if level > l.level {
		l.log.Printf("[%s] %s", Level[level], msg)
	}

}

func (l *Log) Error(msg string) {
	l.printLog(ERROR, msg)
}
func (l *Log) Warning(msg string) {
	l.printLog(WARNING, msg)
}
func (l *Log) Info(msg string) {
	l.printLog(INFO, msg)
}
func (l *Log) Debug(msg string) {
	l.printLog(DEBUG, msg)
}
