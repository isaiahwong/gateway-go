package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is a generic type interface
type Logger interface {
	// Trace()
	// Debug()

	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panicln(args ...interface{})
}

type logger struct {
	*logrus.Logger
}

// NewLogger returns new Logger.
func NewLogger() *logrus.Logger {
	log := logrus.New()

	logrus.Trace()

	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(logrus.InfoLevel)
	// log.SetReportCaller(true)

	return log
}
