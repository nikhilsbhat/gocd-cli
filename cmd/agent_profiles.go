package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func registerAgentProfilesCommand() *cobra.Command {
	registerAgentProfilesCmd := &cobra.Command{
		Use:   "elastic-agent-profile",
		Short: "Command to operate on elastic-agent-profile in GoCD [https://api.gocd.org/current/#elastic-agent-profiles]",
		Long: `Command leverages GoCD elastic-agent-profile apis' [https://api.gocd.org/current/#elastic-agent-profiles] to 
GET/CREATE/UPDATE/DELETE elastic agent profiles in GoCD (make sure you have appropriate plugin is installed before using this)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	registerAgentProfilesCmd.SetUsageTemplate(getUsageTemplate())

	registerAgentProfilesCmd.AddCommand(getAgentProfilesCommand())
	registerAgentProfilesCmd.AddCommand(getAgentProfileCommand())
	registerAgentProfilesCmd.AddCommand(createAgentProfileCommand())
	registerAgentProfilesCmd.AddCommand(updateAgentProfileCommand())
	registerAgentProfilesCmd.AddCommand(deleteAgentProfileCommand())
	registerAgentProfilesCmd.AddCommand(listAgentProfilesCommand())

	for _, command := range registerAgentProfilesCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerAgentProfilesCmd
}

func getAgentProfilesCommand() *cobra.Command {
	getElasticAgentProfilesCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all the elastic agent profiles present in GoCD [https://api.gocd.org/current/#get-all-elastic-agent-profiles]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetElasticAgentProfiles()
			if err != nil {
				return err
			}

			if len(queries) != 0 {
				objectString, err := render.Marshal(response)
				if err != nil {
					return err
				}

				cliLogger.Debugf("since --query is passed, applying query '%v' to the output", queries)
				queriedResponse := objectString.GetQuery(queries)

				return cliRenderer.Render(queriedResponse)
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
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetElasticAgentProfile(args[0])
			if err != nil {
				return err
			}

			if len(queries) != 0 {
				objectString, err := render.Marshal(response)
				if err != nil {
					return err
				}

				cliLogger.Debugf("since --query is passed, applying query '%v' to the output", queries)
				queriedResponse := objectString.GetQuery(queries)

				return cliRenderer.Render(queriedResponse)
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
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteElasticAgentProfile(args[0]); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("elastic agent profile '%s' deleted successfully", args[0]))
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
