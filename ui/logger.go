package ui

import (
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type ContextHook struct{}

func (hook ContextHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook ContextHook) Fire(entry *logrus.Entry) error {
	pc := make([]uintptr, 3, 3)
	cnt := runtime.Callers(6, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i] - 1)
		name := fu.Name()
		if !strings.Contains(name, "github.com/sirupsen/logrus") {
			file, line := fu.FileLine(pc[i] - 1)
			entry.Data["file"] = path.Base(file)
			entry.Data["func"] = path.Base(name)
			entry.Data["line"] = line
			break
		}
	}
	return nil
}

type NotificationsHook struct {
	n *Notifications
}

func NewNotificationsHook(n *Notifications) *NotificationsHook {
	return &NotificationsHook{n: n}
}

func (h NotificationsHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}

func (h NotificationsHook) Fire(entry *logrus.Entry) error {
	d := 15 * time.Second
	if entry.Level == logrus.InfoLevel {
		d = time.Second
	}

	h.n.Show(entry.Level.String(), entry.Message, d)
	return nil
}

var Logger *logrus.Entry

func init() {
	logrus.AddHook(ContextHook{})
	logrus.SetLevel(logrus.DebugLevel)
	Logger = logrus.WithFields(logrus.Fields{})

}