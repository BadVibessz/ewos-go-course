package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	Logf(level logrus.Level, format string, args ...any)
}
