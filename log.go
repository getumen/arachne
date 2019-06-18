//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE -self_package=github.com/getumen/lucy
package lucy

type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Criticalf(format string, v ...interface{})
}
