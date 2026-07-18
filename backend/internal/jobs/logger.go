package jobs

import "go.uber.org/zap"

// zapLogger adapts *zap.Logger to the asynq.Logger interface so worker internals
// log through the same structured pipeline as the rest of the application.
type zapLogger struct {
	sugar *zap.SugaredLogger
}

func newZapLogger(logger *zap.Logger) *zapLogger {
	return &zapLogger{sugar: logger.Sugar()}
}

func (l *zapLogger) Debug(args ...interface{}) { l.sugar.Debug(args...) }
func (l *zapLogger) Info(args ...interface{})  { l.sugar.Info(args...) }
func (l *zapLogger) Warn(args ...interface{})  { l.sugar.Warn(args...) }
func (l *zapLogger) Error(args ...interface{}) { l.sugar.Error(args...) }
func (l *zapLogger) Fatal(args ...interface{}) { l.sugar.Fatal(args...) }
