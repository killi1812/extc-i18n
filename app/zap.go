package app

import (
	"go.uber.org/zap"
)

func devLoggerSetup() error {
	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.PanicLevel))
	if err != nil {
		return err
	}
	_ = zap.ReplaceGlobals(logger)

	return nil
}
