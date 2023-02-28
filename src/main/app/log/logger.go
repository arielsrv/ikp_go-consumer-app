package log

import (
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&nested.Formatter{
		FieldsOrder:     []string{"component", "category"},
		TimestampFormat: "2006-01-02 15:04:05",
		HideKeys:        true,
		TrimMessages:    true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

type ILogger interface {
	Info(v ...any)
	Infof(format string, v ...any)
	Debugf(format string, v ...any)
	Warnf(format string, v ...any)
	Error(v ...any)
	Errorf(format string, v ...any)
	Fatal(v ...any)
}

var logger = &stdLogger{}

type stdLogger struct {
}

func (s stdLogger) Info(v ...any) {
	log.Println(v...)
}

func (s stdLogger) Infof(format string, v ...any) {
	log.Printf(format, v...)
}

func (s stdLogger) Debugf(format string, v ...any) {
	log.Debugf(format, v...)
}

func (s stdLogger) Warn(v ...any) {
	log.Warn(v...)
}

func (s stdLogger) Warnf(format string, v ...any) {
	log.Warnf(format, v...)
}

func (s stdLogger) Error(v ...any) {
	log.Error(v...)
}

func (s stdLogger) Errorf(format string, v ...any) {
	log.Errorf(format, v...)
}

func (s stdLogger) Fatal(v ...any) {
	log.Fatal(v...)
}

func (s stdLogger) SetLogLevel(value string) {
	level, err := log.ParseLevel(value)
	if err != nil {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
}

func Info(v ...any) {
	logger.Info(v...)
}

func Infof(format string, v ...any) {
	logger.Infof(format, v...)
}

func Debugf(format string, v ...any) {
	logger.Debugf(format, v...)
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

func SetLogLevel(level string) {
	logger.SetLogLevel(level)
}
