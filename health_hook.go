package health

import (
	"github.com/sirupsen/logrus"
)

type HealthHook struct {
	Health *Health
}

func NewHealthHook(version string) *HealthHook {
	return &HealthHook{New("1.0.0")}
}

func (hook *HealthHook) Fire(entry *logrus.Entry) error {
	switch entry.Level {
	case logrus.PanicLevel:
		hook.Health.incError()
		return nil
	case logrus.FatalLevel:
		hook.Health.incError()
		return nil
	case logrus.ErrorLevel:
		hook.Health.incError()
		return nil
	case logrus.WarnLevel:
		hook.Health.incWarning()
		return nil
	default:
		return nil
	}
}

func (hook *HealthHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}
