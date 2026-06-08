package query_test

import (
	"testing"

	"github.com/nikhilsbhat/gocd-cli/pkg/query"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:testifylint
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

		expected := []any{
			map[string]any{
				"id": "sample-repo",
				"material": map[string]any{
					"attributes": map[string]any{
						"auto_update":           false,
						"branch":                "master",
						"check_externals":       false,
						"filter":                map[string]any{},
						"ignore_for_scheduling": false,
						"invert_filter":         false,
						"url":                   "https://github.com/TWChennai/gocd-git-path-sample.git",
						"use_tickets":           false,
					},
					"config": map[string]any{
						"attributes": map[string]any{
							"auto_update":           false,
							"check_externals":       false,
							"filter":                map[string]any{},
							"ignore_for_scheduling": false,
							"invert_filter":         false,
							"use_tickets":           false,
						},
					},
					"type": "git",
				},
				"parse_info": map[string]any{},
				"plugin_id":  "json.config.plugin",
				"rules": []any{map[string]any{
					"action":    "refer",
					"directive": "allow",
					"resource":  "*",
					"type":      "*",
				}},
			},
		}

		querySet, err := query.SetQuery(data, "[*] | plugin_id = json.config.plugin")

		require.NoError(t, err)

		response := querySet.RunQuery()

		require.Equal(t, expected, response)
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

		expected := []any{"json.config.plugin", "yaml.config.plugin"}

		querySet, err := query.SetQuery(data, "[*] | plugin_id")
		assert.NoError(t, err)
		assert.Equal(t, "pluck", querySet.GetQueryType())
		response := querySet.RunQuery()
		assert.Equal(t, expected, response)
	})
}
