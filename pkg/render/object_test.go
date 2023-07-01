package render_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	goCdLogger "github.com/nikhilsbhat/gocd-sdk-go/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var log *logrus.Logger

//nolint:gochecknoinits
func init() {
	logger := logrus.New()
	logger.SetLevel(goCdLogger.GetLoglevel("info"))
	logger.WithField("gocd-cli", true)
	logger.SetFormatter(&logrus.JSONFormatter{})
	log = logger
}

func TestObject_CheckFileType(t *testing.T) {
	t.Run("should validate content as json", func(t *testing.T) {
		obj := render.Object(`{"name": "testing"}`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "json", actual)
	})

	t.Run("should validate content as unknown since malformed json passed", func(t *testing.T) {
		obj := render.Object(`{"name": "testing"`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "unknown", actual)
	})

	t.Run("should validate content as yaml", func(t *testing.T) {
		obj := render.Object(`---
name: "testing"`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "yaml", actual)
	})

	t.Run("should validate content as unknown since malformed yaml passed", func(t *testing.T) {
		obj := render.Object(`---
name: "testing`)

		actual := obj.CheckFileType(log)
		assert.Equal(t, "unknown", actual)
	})
}
