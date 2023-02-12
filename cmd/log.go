package cmd

import (
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/sirupsen/logrus"
)

var cliLogger *logrus.Logger

func SetLogger(logLevel string) {
	logger := logrus.New()
	logger.SetLevel(gocd.GetLoglevel(logLevel))
	logger.WithField("gocd-cli", true)
	logger.SetFormatter(&logrus.JSONFormatter{})
	cliLogger = logger
}
