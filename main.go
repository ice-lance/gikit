package main

import (
	"ice_tools/pkg/log"

	"go.uber.org/zap/zapcore"
)

func main() {
	log.InitLog("debug", []zapcore.Field{})
	log.Debug("adfasdfasdfasdfa")
}
