package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init(env string) error {

	var cfg zap.Config

	if env == "prod" {
		cfg = zap.NewProductionConfig()
		// настройка уровня логирования
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		// формат времени в логах обычный-читаемый, также можно
		// установить ключ для поля time - cfg.EncoderConfig.Timekey = "ts"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.OutputPaths = []string{"stdout"}
		l, err := cfg.Build(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
		if err != nil {
			return err
		}
		Logger = l
	} else {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.OutputPaths = []string{"stdout"}
		l, err := cfg.Build(zap.AddCaller())
		if err != nil {
			return err
		}
		Logger = l
	}

	cfg.EncoderConfig.MessageKey = "message"
	cfg.DisableStacktrace = false

	return nil
}

func Sugar() *zap.SugaredLogger {
	return Logger.Sugar()
}

func Close() {
	_ = Logger.Sync()
}
