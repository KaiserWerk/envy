package logging

import (
	"fmt"
	"io"
	"sync"
)

type (
	Logger interface {
		Debug(i interface{})
		Debugf(f string, v ...interface{})
		Info(i interface{})
		Infof(f string, v ...interface{})
		Error(i interface{})
		Errorf(f string, v ...interface{})
	}
	Formatter interface {
		Format(input string, level logLevel) string
	}
	logLevel   string
	envyLogger struct {
		mut       *sync.Mutex
		writer    io.Writer
		formatter Formatter
	}
	envyFormatter struct{}
)

var (
	_ Logger    = &envyLogger{}
	_ Formatter = &envyFormatter{}
)

func (e *envyFormatter) Format(input string, level logLevel) string {
	return string(level) + input
}

const (
	debugLevel logLevel = "[DEBUG] "
	infoLevel  logLevel = "[ INFO] "
	errorLevel logLevel = "[ERROR] "
)

func (e *envyLogger) write(input string, level logLevel) {
	formattedOutput := e.formatter.Format(input, level)
	fmt.Fprintln(e.writer, formattedOutput)
}

func (e *envyLogger) Debug(i interface{}) {
	e.mut.Lock()
	defer e.mut.Unlock()
	e.write(fmt.Sprintf("%v", i), debugLevel)
}

func (e *envyLogger) Debugf(f string, v ...interface{}) {
	e.mut.Lock()
	defer e.mut.Unlock()
	e.write(fmt.Sprintf(f, v...), debugLevel)
}

func (e *envyLogger) Info(i interface{}) {
	e.mut.Lock()
	defer e.mut.Unlock()
	e.write(fmt.Sprintf("%v", i), infoLevel)
}

func (e *envyLogger) Infof(f string, v ...interface{}) {
	e.mut.Lock()
	defer e.mut.Unlock()
	e.write(fmt.Sprintf(f, v...), infoLevel)
}

func (e *envyLogger) Error(i interface{}) {
	e.mut.Lock()
	defer e.mut.Unlock()
	e.write(fmt.Sprintf("%v", i), errorLevel)
}

func (e *envyLogger) Errorf(f string, v ...interface{}) {
	e.mut.Lock()
	defer e.mut.Unlock()
	e.write(fmt.Sprintf(f, v...), errorLevel)
}

func New(w io.Writer) Logger {
	return &envyLogger{
		mut:       new(sync.Mutex),
		writer:    w,
		formatter: &envyFormatter{},
	}
}
