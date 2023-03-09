package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/utils"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func getPluginsCommand() *cobra.Command {
	pluginCommand := &cobra.Command{
		Use:   "plugin",
		Short: "Command to operate on plugins present in GoCD",
		Long: `Command leverages GoCD config repo apis' [https://api.gocd.org/current/#plugin-settings, https://api.gocd.org/current/#plugin-info] to 
GET/CREATE/UPDATE plugins settings or information`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	pluginCommand.SetUsageTemplate(getUsageTemplate())

	pluginCommand.AddCommand(getPluginSettingsCommand())
	pluginCommand.AddCommand(updatePluginSettingsCommand())
	pluginCommand.AddCommand(createPluginSettingsCommand())
	pluginCommand.AddCommand(getPluginsInfoCommand())
	pluginCommand.AddCommand(getPluginInfoCommand())
	pluginCommand.AddCommand(listPluginsCommand())

	return pluginCommand
}

func getPluginSettingsCommand() *cobra.Command {
	getPluginSettingsCmd := &cobra.Command{
		Use:     "get-settings",
		Short:   "Command to GET settings of a specific plugin present in GoCD [https://api.gocd.org/current/#get-plugin-settings]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPluginSettings(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	getPluginSettingsCmd.SetUsageTemplate(getUsageTemplate())

	return getPluginSettingsCmd
}

func createPluginSettingsCommand() *cobra.Command {
	deletePluginSettingsCmd := &cobra.Command{
		Use:     "create-settings",
		Short:   "Command to CREATE settings of a specified plugin present in GoCD [https://api.gocd.org/current/#create-plugin-settings]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var setting gocd.PluginSettings
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case utils.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &setting); err != nil {
					return err
				}
			case utils.FileTypeJSON:
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
		},
	}

	deletePluginSettingsCmd.SetUsageTemplate(getUsageTemplate())

	return deletePluginSettingsCmd
}

func updatePluginSettingsCommand() *cobra.Command {
	updatePluginSettingsCmd := &cobra.Command{
		Use:     "update-settings",
		Short:   "Command to UPDATE settings of a specified plugin present in GoCD [https://api.gocd.org/current/#update-plugin-settings]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var setting gocd.PluginSettings
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case utils.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &setting); err != nil {
					return err
				}
			case utils.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &setting); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
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

	return updatePluginSettingsCmd
}

func getPluginsInfoCommand() *cobra.Command {
	getPluginInformationCmd := &cobra.Command{
		Use:     "get-info-all",
		Short:   "Command to GET information of all plugins present in GoCD [https://api.gocd.org/current/#get-all-plugin-info]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPluginsInfo()
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	getPluginInformationCmd.SetUsageTemplate(getUsageTemplate())

	return getPluginInformationCmd
}

func getPluginInfoCommand() *cobra.Command {
	getPluginInofCmd := &cobra.Command{
		Use:     "get-info",
		Short:   "Command to GET information of a specific plugin present in GoCD [https://api.gocd.org/current/#get-plugin-info]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPluginInfo(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	getPluginInofCmd.SetUsageTemplate(getUsageTemplate())

	return getPluginInofCmd
}

func listPluginsCommand() *cobra.Command {
	listPluginsCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all plugins present in GoCD",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPluginsInfo()
			if err != nil {
				return err
			}

			var pluginList []string

			for _, plugin := range response.Plugins {
				pluginList = append(pluginList, plugin.ID)
			}

			return cliRenderer.Render(strings.Join(pluginList, "\n"))
		},
	}

	listPluginsCmd.SetUsageTemplate(getUsageTemplate())

	return listPluginsCmd
}
