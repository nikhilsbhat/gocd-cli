package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/utils"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func getConfigRepoCommand() *cobra.Command {
	configRepoCommand := &cobra.Command{
		Use:   "configrepo",
		Short: "Command to operate on configrepo present in GoCD [https://api.gocd.org/current/#config-repo]",
		Long: `Command leverages GoCD config repo apis' [https://api.gocd.org/current/#config-repo] to 
GET/CREATE/UPDATE/DELETE and trigger update on the same`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	configRepoCommand.SetUsageTemplate(getUsageTemplate())

	configRepoCommand.AddCommand(getConfigRepoTriggerUpdateCommand())
	configRepoCommand.AddCommand(getConfigRepoStatusCommand())
	configRepoCommand.AddCommand(getGetConfigRepoCommand())
	configRepoCommand.AddCommand(getCreateConfigRepoCommand())
	configRepoCommand.AddCommand(getUpdateConfigRepoCommand())
	configRepoCommand.AddCommand(getDeleteConfigRepoCommand())

	return configRepoCommand
}

func getGetConfigRepoCommand() *cobra.Command {
	configGetCommand := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET the config information with a specified ID.",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.ConfigRepoError{Message: "information of only one config-repo can be fetched"}
			}

			response, err := client.GetConfigRepo(args[0])
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(response); err != nil {
				return err
			}

			return nil
		},
	}

	configGetCommand.SetUsageTemplate(getUsageTemplate())

	return configGetCommand
}

func getCreateConfigRepoCommand() *cobra.Command {
	configCreateStatusCommand := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE the config-repo with specified configuration",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var configRepo gocd.ConfigRepo
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case utils.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &configRepo); err != nil {
					return err
				}
			case utils.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &configRepo); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			if err = client.CreateConfigRepo(configRepo); err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("config repo %s created successfully", configRepo.ID)); err != nil {
				return err
			}

			return nil
		},
	}

	configCreateStatusCommand.SetUsageTemplate(getUsageTemplate())

	return configCreateStatusCommand
}

func getUpdateConfigRepoCommand() *cobra.Command {
	configCreateStatusCommand := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE the config-repo",
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var configRepo gocd.ConfigRepo
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case utils.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &configRepo); err != nil {
					return err
				}
			case utils.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &configRepo); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			response, err := client.UpdateConfigRepo(configRepo)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(response); err != nil {
				return err
			}

			return nil
		},
	}

	configCreateStatusCommand.SetUsageTemplate(getUsageTemplate())

	return configCreateStatusCommand
}

func getDeleteConfigRepoCommand() *cobra.Command {
	configUpdateStatusCommand := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE the specified config-repo",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.ConfigRepoError{Message: "information of only one config-repo can be fetched"}
			}

			if err := client.DeleteConfigRepo(args[0]); err != nil {
				return err
			}

			if err := cliRenderer.Render(fmt.Sprintf("config repo deleted: %s", args[0])); err != nil {
				return err
			}

			return nil
		},
	}

	return configUpdateStatusCommand
}

func getConfigRepoStatusCommand() *cobra.Command {
	configStatusCommand := &cobra.Command{
		Use:     "status",
		Short:   "Command to GET the status of config-repo update operation.",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.ConfigRepoError{Message: "status of only one config-repo update operation can be fetched"}
			}

			response, err := client.ConfigRepoStatus(args[0])
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(response); err != nil {
				return err
			}

			return nil
		},
	}

	configStatusCommand.SetUsageTemplate(getUsageTemplate())

	return configStatusCommand
}

func getConfigRepoTriggerUpdateCommand() *cobra.Command {
	configTriggerUpdateCommand := &cobra.Command{
		Use:     "trigger-update",
		Short:   "Command to TRIGGER the update for config-repo to get latest revisions.",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.ConfigRepoError{Message: "updates for only one config-repo can be triggered"}
			}

			response, err := client.ConfigRepoTriggerUpdate(args[0])
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(response); err != nil {
				return err
			}

			return nil
		},
	}

	configTriggerUpdateCommand.SetUsageTemplate(getUsageTemplate())

	return configTriggerUpdateCommand
}
