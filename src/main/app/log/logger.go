package log

import (
	"io"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&nested.Formatter{
		FieldsOrder:      []string{"component", "category"},
		TimestampFormat:  "2006-01-02 15:04:05",
		HideKeys:         true,
		NoUppercaseLevel: true,
		TrimMessages:     true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
}

type ILogger interface {
	Info(v ...any)
	Infof(format string, v ...any)
	Debugf(format string, v ...any)
	Warnf(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	Fatal(v ...any)
	GetWriter() *io.PipeWriter
}

func Info(v ...any) {
	logger.Info(v...)
}

func Infof(format string, v ...any) {
	logger.Infof(format, v...)
}

func Warnf(format string, v ...any) {
	logger.Warnf(format, v...)
}

func Error(v ...any) {
	logger.Error(v...)
}

func Errorf(format string, v ...any) {
	logger.Errorf(format, v...)
}

func Fatal(v ...any) {
	logger.Fatal(v...)
}

func SetLogLevel(value string) {
	level, err := logrus.ParseLevel(value)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(level)
	}
}

func GetWriter() *io.PipeWriter {
	return logger.Writer()
}
