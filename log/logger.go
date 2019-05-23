package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Logger  *logrus.Logger

func init()  {
	Logger = logrus.New()
	Logger.Out = os.Stdout
	Logger.SetFormatter(&logrus.TextFormatter{DisableColors:true})
}
