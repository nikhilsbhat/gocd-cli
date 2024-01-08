package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

func registerAgentProfilesCommand() *cobra.Command {
	registerAgentProfilesCmd := &cobra.Command{
		Use:   "elastic-agent-profile",
		Short: "Command to operate on elastic-agent-profile in GoCD [https://api.gocd.org/current/#elastic-agent-profiles]",
		Long: `Command leverages GoCD elastic-agent-profile apis' [https://api.gocd.org/current/#elastic-agent-profiles] to 
GET/CREATE/UPDATE/DELETE elastic agent profiles in GoCD (make sure you have appropriate plugin is installed before using this)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	registerAgentProfilesCmd.SetUsageTemplate(getUsageTemplate())

	registerAgentProfilesCmd.AddCommand(getAgentProfilesCommand())
	registerAgentProfilesCmd.AddCommand(getAgentProfileCommand())
	registerAgentProfilesCmd.AddCommand(createAgentProfileCommand())
	registerAgentProfilesCmd.AddCommand(updateAgentProfileCommand())
	registerAgentProfilesCmd.AddCommand(deleteAgentProfileCommand())
	registerAgentProfilesCmd.AddCommand(listAgentProfilesCommand())
	registerAgentProfilesCmd.AddCommand(getAgentProfilesUsageCommand())

	for _, command := range registerAgentProfilesCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerAgentProfilesCmd
}

func getAgentProfilesCommand() *cobra.Command {
	getElasticAgentProfilesCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all the elastic agent profiles present in GoCD [https://api.gocd.org/current/#get-all-elastic-agent-profiles]",
		Example: "gocd-cli elastic-agent-profile get-all",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetElasticAgentProfiles()
			if err != nil {
				return err
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := render.SetQuery(response, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debugf(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(response.CommonConfigs)
		},
	}

	return getElasticAgentProfilesCmd
}

func getAgentProfileCommand() *cobra.Command {
	getElasticAgentProfileCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a specific elastic agent profile present in GoCD [https://api.gocd.org/current/#get-an-elastic-agent-profile]",
		Example: "gocd-cli elastic-agent-profile get sample_kubernetes",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetElasticAgentProfile(args[0])
			if err != nil {
				return err
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := render.SetQuery(response, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debugf(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(response)
		},
	}

	return getElasticAgentProfileCmd
}

func createAgentProfileCommand() *cobra.Command {
	createElasticAgentProfileCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE a elastic agent profile with all specified configurations in GoCD [https://api.gocd.org/current/#create-an-elastic-agent-profile]",
		Example: "gocd-cli elastic-agent-profile create sample_ec2 --from-file sample-ec2.yaml --log-level debug",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var commonCfg gocd.CommonConfig
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &commonCfg); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &commonCfg); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			response, err := client.CreateElasticAgentProfile(commonCfg)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("elastic agent profile %s created successfully", commonCfg.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return createElasticAgentProfileCmd
}

func updateAgentProfileCommand() *cobra.Command {
	updateElasticAgentProfileCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE a elastic agent profile with all specified configurations in GoCD [https://api.gocd.org/current/#update-an-elastic-agent-profile]",
		Example: "gocd-cli elastic-agent-profile update sample_ec2 --from-file sample-ec2.yaml --log-level debug",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var commonCfg gocd.CommonConfig
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &commonCfg); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &commonCfg); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			response, err := client.UpdateElasticAgentProfile(commonCfg)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("elastic agent profile %s updated successfully", commonCfg.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return updateElasticAgentProfileCmd
}

func deleteAgentProfileCommand() *cobra.Command {
	deleteElasticAgentProfileCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE a specific elastic agent profile present in GoCD [https://api.gocd.org/current/#delete-an-elastic-agent-profile]",
		Example: "gocd-cli elastic-agent-profile delete sample_kubernetes",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete elastic-agent-profile '%s'", profileName)

			if !cliCfg.Yes {
				contains, option := cliShellReadConfig.Reader()
				if !contains {
					cliLogger.Fatalln("user input validation failed, cannot proceed further")
				}

				if option.Short == "n" {
					cliLogger.Warn("not proceeding further since 'no' was opted")

					os.Exit(0)
				}
			}

			if err := client.DeleteElasticAgentProfile(profileName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("elastic agent profile '%s' deleted successfully", profileName))
		},
	}

	return deleteElasticAgentProfileCmd
}

func listAgentProfilesCommand() *cobra.Command {
	listElasticAgentProfilesCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all elastic agent profiles present in GoCD [https://api.gocd.org/current/#get-all-elastic-agent-profiles]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetElasticAgentProfiles()
			if err != nil {
				return err
			}

			var elasticAgentProfiles []string

			for _, commonConfig := range response.CommonConfigs {
				elasticAgentProfiles = append(elasticAgentProfiles, commonConfig.ID)
			}

			return cliRenderer.Render(strings.Join(elasticAgentProfiles, "\n"))
		},
	}

	return listElasticAgentProfilesCmd
}

func getAgentProfilesUsageCommand() *cobra.Command {
	getAgentProfilesUsageCmd := &cobra.Command{
		Use:     "usage",
		Short:   "Command to GET an information about pipelines using elastic agent profiles",
		Example: "gocd-cli elastic-agent-profile usage sample_ec2",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetElasticAgentProfileUsage(args[0])
			if err != nil {
				return err
			}

			if rawOutput {
				return cliRenderer.Render(response)
			}

			var elasticAgentProfilesUsage []string

			for _, usage := range response {
				if !funk.Contains(elasticAgentProfilesUsage, usage.PipelineName) {
					elasticAgentProfilesUsage = append(elasticAgentProfilesUsage, usage.PipelineName)
				}
			}

			return cliRenderer.Render(strings.Join(elasticAgentProfilesUsage, "\n"))
		},
	}

	registerRawFlags(getAgentProfilesUsageCmd)

	return getAgentProfilesUsageCmd
}
