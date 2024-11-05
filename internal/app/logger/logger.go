package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

func CreateLogger(level string) Logger {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, _ := zap.ParseAtomicLevel(level)

	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, _ := cfg.Build()
	return Logger{
		log: zl,
	}
}

func (logger Logger) GetLogger() zap.Logger {
	return *logger.log
}
