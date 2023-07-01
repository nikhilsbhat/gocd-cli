package cmd

import (
	goCdLogger "github.com/nikhilsbhat/gocd-sdk-go/pkg/logger"
	"github.com/sirupsen/logrus"
)

var cliLogger *logrus.Logger

func SetLogger(logLevel string) {
	logger := logrus.New()
	logger.SetLevel(goCdLogger.GetLoglevel(logLevel))
	logger.WithField("gocd-cli", true)
	logger.SetFormatter(&logrus.JSONFormatter{})
	cliLogger = logger
}
