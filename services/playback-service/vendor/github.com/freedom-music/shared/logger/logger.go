package logger

import (
    "strings"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

// New builds a zap logger with sane defaults for JSON structured logging.
func New(level string) (*zap.Logger, error) {
    lvl := zap.InfoLevel
    if level != "" {
        if err := lvl.UnmarshalText([]byte(strings.ToLower(level))); err != nil {
            return nil, err
        }
    }

    cfg := zap.Config{
        Level:       zap.NewAtomicLevelAt(lvl),
        Development: false,
        Encoding:    "json",
        EncoderConfig: zapcore.EncoderConfig{
            TimeKey:        "ts",
            LevelKey:       "level",
            NameKey:        "logger",
            CallerKey:      "caller",
            MessageKey:     "msg",
            StacktraceKey:  "stacktrace",
            EncodeLevel:    zapcore.LowercaseLevelEncoder,
            EncodeTime:     zapcore.RFC3339TimeEncoder,
            EncodeCaller:   zapcore.ShortCallerEncoder,
            LineEnding:     zapcore.DefaultLineEnding,
            EncodeDuration: zapcore.SecondsDurationEncoder,
        },
        OutputPaths:      []string{"stdout"},
        ErrorOutputPaths: []string{"stderr"},
    }

    return cfg.Build()
}
