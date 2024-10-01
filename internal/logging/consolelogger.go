package logging

import (
	"fmt"
	"io"
)

type consoleLogger struct {
	writer    io.Writer
	formatter Formatter
}

func (e *consoleLogger) Debug(i interface{}) {
	e.write(fmt.Sprintf("%v", i), debugLevel)
}

func (e *consoleLogger) Debugf(f string, v ...interface{}) {
	e.write(fmt.Sprintf(f, v...), debugLevel)
}

func (e *consoleLogger) Info(i interface{}) {
	e.write(fmt.Sprintf("%v", i), infoLevel)
}

func (e *consoleLogger) Infof(f string, v ...interface{}) {
	e.write(fmt.Sprintf(f, v...), infoLevel)
}

func (e *consoleLogger) Error(i interface{}) {
	e.write(fmt.Sprintf("%v", i), errorLevel)
}

func (e *consoleLogger) Errorf(f string, v ...interface{}) {
	e.write(fmt.Sprintf(f, v...), errorLevel)
}

func (e *consoleLogger) write(input string, level logLevel) {
	formattedOutput := e.formatter.Format(input, level)
	_, _ = fmt.Fprintln(e.writer, formattedOutput)
}
