package logrusfluent

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/fluent/fluent-logger-golang/fluent"
)

// FluentHook to send logs via fluentd.
type FluentHook struct {
	Fluent     *fluent.Fluent
	DefaultTag string
}

// NewFluentHook creates a new hook to send to fluentd.
func NewFluentHook(config fluent.Config) (*FluentHook, error) {
	logger, err := fluent.New(config)
	if err != nil {
		return nil, err
	}
	return &FluentHook{Fluent: logger, DefaultTag: "app"}, nil
}

// Fire implements logrus.Hook interface Fire method.
func (f *FluentHook) Fire(entry *logrus.Entry) error {
	msg := f.buildMessage(entry)
	tag := f.DefaultTag
	rawTag, ok := entry.Data["tag"]
	if ok {
		tag = fmt.Sprint(rawTag)
	}
	f.Fluent.Post(tag, msg)
	return nil
}

// Levels implements logrus.Hook interface Levels method.
func (f *FluentHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (f *FluentHook) buildMessage(entry *logrus.Entry) map[string]interface{} {
	data := make(map[string]interface{})

	for k, v := range entry.Data {
		if k == "tag" {
			continue
		}
		switch v.(type) {
		case uint8, uint16, uint32, uint64, int8, int16, int32, int64, float32, float64, complex64, complex128, uint, int, string:
			// For permitive types, assign directly to preserve original type
			data[k] = v
		default:
			data[k] = fmt.Sprint(v)
		}
	}
	data["msg"] = entry.Message
	data["level"] = entry.Level.String()

	return data
}
