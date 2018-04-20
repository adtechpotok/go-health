package health

import (
	"github.com/sirupsen/logrus"
)

//Hook for logrus
type HealthHook struct {
	Health *Health
}

//Create and return a new instance of the hook
func NewHealthHook(version string) *HealthHook {
	return &HealthHook{New(version)}
}

//Execute hook
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

//Return slice of logrus levels with witch hook work
func (hook *HealthHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}
}
