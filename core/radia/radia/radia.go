package radia

import (
	"fmt"
	"sync"

	"github.com/piotrwyrw/radia/internal/rconst"
	"github.com/piotrwyrw/radia/radia/rregistry"
	"github.com/sirupsen/logrus"
)

var initOnce sync.Once = sync.Once{}

func Initialize() {
	initOnce.Do(func() {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})

		fmt.Println(rconst.RadiaBanner)
		fmt.Printf("\t\t( Radia Renderer %s )\n\n", rconst.RadiaVersion)

		logrus.Infof("Radia Initialized")
		logrus.SetLevel(logrus.DebugLevel)

		_ = rregistry.GetCentralRegistry()
	})
}
