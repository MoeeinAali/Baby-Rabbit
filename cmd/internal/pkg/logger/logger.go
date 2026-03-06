package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func Init() {
	logger, _ := zap.NewProduction()
	Log = logger.Sugar()
}

func Sync() {
	Log.Sync()
}
