package log

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	Driver zerolog.Logger
	Rotate *lumberjack.Logger
}

func NewRotate(levelStr string, filename string, maxSize, maxBackups int) zerolog.Logger {
	level, _ := zerolog.ParseLevel(levelStr)
	return zerolog.New(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    10, //默认1MB
		MaxBackups: 10,
		LocalTime:  true,
	}).Level(level).With().Timestamp().Caller().Logger()
}
