package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global       *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
)

func init() {
	SetLogger(New(defaultLevel))
}

// New создает экземпляр *zap.SugaredLogger со стандартным json выводом.
// Если уровень логгирования не передан - будет использоваться уровень
// по умолчанию (zap.ErrorLevel)
func New(level zapcore.LevelEnabler, options ...zap.Option) *zap.SugaredLogger {
	return NewWithSink(level, os.Stdout, options...)
}

func NewWithSink(level zapcore.LevelEnabler, sink io.Writer, options ...zap.Option) *zap.SugaredLogger {
	if level == nil {
		level = defaultLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "ts",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			SkipLineEnding: false,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(sink),
		level)

	return zap.New(core, options...).Sugar() // wrapCoreWithMetrics(core)
}

// SetLogger устанавливает глобальный логгер. Функция непотокобезопасна.
func SetLogger(l *zap.SugaredLogger) {
	global = l
}

// SetLevel устанавливает уровень логгирования глобального логгера.
func SetLevel(level zapcore.Level) {
	defaultLevel.SetLevel(level)
}
