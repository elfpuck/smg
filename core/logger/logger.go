package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Level uint32

var config *Config

func ParseLevel(level string) Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN":
		return WarnLevel
	case "Error":
		return ErrorLevel
	case "FATAL":
		return FatalLevel
	default:
		return DebugLevel
	}
}
func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "Error"
	case FatalLevel:
		return "FATAL"
	}
	return ""
}

const (
	FatalLevel Level = iota
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

func Init(cfg *Config) {
	config = cfg
	if config.Output == nil {
		config.Output = os.Stdout
	}
	if config.TimeFormat == "" {
		config.TimeFormat = "2006-01-02 15:04:05"
	}
}

func print(l Level, msg ...any) {
	if config.Level >= l {
		fmt.Fprintf(config.Output, "%s [%s] %v\n", time.Now().Format(config.TimeFormat), l, fmt.Sprint(msg...))
	}
}

func Println(msg ...any) {
	fmt.Fprintln(config.Output, msg...)
}

func Debug(msg ...any) {
	print(DebugLevel, msg...)
}
func Info(msg ...any) {
	print(InfoLevel, msg...)
}
func Warn(msg ...any) {
	print(WarnLevel, msg...)
}
func Error(msg ...any) {
	print(ErrorLevel, msg...)
}
func Fatal(msg ...any) {
	print(FatalLevel, msg...)
	os.Exit(1)
}

func CommonInfo(msg ...any) {
	print(InfoLevel, msg...)
	fmt.Fprintln(os.Stdout, msg...)
}

func CommonFatal(msg ...any) {
	print(FatalLevel, msg...)
	fmt.Fprintln(os.Stdout, msg...)
	os.Exit(1)
}

type Config struct {
	Level      Level
	TimeFormat string
	Output     io.Writer
}
