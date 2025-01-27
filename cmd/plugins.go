package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/query"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type goCdPlugin struct {
	pluginID string
	groovy   bool
	json     bool
	yaml     bool
}

var goCdPluginObj goCdPlugin

func registerPluginsCommand() *cobra.Command {
	pluginCommand := &cobra.Command{
		Use:   "plugin",
		Short: "Command to operate on plugins present in GoCD",
		Long: `Command leverages GoCD config repo apis' [https://api.gocd.org/current/#plugin-settings, https://api.gocd.org/current/#plugin-info] to 
GET/CREATE/UPDATE plugins settings or information`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	pluginCommand.SetUsageTemplate(getUsageTemplate())

	pluginCommand.AddCommand(getPluginSettingsCommand())
	pluginCommand.AddCommand(updatePluginSettingsCommand())
	pluginCommand.AddCommand(createPluginSettingsCommand())
	pluginCommand.AddCommand(getPluginsInfoCommand())
	pluginCommand.AddCommand(getPluginInfoCommand())
	pluginCommand.AddCommand(listPluginsCommand())

	for _, command := range pluginCommand.Commands() {
		command.SilenceUsage = true
	}

	return pluginCommand
}

func getPluginSettingsCommand() *cobra.Command {
	getPluginSettingsCmd := &cobra.Command{
		Use:     "get-settings",
		Short:   "Command to GET settings of a specific plugin present in GoCD [https://api.gocd.org/current/#get-plugin-settings]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			for {
				response, err := client.GetPluginSettings(args[0])
				if err != nil {
					return err
				}

				if len(jsonQuery) != 0 {
					cliLogger.Debugf(queryEnabledMessage, jsonQuery)

					baseQuery, err := query.SetQuery(response, jsonQuery)
					if err != nil {
						return err
					}

					cliLogger.Debug(baseQuery.Print())

					return cliRenderer.Render(baseQuery.RunQuery())
				}

				if err = cliRenderer.Render(response); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	getPluginSettingsCmd.SetUsageTemplate(getUsageTemplate())

	return getPluginSettingsCmd
}

func createPluginSettingsCommand() *cobra.Command {
	createPluginSettingsCmd := &cobra.Command{
		Use:     "create-settings",
		Short:   "Command to CREATE settings of a specified plugin present in GoCD [https://api.gocd.org/current/#create-plugin-settings]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE:    createPluginSettings,
	}

	createPluginSettingsCmd.SetUsageTemplate(getUsageTemplate())

	return createPluginSettingsCmd
}

func updatePluginSettingsCommand() *cobra.Command {
	updatePluginSettingsCmd := &cobra.Command{
		Use:     "update-settings",
		Short:   "Command to UPDATE settings of a specified plugin present in GoCD [https://api.gocd.org/current/#update-plugin-settings]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var setting gocd.PluginSettings
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &setting); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &setting); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			pluginSettingsFetched, err := client.GetPluginSettings(setting.ID)
			if err != nil && !strings.Contains(err.Error(), "404") {
				return err
			}

			if create {
				if reflect.DeepEqual(pluginSettingsFetched, gocd.PluginSettings{}) {
					return createPluginSettings(cmd, nil)
				}
			}

			if len(setting.ETAG) == 0 {
				setting.ETAG = pluginSettingsFetched.ETAG
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "pipeline-settings", setting.ID)

			existing, err := diffCfg.String(pluginSettingsFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			response, err := client.UpdatePluginSettings(setting)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("setting for plugin %s updated successfully", setting.ID)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	updatePluginSettingsCmd.SetUsageTemplate(getUsageTemplate())
	updatePluginSettingsCmd.PersistentFlags().BoolVarP(&create, "create", "", false,
		"if a plugin setting for respective plugin doesn't already exist, run create")

	return updatePluginSettingsCmd
}

func getPluginsInfoCommand() *cobra.Command {
	getPluginInformationCmd := &cobra.Command{
		Use:     "get-info-all",
		Short:   "Command to GET information of all plugins present in GoCD [https://api.gocd.org/current/#get-all-plugin-info]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetPluginsInfo()
				if err != nil {
					return err
				}

				if len(jsonQuery) != 0 {
					cliLogger.Debugf(queryEnabledMessage, jsonQuery)

					baseQuery, err := query.SetQuery(response, jsonQuery)
					if err != nil {
						return err
					}

					cliLogger.Debug(baseQuery.Print())

					return cliRenderer.Render(baseQuery.RunQuery())
				}

				if err = cliRenderer.Render(response); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	getPluginInformationCmd.SetUsageTemplate(getUsageTemplate())

	return getPluginInformationCmd
}

func getPluginInfoCommand() *cobra.Command {
	getPluginInfoCmd := &cobra.Command{
		Use:     "get-info",
		Short:   "Command to GET information of a specific plugin present in GoCD [https://api.gocd.org/current/#get-plugin-info]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			for {
				response, err := client.GetPluginInfo(args[0])
				if err != nil {
					return err
				}

				if len(jsonQuery) != 0 {
					cliLogger.Debugf(queryEnabledMessage, jsonQuery)

					baseQuery, err := query.SetQuery(response, jsonQuery)
					if err != nil {
						return err
					}

					cliLogger.Debug(baseQuery.Print())

					return cliRenderer.Render(baseQuery.RunQuery())
				}

				if err = cliRenderer.Render(response); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	getPluginInfoCmd.SetUsageTemplate(getUsageTemplate())

	return getPluginInfoCmd
}

func listPluginsCommand() *cobra.Command {
	listPluginsCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all plugins present in GoCD",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetPluginsInfo()
				if err != nil {
					return err
				}

				var pluginList []string

				for _, plugin := range response.Plugins {
					pluginList = append(pluginList, plugin.ID)
				}

				if err = cliRenderer.Render(strings.Join(pluginList, "\n")); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	listPluginsCmd.SetUsageTemplate(getUsageTemplate())

	return listPluginsCmd
}

func createPluginSettings(cmd *cobra.Command, _ []string) error {
	var setting gocd.PluginSettings

	object, err := readObject(cmd)
	if err != nil {
		return err
	}

	switch objType := object.CheckFileType(cliLogger); objType {
	case content.FileTypeYAML:
		if err = yaml.Unmarshal([]byte(object), &setting); err != nil {
			return err
		}
	case content.FileTypeJSON:
		if err = json.Unmarshal([]byte(object), &setting); err != nil {
			return err
		}
	default:
		return &errors.UnknownObjectTypeError{Name: objType}
	}

	response, err := client.CreatePluginSettings(setting)
	if err != nil {
		return err
	}

	if err = cliRenderer.Render(fmt.Sprintf("setting for plugin %s created successfully", setting.ID)); err != nil {
		return err
	}

	return cliRenderer.Render(response)
}
