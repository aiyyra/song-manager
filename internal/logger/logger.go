// Package logger to keep all logs logic
package logger

import "go.uber.org/zap"

var Log *zap.Logger

func Init() {
	var err error
	Log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}
