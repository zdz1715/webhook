package util

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
)

type Logger struct {
	Driver zerolog.Logger
	Rotate *lumberjack.Logger
}

func NewRotate(filename string, maxSize, maxBackups int) io.Writer {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize, //默认1MB
		MaxBackups: maxBackups,
		LocalTime:  true,
	}
}

func NewLogRotate(filename string, maxSize, maxBackups int) zerolog.Logger {
	return zerolog.New(NewRotate(filename, maxSize, maxBackups)).With().Timestamp().Logger()
}

func NewMultiLevelWriter(writer ...io.Writer) zerolog.Logger {
	multi := zerolog.MultiLevelWriter(writer...)
	return zerolog.New(multi).With().Timestamp().Logger()
}
