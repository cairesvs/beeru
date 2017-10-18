// Wrapper of logrus logging infrastructure.
//
// It adds common fields as application, environment and version and provides a
// simple API:
//
//   logger.Info("Lorem Ipsum")
//   // {"application":"urls","environment":"prod","fields.level":6,"level":6,"level_name":"info","full_message":"Lorem Ipsum","time":"2017-06-30T16:10:50-03:00","timestamp":1498849850178,"version":"detached"}
//
//   logger.WithFields(logger.Fields{"custom": "field"}).Info("Lorem Ipsum")
//   // {"application":"urls","custom":"field","environment":"prod","fields.level":6,"level":6,"level_name":"info","full_message":"Lorem Ipsum","time":"2017-06-30T16:10:50-03:00","timestamp":1498849850178,"version":"detached"}
//
//   logger.Infof("Lorem %", "Ipsum")
//   // {"application":"urls","environment":"prod","fields.level":6,"level":6,"level_name":"info","full_message":"Lorem Ipsum","time":"2017-06-30T16:10:50-03:00","timestamp":1498849850178,"version":"detached"}
//
package logger

import (
	"fmt"
	"log/syslog"
	"time"

	"github.com/cairesvs/beeru/pkg/config"
	log "github.com/sirupsen/logrus"
)

type Fields log.Fields

type entry struct {
	*log.Entry
}

func init() {
	log.SetLevel(levelFromString(config.Get("logLevel")))

	// For better developer experience log in plain text instead JSON
	// in development
	if config.Get("environment") != string(config.DevelopmentEnvironment) {
		log.SetFormatter(&log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyMsg:   "full_message",
				log.FieldKeyLevel: "level_name",
			},
		})
	}

	log.RegisterExitHandler(func() {
		Info("application will stop :(")
	})
}

func WithFields(fields Fields) *entry {
	fields["application"] = config.Get("application")
	fields["environment"] = config.Get("environment")
	fields["version"] = config.Get("version")

	return &entry{log.WithFields(log.Fields(fields))}
}

func Debug(msg interface{}) {
	WithFields(Fields{}).Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	WithFields(Fields{}).Debug(fmt.Sprintf(format, args...))
}

func Info(msg interface{}) {
	WithFields(Fields{}).Info(msg)
}

func Infof(format string, args ...interface{}) {
	WithFields(Fields{}).Info(fmt.Sprintf(format, args...))
}

func Warn(msg interface{}) {
	WithFields(Fields{}).Warn(msg)
}

func Warnf(format string, args ...interface{}) {
	WithFields(Fields{}).Warn(fmt.Sprintf(format, args...))
}

func Error(msg interface{}) {
	WithFields(Fields{}).Error(msg)
}

func Errorf(format string, args ...interface{}) {
	WithFields(Fields{}).Error(fmt.Sprintf(format, args...))
}

func Fatal(msg interface{}) {
	WithFields(Fields{}).Fatal(msg)
}

func Fatalf(format string, args ...interface{}) {
	WithFields(Fields{}).Fatal(fmt.Sprintf(format, args...))
}

func Panic(msg interface{}) {
	WithFields(Fields{}).Panic(msg)
}

func Panicf(format string, args ...interface{}) {
	WithFields(Fields{}).Panic(fmt.Sprintf(format, args...))
}

func (entry *entry) Debug(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_DEBUG).
		WithField("timestamp", nowMillis()).
		Debug(msg)
}

func (entry *entry) Info(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_INFO).
		WithField("timestamp", nowMillis()).
		Info(msg)
}

func (entry *entry) Warn(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_WARNING).
		WithField("timestamp", nowMillis()).
		Warn(msg)
}

func (entry *entry) Error(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_ERR).
		WithField("timestamp", nowMillis()).
		Error(msg)
}

func (entry *entry) Fatal(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_CRIT).
		WithField("timestamp", nowMillis()).
		Fatal(msg)
}

func (entry *entry) Panic(msg interface{}) {
	entry.Entry.
		WithField("level", syslog.LOG_EMERG).
		WithField("timestamp", nowMillis()).
		Panic(msg)
}

// Convert the level string to a logrus Level. E.g. "panic" becomes "PanicLevel".
func levelFromString(level string) log.Level {
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "warning":
		return log.WarnLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	}

	return log.DebugLevel
}

// Unix timestamp in milliseconds resolution
func nowMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
