package bot

import (
	"fmt"
	"log/slog"
)

type logger struct {
	log *slog.Logger
}

func (l *logger) Debugf(format string, args ...any) {
	l.log.Debug(fmt.Sprintf(format, args...))
}

func (l *logger) Errorf(format string, args ...any) {
	l.log.Error(fmt.Sprintf(format, args...))
}
