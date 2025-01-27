package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/query"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

var elasticProfiles []string

func registerAgentProfilesCommand() *cobra.Command {
	registerAgentProfilesCmd := &cobra.Command{
		Use:   "elastic-agent-profile",
		Short: "Command to operate on elastic-agent-profile in GoCD [https://api.gocd.org/current/#elastic-agent-profiles]",
		Long: `Command leverages GoCD elastic-agent-profile apis' [https://api.gocd.org/current/#elastic-agent-profiles] to 
GET/CREATE/UPDATE/DELETE elastic agent profiles in GoCD (make sure you have appropriate plugin is installed before using this)`,
		RunE: func(cmd *cobra.Command, _ []string) error {
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
		RunE: func(_ *cobra.Command, _ []string) error {
			response, err := client.GetElasticAgentProfiles()
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
		RunE: func(_ *cobra.Command, args []string) error {
			for {
				response, err := client.GetElasticAgentProfile(args[0])
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

	return getElasticAgentProfileCmd
}

func createAgentProfileCommand() *cobra.Command {
	createElasticAgentProfileCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE a elastic agent profile with all specified configurations in GoCD [https://api.gocd.org/current/#create-an-elastic-agent-profile]",
		Example: "gocd-cli elastic-agent-profile create sample_ec2 --from-file sample-ec2.yaml --log-level debug",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE:    createAgentProfiles,
	}

	return createElasticAgentProfileCmd
}

func updateAgentProfileCommand() *cobra.Command {
	updateElasticAgentProfileCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE a elastic agent profile with all specified configurations in GoCD [https://api.gocd.org/current/#update-an-elastic-agent-profile]",
		Example: "gocd-cli elastic-agent-profile update sample_ec2 --from-file sample-ec2.yaml --log-level debug",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var commonCfg gocd.CommonConfig
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &commonCfg); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &commonCfg); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			elasticAgentProfileFetched, err := client.GetElasticAgentProfile(commonCfg.ID)
			if err != nil && !strings.Contains(err.Error(), "404") {
				return err
			}

			if create {
				if reflect.DeepEqual(elasticAgentProfileFetched, gocd.CommonConfig{}) {
					return createAgentProfiles(cmd, nil)
				}
			}

			if len(commonCfg.ETAG) == 0 {
				commonCfg.ETAG = elasticAgentProfileFetched.ETAG
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "elastic-agent-profile", elasticAgentProfileFetched.ID)

			existing, err := diffCfg.String(elasticAgentProfileFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			response, err := client.UpdateElasticAgentProfile(commonCfg)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("elastic agent profile %s updated successfully", commonCfg.ID)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	updateElasticAgentProfileCmd.SetUsageTemplate(getUsageTemplate())
	updateElasticAgentProfileCmd.PersistentFlags().BoolVarP(&create, "create", "", false,
		"if a elastic agent profile by this name doesn't already exist, run create")

	return updateElasticAgentProfileCmd
}

func deleteAgentProfileCommand() *cobra.Command {
	deleteElasticAgentProfileCmd := &cobra.Command{
		Use:   "delete",
		Short: "Command to DELETE a specific elastic agent profile present in GoCD [https://api.gocd.org/current/#delete-an-elastic-agent-profile]",
		Example: `gocd-cli elastic-agent-profile delete sample_kubernetes
gocd-cli elastic-agent-profile delete sample_kubernetes -y`,
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			profileName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete elastic-agent-profile '%s' [y/n]", profileName)

			if !cliCfg.Yes {
				contains, option := cliShellReadConfig.Reader()
				if !contains {
					cliLogger.Fatalln(inputValidationFailureMessage)
				}

				if option.Short == "n" {
					cliLogger.Warn(optingOutMessage)

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
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetElasticAgentProfiles()
				if err != nil {
					return err
				}

				var elasticAgentProfiles []string

				for _, commonConfig := range response.CommonConfigs {
					elasticAgentProfiles = append(elasticAgentProfiles, commonConfig.ID)
				}

				if err = cliRenderer.Render(strings.Join(elasticAgentProfiles, "\n")); err != nil {
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

	return listElasticAgentProfilesCmd
}

func getAgentProfilesUsageCommand() *cobra.Command {
	getAgentProfilesUsageCmd := &cobra.Command{
		Use:     "usage",
		Short:   "Command to GET an information about pipelines using elastic agent profiles",
		Example: "gocd-cli elastic-agent-profile usage sample_ec2",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				usageResponse := make([]gocd.ElasticProfileUsage, 0)

				for _, elasticProfile := range elasticProfiles {
					response, err := client.GetElasticAgentProfileUsage(elasticProfile)
					if err != nil {
						return err
					}

					usageResponse = append(usageResponse, response...)
				}

				var elasticAgentProfilesUsage []string

				for _, usage := range usageResponse {
					if !funk.Contains(elasticAgentProfilesUsage, usage.PipelineName) {
						elasticAgentProfilesUsage = append(elasticAgentProfilesUsage, usage.PipelineName)
					}
				}

				switch rawOutput {
				case true:
					if err := cliRenderer.Render(usageResponse); err != nil {
						return err
					}
				default:
					if err := cliRenderer.Render(strings.Join(elasticAgentProfilesUsage, "\n")); err != nil {
						return err
					}
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	registerElasticProfilesFlags(getAgentProfilesUsageCmd)
	registerRawFlags(getAgentProfilesUsageCmd)

	return getAgentProfilesUsageCmd
}

func createAgentProfiles(cmd *cobra.Command, _ []string) error {
	var commonCfg gocd.CommonConfig

	object, err := readObject(cmd)
	if err != nil {
		return err
	}

	switch objType := object.CheckFileType(cliLogger); objType {
	case content.FileTypeYAML:
		if err = yaml.Unmarshal([]byte(object), &commonCfg); err != nil {
			return err
		}
	case content.FileTypeJSON:
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
}
