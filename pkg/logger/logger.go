package logger

import (
	colorable "github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(env string) {
	var err error
	if env == "production" {
		Log, err = zap.NewProduction()
	} else {
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.AddSync(colorable.NewColorableStdout()),
			zapcore.DebugLevel,
		)

		Log = zap.New(core,
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	}
	if err != nil {
		panic(err)
	}
}
