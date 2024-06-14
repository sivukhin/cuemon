package lib

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	Logger = zap.Must(zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		Encoding:          "console",
		EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: true,
	}.Build()).Sugar()
}
