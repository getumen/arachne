package arachne

//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/arachne

// Logger is interface for logging in arachne crawler.
type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Criticalf(format string, v ...interface{})
}

// LogLevel : DebugLevel, InfoLevel, WarnLevel, ErrorLevel, CriticalLevel
type LogLevel int

const (
	// DebugLevel : debug or higher
	DebugLevel LogLevel = iota
	// InfoLevel : info or higher
	InfoLevel
	// WarnLevel : warn or higher
	WarnLevel
	// ErrorLevel : error or higher
	ErrorLevel
	// CriticalLevel : critical
	CriticalLevel
)
