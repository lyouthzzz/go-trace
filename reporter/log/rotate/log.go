package rotate

import (
	"io"

	"github.com/natefinch/lumberjack"
)

type Option func(*Config)

func MaxSizeOption(maxSize int) Option {
	return func(conf *Config) {
		conf.MaxSize = maxSize
	}
}

func MaxAgeOption(maxAge int) Option {
	return func(conf *Config) {
		conf.MaxAge = maxAge
	}
}

func MaxBackupsOption(maxBackups int) Option {
	return func(conf *Config) {
		conf.MaxBackups = maxBackups
	}
}

type Config struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

func NewWriter(filename string, opts ...Option) io.Writer {
	config := &Config{
		Filename:   filename,
		MaxSize:    100,
		MaxAge:     30,
		MaxBackups: 0,
		Compress:   false,
	}

	for _, opt := range opts {
		opt(config)
	}

	return &lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
}
