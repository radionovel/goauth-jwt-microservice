package logger

import "go.uber.org/zap"

func NewConsoleLogger() *zap.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	return l
}
