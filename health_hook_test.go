package health

import (
	"testing"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestFire(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		entry    *logrus.Entry
		warnings uint64
		errors   uint64
	}{
		{&logrus.Entry{Level: logrus.PanicLevel}, 0, 1},
		{&logrus.Entry{Level: logrus.FatalLevel}, 0, 1},
		{&logrus.Entry{Level: logrus.ErrorLevel}, 0, 1},
		{&logrus.Entry{Level: logrus.WarnLevel}, 1, 0},
		{&logrus.Entry{Level: logrus.InfoLevel}, 0, 0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("entry:%v, warnings:%d, errors:%d", tc.entry, tc.warnings, tc.errors), func(t *testing.T) {
			hook := NewHealthHook("1.0.0")
			hook.Fire(tc.entry)
			assert.Equal(tc.warnings, hook.Health.Warnings)
			assert.Equal(tc.errors, hook.Health.Errors)

		})
	}
}

func TestLevels(t *testing.T) {
	hook := NewHealthHook("1.0.0")
	assert := assert.New(t)
	assert.Equal([]logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	}, hook.Levels())
}
