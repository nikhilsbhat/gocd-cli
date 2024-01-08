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
	"gopkg.in/yaml.v3"
)

func registerClusterProfilesCommand() *cobra.Command {
	registerClusterProfilesCmd := &cobra.Command{
		Use:   "cluster-profile",
		Short: "Command to operate on cluster-profile present in GoCD [https://api.gocd.org/current/#cluster-profiles]",
		Long: `Command leverages GoCD cluster-profile apis' [https://api.gocd.org/current/#cluster-profiles] to 
GET/CREATE/UPDATE/DELETE cluster profiles present in GoCD`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	registerClusterProfilesCmd.SetUsageTemplate(getUsageTemplate())

	registerClusterProfilesCmd.AddCommand(getClusterProfilesCommand())
	registerClusterProfilesCmd.AddCommand(getClusterProfileCommand())
	registerClusterProfilesCmd.AddCommand(createClusterProfileCommand())
	registerClusterProfilesCmd.AddCommand(updateClusterProfileCommand())
	registerClusterProfilesCmd.AddCommand(deleteClusterProfileCommand())
	registerClusterProfilesCmd.AddCommand(listClusterProfilesCommand())

	for _, command := range registerClusterProfilesCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerClusterProfilesCmd
}

func getClusterProfilesCommand() *cobra.Command {
	getClusterProfilesCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all the cluster profiles present in GoCD [https://api.gocd.org/current/#get-all-cluster-profiles]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetClusterProfiles()
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

			return cliRenderer.Render(response.ClusterProfilesConfig)
		},
	}

	return getClusterProfilesCmd
}

func getClusterProfileCommand() *cobra.Command {
	getClusterProfileCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a specific cluster profile present in GoCD [https://api.gocd.org/current/#get-a-cluster-profile]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetClusterProfile(args[0])
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

	return getClusterProfileCmd
}

func createClusterProfileCommand() *cobra.Command {
	createClusterProfileCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE a cluster profile with all specified configurations in GoCD [https://api.gocd.org/current/#create-a-cluster-profile]",
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

			response, err := client.CreateClusterProfile(commonCfg)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("cluster profile %s created successfully", commonCfg.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return createClusterProfileCmd
}

func updateClusterProfileCommand() *cobra.Command {
	updateClusterProfileCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE a cluster profile with all specified configurations in GoCD [https://api.gocd.org/current/#update-a-cluster-profile]",
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

			response, err := client.UpdateClusterProfile(commonCfg)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("cluster profile %s updated successfully", commonCfg.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return updateClusterProfileCmd
}

func deleteClusterProfileCommand() *cobra.Command {
	deleteClusterProfileCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE a specific cluster profile present in GoCD [https://api.gocd.org/current/#delete-a-cluster-profile]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			clusterProfile := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete cluster profile '%s'", clusterProfile)

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

			if err := client.DeleteClusterProfile(clusterProfile); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("cluster profile '%s' deleted successfully", clusterProfile))
		},
	}

	return deleteClusterProfileCmd
}

func listClusterProfilesCommand() *cobra.Command {
	listElasticAgentProfilesCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all cluster profiles present in GoCD [https://api.gocd.org/current/#get-all-cluster-profiles]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetClusterProfiles()
			if err != nil {
				return err
			}

			var clusterProfiles []string

			for _, commonConfig := range response.ClusterProfilesConfig {
				clusterProfiles = append(clusterProfiles, commonConfig.ID)
			}

			return cliRenderer.Render(strings.Join(clusterProfiles, "\n"))
		},
	}

	return listElasticAgentProfilesCmd
}
