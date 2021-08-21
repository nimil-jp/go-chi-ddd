package log

import (
	"go-chi-ddd/config"
	"go-chi-ddd/domain/repository"
	"go.uber.org/zap"
)

func init() {
	var zConfig zap.Config
	var zapLog *zap.Logger

	if config.IsDebugging() {
		zConfig = zap.NewDevelopmentConfig()
	} else {
		zConfig = zap.NewProductionConfig()
	}

	zConfig.DisableStacktrace = true

	if config.IsDebugging() {
		zapLog, _ = zConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	} else {
		zapLog, _ = zConfig.Build(zap.AddCallerSkip(1))
	}

	log = logger{
		zapLog: zapLog,
	}
}

type logger struct {
	zapLog *zap.Logger
}

var log logger

func Logger() repository.ILogger {
	return &log
}

func (l *logger) Sync() {
	_ = l.zapLog.Sync()
}

// normal

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.zapLog.Debug(msg, fields...)
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.zapLog.Info(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.zapLog.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.zapLog.Error(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.zapLog.Fatal(msg, fields...)
}

// format

func (l *logger) Debugf(template string, args ...interface{}) {
	l.zapLog.Sugar().Debugf(template, args...)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.zapLog.Sugar().Infof(template, args...)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.zapLog.Sugar().Warnf(template, args...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.zapLog.Sugar().Errorf(template, args...)
}

func (l *logger) Fatalf(template string, args ...interface{}) {
	l.zapLog.Sugar().Fatalf(template, args...)
}
