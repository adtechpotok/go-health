package health

import (
	"github.com/sirupsen/logrus"
)

type HealthHook struct {
	Health *Health
}

func NewHealthHook(version string) *HealthHook {
	return &HealthHook{New(version)}
}

func (hook *HealthHook) Fire(entry *logrus.Entry) error {
	switch entry.Level {
	case logrus.PanicLevel:
		hook.Health.incError()
	case logrus.FatalLevel:
		hook.Health.incError()
	case logrus.ErrorLevel:
		hook.Health.incError()
	case logrus.WarnLevel:
		hook.Health.incWarning()
	}

	return nil
}

func (hook *HealthHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}
