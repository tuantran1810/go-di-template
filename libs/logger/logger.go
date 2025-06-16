package logger

import "go.uber.org/zap"

func MustNamedLogger(name string) *zap.SugaredLogger {
	logger, err := zap.NewProduction(
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		panic(err)
	}

	return logger.Named(name).Sugar()
}
