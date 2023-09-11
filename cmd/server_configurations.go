package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func registerServerConfigCommand() *cobra.Command {
	serverCommand := &cobra.Command{
		Use:   "server-config",
		Short: "Command to operate on GoCD server's configurations",
		Long: `Command leverages GoCD apis':
https://api.gocd.org/current/#update-artifacts-config,
https://api.gocd.org/current/#create-or-update-mailserver-config,
https://api.gocd.org/current/#update-job-timeout-config,
https://api.gocd.org/current/#create-or-update-siteurls-config

to operate on GoCD server's configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	serverCommand.SetUsageTemplate(getUsageTemplate())
	serverCommand.AddCommand(registerSiteCommand())
	serverCommand.AddCommand(registerArtifactManagementCommand())
	serverCommand.AddCommand(registerJobTimeoutCommand())
	serverCommand.AddCommand(registerMailServerConfigCommand())

	for _, command := range serverCommand.Commands() {
		command.SilenceUsage = true
	}

	return serverCommand
}

func registerArtifactManagementCommand() *cobra.Command {
	registerArtifactManagementCmd := &cobra.Command{
		Use:   "artifact-config",
		Short: "Command to operate on GoCD server's artifact configuration",
		Long: `Command leverages GoCD apis':
https://api.gocd.org/current/#get-artifacts-config,
https://api.gocd.org/current/#update-artifacts-config
to operate on artifacts in GoCD`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	registerArtifactManagementCmd.SetUsageTemplate(getUsageTemplate())
	registerArtifactManagementCmd.AddCommand(artifactConfigGetCommand())
	registerArtifactManagementCmd.AddCommand(artifactConfigUpdateCommand())

	for _, command := range registerArtifactManagementCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerArtifactManagementCmd
}

func registerSiteCommand() *cobra.Command {
	serverCommand := &cobra.Command{
		Use:   "site-url",
		Short: "Command to operate on GoCD server's site-url configuration",
		Long: `Command leverages GoCD apis':
https://api.gocd.org/current/#update-artifacts-config,
https://api.gocd.org/current/#create-or-update-mailserver-config

to operate on GoCD's site-url configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	serverCommand.SetUsageTemplate(getUsageTemplate())
	serverCommand.AddCommand(getSiteURLSCommand())
	serverCommand.AddCommand(createUpdateSiteURLSCommand())

	for _, command := range serverCommand.Commands() {
		command.SilenceUsage = true
	}

	return serverCommand
}

func registerJobTimeoutCommand() *cobra.Command {
	registerJobTimeoutCmd := &cobra.Command{
		Use:   "job-timeout",
		Short: "Command to operate on GoCD server's default job timeout",
		Long: `Command leverages GoCD apis':
https://api.gocd.org/current/#get-jobtimeout-config,
https://api.gocd.org/current/#update-job-timeout-config

to operate on GoCD's default job timeout configuration`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	registerJobTimeoutCmd.SetUsageTemplate(getUsageTemplate())
	registerJobTimeoutCmd.AddCommand(getJobTimeoutCommand())
	registerJobTimeoutCmd.AddCommand(updateJobTimeoutCommand())

	for _, command := range registerJobTimeoutCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerJobTimeoutCmd
}

func registerMailServerConfigCommand() *cobra.Command {
	registerJobTimeoutCmd := &cobra.Command{
		Use:   "mail-server-config",
		Short: "Command to operate on GoCD server's mail-server configuration",
		Long: `Command leverages GoCD apis':
https://api.gocd.org/current/#get-mailserver-config,
https://api.gocd.org/current/#create-or-update-mailserver-config,
https://api.gocd.org/current/#delete-mailserver-config

to operate on GoCD server's mail-server configuration'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	registerJobTimeoutCmd.SetUsageTemplate(getUsageTemplate())
	registerJobTimeoutCmd.AddCommand(getMailServerConfigCommand())
	registerJobTimeoutCmd.AddCommand(createUpdateMailServerConfigCommand())
	registerJobTimeoutCmd.AddCommand(deleteMailServerConfigCommand())

	for _, command := range registerJobTimeoutCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerJobTimeoutCmd
}

func getSiteURLSCommand() *cobra.Command {
	getSiteURLSCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to get site urls configured in GoCD server [https://api.gocd.org/current/#get-siteurls-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetSiteURL()
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return getSiteURLSCmd
}

func createUpdateSiteURLSCommand() *cobra.Command {
	var siteURL, secureSiteURL string

	createUpdateSiteURLSCmd := &cobra.Command{
		Use:     "create-or-update",
		Short:   "Command to create/update site urls configured in GoCD server [https://api.gocd.org/current/#create-or-update-siteurls-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			siteURLConfig := gocd.SiteURLConfig{
				SiteURL:       siteURL,
				SecureSiteURL: secureSiteURL,
			}

			response, err := client.CreateOrUpdateSiteURL(siteURLConfig)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	createUpdateSiteURLSCmd.PersistentFlags().StringVarP(&siteURL, "url", "", "",
		"the site URL to be used by GoCD Server to generate links for emails, feeds, etc. Format: [protocol]://[host]:[port]")
	createUpdateSiteURLSCmd.PersistentFlags().StringVarP(&secureSiteURL, "secure-url", "", "",
		"if you wish that your primary site URL be HTTP, but still want to have HTTPS endpoints for the features that require SSL")

	if err := createUpdateSiteURLSCmd.MarkPersistentFlagRequired("url"); err != nil {
		cliLogger.Fatalf("%v", err)
	}

	return createUpdateSiteURLSCmd
}

func artifactConfigGetCommand() *cobra.Command {
	artifactConfigGetCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to get artifact configurations from GoCD server [https://api.gocd.org/current/#get-artifacts-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetArtifactConfig()
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return artifactConfigGetCmd
}

func artifactConfigUpdateCommand() *cobra.Command {
	var artifactDir string
	var purgeStartDiskSpace, purgeUptoDiskSpace float64

	artifactConfigUpdateCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE the artifact configuration in GoCd [https://api.gocd.org/current/#update-artifacts-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetArtifactConfig()
			if err != nil {
				return err
			}

			if len(artifactDir) != 0 {
				response.ArtifactsDir = artifactDir
			}

			response.PurgeSettings.PurgeStartDiskSpace = purgeStartDiskSpace
			response.PurgeSettings.PurgeUptoDiskSpace = purgeUptoDiskSpace

			env, err := client.UpdateArtifactConfig(response)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render("artifact config updated successfully"); err != nil {
				return err
			}

			return cliRenderer.Render(env)
		},
	}

	artifactConfigUpdateCmd.PersistentFlags().StringVarP(&artifactDir, "artifacts-dir", "", "",
		"the directory where GoCD has to store its information, including artefacts published by jobs, "+
			"can be an absolute path or a relative path, which will take the server-installed directory as the root")
	artifactConfigUpdateCmd.PersistentFlags().Float64VarP(&purgeStartDiskSpace, "purge-start-disk-space", "", 0,
		"GoCD starts purging artefacts when disk space is lower than this value")
	artifactConfigUpdateCmd.PersistentFlags().Float64VarP(&purgeUptoDiskSpace, "purge-upto-disk-space", "", 0,
		"GoCD purges artefacts until available disk space is greater than this value")

	return artifactConfigUpdateCmd
}

func getJobTimeoutCommand() *cobra.Command {
	getJobTimeoutCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to get default job timeout GoCD server [https://api.gocd.org/current/#get-jobtimeout-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetDefaultJobTimeout()
			if err != nil {
				return fmt.Errorf("fetching default job timeout errored with: %w", err)
			}

			return cliRenderer.Render(response)
		},
	}

	return getJobTimeoutCmd
}

func updateJobTimeoutCommand() *cobra.Command {
	var jobTimeout int

	updateJobTimeoutCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to update default job timeout in GoCD server [https://api.gocd.org/current/#update-job-timeout-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.UpdateDefaultJobTimeout(jobTimeout); err != nil {
				return fmt.Errorf("updating default job timeout errored with: %w", err)
			}

			return cliRenderer.Render("default job timeout updated successfully")
		},
	}

	updateJobTimeoutCmd.PersistentFlags().IntVarP(&jobTimeout, "timeout", "", 0,
		"default timeout value in minutes for cancelling the hung jobs.")

	if err := updateJobTimeoutCmd.MarkPersistentFlagRequired("timeout"); err != nil {
		cliLogger.Fatalf("%v", err)
	}

	return updateJobTimeoutCmd
}

func getMailServerConfigCommand() *cobra.Command {
	getMailServerConfigCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to get mail server configuration in GoCD server [https://api.gocd.org/current/#get-mailserver-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetMailServerConfig()
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return getMailServerConfigCmd
}

func createUpdateMailServerConfigCommand() *cobra.Command {
	createUpdateMailServerConfigCmd := &cobra.Command{
		Use:     "create-or-update",
		Short:   "Command to create/update mail server config in GoCD server [https://api.gocd.org/current/#create-or-update-mailserver-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var mailConfig gocd.MailServerConfig
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &mailConfig); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &mailConfig); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			output, err := client.CreateOrUpdateMailServerConfig(mailConfig)
			if err != nil {
				return err
			}

			return cliRenderer.Render(output)
		},
	}

	return createUpdateMailServerConfigCmd
}

func deleteMailServerConfigCommand() *cobra.Command {
	deleteMailServerConfigCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to delete mail server configuration in GoCD server [https://api.gocd.org/current/#delete-mailserver-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteMailServerConfig(); err != nil {
				return err
			}

			return cliRenderer.Render("mail server configuration was deleted successfully from GoCD server")
		},
	}

	return deleteMailServerConfigCmd
}
