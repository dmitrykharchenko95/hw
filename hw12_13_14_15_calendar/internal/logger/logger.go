package logger

import (
	"github.com/sirupsen/logrus"
)

func New(level string) (*logrus.Logger, error) {
	l := logrus.New()

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	l.SetLevel(lvl)

	l.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return l, nil
}
