package logger

import "github.com/sirupsen/logrus"

type Level = logrus.Level // todo: bad

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type Logger interface {
	Logf(level Level, format string, args ...any)
	SetLevel(level Level)
	GetLevel() Level
}
