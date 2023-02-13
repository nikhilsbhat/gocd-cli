package cmd

import (
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/spf13/cobra"
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
	registerConfigRepoFlags(configRepoCommand)
	configRepoCommand.SetUsageTemplate(getUsageTemplate())

	configRepoCommand.AddCommand(getConfigRepoTriggerUpdateCommand())
	configRepoCommand.AddCommand(getConfigRepoStatusCommand())
	// configRepoCommand.AddCommand(getGetConfigRepoCommand())
	// configRepoCommand.AddCommand(getCreateConfigRepoCommand())
	// configRepoCommand.AddCommand(getUpdateConfigRepoCommand())
	// configRepoCommand.AddCommand(getDeleteConfigRepoCommand())

	return configRepoCommand
}

func getGetConfigRepoCommand() *cobra.Command {
	return &cobra.Command{}
}

func getUpdateConfigRepoCommand() *cobra.Command {
	return &cobra.Command{}
}

func getDeleteConfigRepoCommand() *cobra.Command {
	return &cobra.Command{}
}

func getCreateConfigRepoCommand() *cobra.Command {
	return &cobra.Command{}
}

func getConfigRepoStatusCommand() *cobra.Command {
	configTriggerStatusCommand := &cobra.Command{
		Use:     "status",
		Short:   "Command to get the status of config repository update operation.",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setGoCDClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.AuthError{Message: "args cannot be more than one. At once, status of only one configrepo update operation can be fetched"}
			}

			response, err := client.ConfigRepoStatus(args[0])
			if err != nil {
				return err
			}

			if err = render(response); err != nil {
				return err
			}

			return nil
		},
	}

	configTriggerStatusCommand.SetUsageTemplate(getUsageTemplate())

	return configTriggerStatusCommand
}

func getConfigRepoTriggerUpdateCommand() *cobra.Command {
	configTriggerUpdateCommand := &cobra.Command{
		Use:     "trigger-update",
		Short:   "Command to trigger the update for config repository to get latest revisions.",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: setGoCDClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return &errors.AuthError{Message: "args cannot be more than one. At once, updates for only one configrepo can be triggered"}
			}

			response, err := client.ConfigRepoTriggerUpdate(args[0])
			if err != nil {
				return err
			}

			if err = render(response); err != nil {
				return err
			}

			return nil
		},
	}

	configTriggerUpdateCommand.SetUsageTemplate(getUsageTemplate())

	return configTriggerUpdateCommand
}
