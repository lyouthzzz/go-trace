package log

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	DebugLevel  Level = zapcore.DebugLevel
	InfoLevel   Level = zapcore.InfoLevel
	WarnLevel   Level = zapcore.WarnLevel
	ErrorLevel  Level = zapcore.ErrorLevel
	DPanicLevel Level = zapcore.DPanicLevel
	PanicLevel  Level = zapcore.PanicLevel
	FatalLevel  Level = zapcore.FatalLevel
)

func NewMultiWriter(level Level, writers ...io.Writer) *zap.Logger {
	ws := make([]zapcore.WriteSyncer, len(writers))
	for i, w := range writers {
		ws[i] = zapcore.AddSync(w)
	}
	mws := zapcore.NewMultiWriteSyncer(ws...)

	core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), mws, level)

	return zap.New(core)
}
