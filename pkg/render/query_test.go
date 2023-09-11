package render_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
)

func TestObject_GetQuery(t *testing.T) {
	t.Run("should be able to identify query as where and execute the query successfully", func(t *testing.T) {
		data := []gocd.ConfigRepo{
			{
				ID:       "sample-repo",
				PluginID: "json.config.plugin",
				Material: gocd.Material{
					Type: "git",
					Attributes: gocd.Attribute{
						URL:        "https://github.com/TWChennai/gocd-git-path-sample.git",
						Branch:     "master",
						AutoUpdate: false,
					},
				},
				Rules: []map[string]string{
					{
						"action":    "refer",
						"directive": "allow",
						"resource":  "*",
						"type":      "*",
					},
				},
			},
			{
				ID:       "gocd-go-sdk",
				PluginID: "yaml.config.plugin",
				Material: gocd.Material{
					Type: "git",
					Attributes: gocd.Attribute{
						URL:        "https://github.com/nikhilsbhat/gocd-sdk-go.git",
						Branch:     "master",
						AutoUpdate: false,
					},
				},
				Rules: []map[string]string{
					{
						"action":    "refer",
						"directive": "allow",
						"resource":  "*",
						"type":      "*",
					},
				},
			},
		}

		expected := []interface{}{
			map[string]interface{}{
				"id": "sample-repo",
				"material": map[string]interface{}{
					"attributes": map[string]interface{}{
						"branch": "master",
						"filter": map[string]interface{}{},
						"url":    "https://github.com/TWChennai/gocd-git-path-sample.git",
					},
					"type": "git",
				},
				"plugin_id": "json.config.plugin",
				"rules": []interface{}{map[string]interface{}{
					"action":    "refer",
					"directive": "allow",
					"resource":  "*",
					"type":      "*",
				}},
			},
		}

		query, err := render.SetQuery(data, "[*] | plugin_id = json.config.plugin")
		assert.NoError(t, err)
		response := query.RunQuery()
		assert.Equal(t, expected, response)
	})

	t.Run("should be able to identify query as find executes the query successfully", func(t *testing.T) {
		data := []gocd.ConfigRepo{
			{
				ID:       "sample-repo",
				PluginID: "json.config.plugin",
				Material: gocd.Material{
					Type: "git",
					Attributes: gocd.Attribute{
						URL:        "https://github.com/TWChennai/gocd-git-path-sample.git",
						Branch:     "master",
						AutoUpdate: false,
					},
				},
				Rules: []map[string]string{
					{
						"action":    "refer",
						"directive": "allow",
						"resource":  "*",
						"type":      "*",
					},
				},
			},
			{
				ID:       "gocd-go-sdk",
				PluginID: "yaml.config.plugin",
				Material: gocd.Material{
					Type: "git",
					Attributes: gocd.Attribute{
						URL:        "https://github.com/nikhilsbhat/gocd-sdk-go.git",
						Branch:     "master",
						AutoUpdate: false,
					},
				},
				Rules: []map[string]string{
					{
						"action":    "refer",
						"directive": "allow",
						"resource":  "*",
						"type":      "*",
					},
				},
			},
		}

		expected := []interface{}{"json.config.plugin", "yaml.config.plugin"}

		query, err := render.SetQuery(data, "[*] | plugin_id")
		assert.NoError(t, err)
		assert.Equal(t, "pluck", query.GetQueryType())
		response := query.RunQuery()
		assert.Equal(t, expected, response)
	})
}
