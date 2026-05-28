package logger

import "go.uber.org/zap"

// Log is the package-global sugared logger. Init must be called once
// at start-up.
var Log *zap.SugaredLogger

func Init() {
	l, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	Log = l.Sugar()
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
