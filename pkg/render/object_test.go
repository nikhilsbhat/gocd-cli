package render_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	goCdLogger "github.com/nikhilsbhat/gocd-sdk-go/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
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

func TestYAML(t *testing.T) {
	t.Run("", func(t *testing.T) {
		fileData, err := os.ReadFile("/Users/nikhil.bhat/my-opensource/gocd-cli/sampl_pipeline.gocd.yaml")
		assert.NoError(t, err)

		yamlMap := make(map[interface{}]interface{})
		err = yaml.Unmarshal(fileData, &yamlMap)
		assert.NoError(t, err)

		plainYaml, err := yaml.Marshal(yamlMap)
		assert.NoError(t, err)

		var out map[string]interface{}
		err = yaml.Unmarshal(plainYaml, &out)
		assert.NoError(t, err)

		fmt.Printf("%v", out)

		assert.Equal(t, "", out)
	})
}
