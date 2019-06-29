package lucy

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/lucy

import "log"

// Logger is interface for logging in lucy crawler.
type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Criticalf(format string, v ...interface{})
}

// StdoutLogger is a sample logger that prints all logs to stdout
type StdoutLogger struct{}

// Debugf prints args to stdout
func (StdoutLogger) Debugf(format string, v ...interface{}) { log.Printf(format, v...) }

// Infof prints args to stdout
func (StdoutLogger) Infof(format string, v ...interface{}) { log.Printf(format, v...) }

// Warnf prints args to stdout
func (StdoutLogger) Warnf(format string, v ...interface{}) { log.Printf(format, v...) }

// Errorf prints args to stdout
func (StdoutLogger) Errorf(format string, v ...interface{}) { log.Printf(format, v...) }

// Criticalf prints args to stdout
func (StdoutLogger) Criticalf(format string, v ...interface{}) { log.Printf(format, v...) }
