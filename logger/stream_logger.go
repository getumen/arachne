package logger

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/arachne/logger

import (
	"log"
	"os"

	"github.com/getumen/arachne"
)

type logger interface {
	Printf(format string, v ...interface{})
}

// StreamLogger is a sample logger that prints all logs to stdout
type StreamLogger struct {
	LogLevel       arachne.LogLevel
	debugLogger    logger
	infoLogger     logger
	warnLogger     logger
	errorLogger    logger
	criticalLogger logger
}

// NewStdoutLogger : standard stream logger
func NewStdoutLogger(logLevel arachne.LogLevel) arachne.Logger {
	return &StreamLogger{
		LogLevel:       logLevel,
		debugLogger:    log.New(os.Stderr, "[DEBUG]\t", log.LstdFlags),
		infoLogger:     log.New(os.Stderr, "[Info]\t", log.LstdFlags),
		warnLogger:     log.New(os.Stderr, "[Warn]\t", log.LstdFlags),
		errorLogger:    log.New(os.Stderr, "[Error]\t", log.LstdFlags),
		criticalLogger: log.New(os.Stderr, "[Critical]\t", log.LstdFlags),
	}
}

// Debugf prints args to stdout
func (s *StreamLogger) Debugf(format string, v ...interface{}) {
	if s.LogLevel <= arachne.DebugLevel {
		s.debugLogger.Printf(format, v...)
	}
}

// Infof prints args to stdout
func (s *StreamLogger) Infof(format string, v ...interface{}) {
	if s.LogLevel <= arachne.InfoLevel {
		s.infoLogger.Printf(format, v...)
	}
}

// Warnf prints args to stdout
func (s *StreamLogger) Warnf(format string, v ...interface{}) {
	if s.LogLevel <= arachne.WarnLevel {
		s.warnLogger.Printf(format, v...)
	}
}

// Errorf prints args to stdout
func (s *StreamLogger) Errorf(format string, v ...interface{}) {
	if s.LogLevel <= arachne.ErrorLevel {
		s.errorLogger.Printf(format, v...)
	}
}

// Criticalf prints args to stdout
func (s *StreamLogger) Criticalf(format string, v ...interface{}) {
	if s.LogLevel <= arachne.CriticalLevel {
		s.criticalLogger.Printf(format, v...)
	}
}
