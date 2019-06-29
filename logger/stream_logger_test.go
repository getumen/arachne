package logger

import (
	"testing"

	"github.com/getumen/lucy"
	gomock "github.com/golang/mock/gomock"
)

func TestStreamLogger_Debugf(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		logLevel lucy.LogLevel
		callNum  int
	}{
		{
			lucy.DebugLevel,
			1,
		},
		{
			lucy.InfoLevel,
			0,
		},
		{
			lucy.WarnLevel,
			0,
		},
		{
			lucy.ErrorLevel,
			0,
		},
		{
			lucy.CriticalLevel,
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
		logLevel lucy.LogLevel
		callNum  int
	}{
		{
			lucy.DebugLevel,
			1,
		},
		{
			lucy.InfoLevel,
			1,
		},
		{
			lucy.WarnLevel,
			0,
		},
		{
			lucy.ErrorLevel,
			0,
		},
		{
			lucy.CriticalLevel,
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
		logLevel lucy.LogLevel
		callNum  int
	}{
		{
			lucy.DebugLevel,
			1,
		},
		{
			lucy.InfoLevel,
			1,
		},
		{
			lucy.WarnLevel,
			1,
		},
		{
			lucy.ErrorLevel,
			0,
		},
		{
			lucy.CriticalLevel,
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
		logLevel lucy.LogLevel
		callNum  int
	}{
		{
			lucy.DebugLevel,
			1,
		},
		{
			lucy.InfoLevel,
			1,
		},
		{
			lucy.WarnLevel,
			1,
		},
		{
			lucy.ErrorLevel,
			1,
		},
		{
			lucy.CriticalLevel,
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
		logLevel lucy.LogLevel
		callNum  int
	}{
		{
			lucy.DebugLevel,
			1,
		},
		{
			lucy.InfoLevel,
			1,
		},
		{
			lucy.WarnLevel,
			1,
		},
		{
			lucy.ErrorLevel,
			1,
		},
		{
			lucy.CriticalLevel,
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
