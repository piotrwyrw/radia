package radia

import (
	"sync"

	"github.com/sirupsen/logrus"
)

var initOnce sync.Once = sync.Once{}

func Initialize() {
	initOnce.Do(func() {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
		logrus.Infof("Radia Initialized")
		logrus.SetLevel(logrus.DebugLevel)
	})
}
