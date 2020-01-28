package service

import (
	"flag"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logLevelCfg = "log-level"
)

var (
	// Command line pointer to logger level flag configuration.
	loggerLevelPtr *string
)

func loggerFlags(flags *flag.FlagSet) {
	loggerLevelPtr = flags.String(logLevelCfg, "INFO",
		"Output level of logs (DEBUG, INFO, WARN, ERROR, DPANIC, PANIC, FATAL)")
}

func newLogger() (*zap.Logger, error) {
	var level zapcore.Level
	err := (&level).UnmarshalText([]byte(*loggerLevelPtr))
	if err != nil {
		return nil, err
	}
	conf := zap.NewProductionConfig()
	conf.Level.SetLevel(level)
	return conf.Build()
}
