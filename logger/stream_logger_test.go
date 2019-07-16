package logger

import (
	"testing"

	"github.com/getumen/arachne"
	gomock "github.com/golang/mock/gomock"
)

func TestStreamLogger_Debugf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		logLevel arachne.LogLevel
		callNum  int
	}{
		{
			arachne.DebugLevel,
			1,
		},
		{
			arachne.InfoLevel,
			0,
		},
		{
			arachne.WarnLevel,
			0,
		},
		{
			arachne.ErrorLevel,
			0,
		},
		{
			arachne.CriticalLevel,
			0,
		},
	}
	for _, test := range tests {
		logger := NewMocklogger(ctrl)
		logger.EXPECT().Printf(gomock.Any()).Times(test.callNum)
		streamLogger := &StreamLogger{
			LogLevel:    test.logLevel,
			debugLogger: logger,
		}
		streamLogger.Debugf("test")
	}
}

func TestStreamLogger_Infof(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		logLevel arachne.LogLevel
		callNum  int
	}{
		{
			arachne.DebugLevel,
			1,
		},
		{
			arachne.InfoLevel,
			1,
		},
		{
			arachne.WarnLevel,
			0,
		},
		{
			arachne.ErrorLevel,
			0,
		},
		{
			arachne.CriticalLevel,
			0,
		},
	}
	for _, test := range tests {
		logger := NewMocklogger(ctrl)
		logger.EXPECT().Printf(gomock.Any()).Times(test.callNum)
		streamLogger := &StreamLogger{
			LogLevel:   test.logLevel,
			infoLogger: logger,
		}
		streamLogger.Infof("test")
	}
}

func TestStreamLogger_Warnf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		logLevel arachne.LogLevel
		callNum  int
	}{
		{
			arachne.DebugLevel,
			1,
		},
		{
			arachne.InfoLevel,
			1,
		},
		{
			arachne.WarnLevel,
			1,
		},
		{
			arachne.ErrorLevel,
			0,
		},
		{
			arachne.CriticalLevel,
			0,
		},
	}
	for _, test := range tests {
		logger := NewMocklogger(ctrl)
		logger.EXPECT().Printf(gomock.Any()).Times(test.callNum)
		streamLogger := &StreamLogger{
			LogLevel:   test.logLevel,
			warnLogger: logger,
		}
		streamLogger.Warnf("test")
	}
}

func TestStreamLogger_Errorf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		logLevel arachne.LogLevel
		callNum  int
	}{
		{
			arachne.DebugLevel,
			1,
		},
		{
			arachne.InfoLevel,
			1,
		},
		{
			arachne.WarnLevel,
			1,
		},
		{
			arachne.ErrorLevel,
			1,
		},
		{
			arachne.CriticalLevel,
			0,
		},
	}
	for _, test := range tests {
		logger := NewMocklogger(ctrl)
		logger.EXPECT().Printf(gomock.Any()).Times(test.callNum)
		streamLogger := &StreamLogger{
			LogLevel:    test.logLevel,
			errorLogger: logger,
		}
		streamLogger.Errorf("test")
	}
}

func TestStreamLogger_Criticalf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		logLevel arachne.LogLevel
		callNum  int
	}{
		{
			arachne.DebugLevel,
			1,
		},
		{
			arachne.InfoLevel,
			1,
		},
		{
			arachne.WarnLevel,
			1,
		},
		{
			arachne.ErrorLevel,
			1,
		},
		{
			arachne.CriticalLevel,
			1,
		},
	}
	for _, test := range tests {
		logger := NewMocklogger(ctrl)
		logger.EXPECT().Printf(gomock.Any()).Times(test.callNum)
		streamLogger := &StreamLogger{
			LogLevel:       test.logLevel,
			criticalLogger: logger,
		}
		streamLogger.Criticalf("test")
	}
}
