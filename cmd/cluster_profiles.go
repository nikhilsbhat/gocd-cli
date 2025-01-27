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
	"gopkg.in/yaml.v3"
)

func registerClusterProfilesCommand() *cobra.Command {
	registerClusterProfilesCmd := &cobra.Command{
		Use:   "cluster-profile",
		Short: "Command to operate on cluster-profile present in GoCD [https://api.gocd.org/current/#cluster-profiles]",
		Long: `Command leverages GoCD cluster-profile apis' [https://api.gocd.org/current/#cluster-profiles] to 
GET/CREATE/UPDATE/DELETE cluster profiles present in GoCD`,
		RunE: func(cmd *cobra.Command, _ []string) error {
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
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetClusterProfiles()
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

				if err = cliRenderer.Render(response.ClusterProfilesConfig); err != nil {
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

	return getClusterProfilesCmd
}

func getClusterProfileCommand() *cobra.Command {
	getClusterProfileCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a specific cluster profile present in GoCD [https://api.gocd.org/current/#get-a-cluster-profile]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			for {
				response, err := client.GetClusterProfile(args[0])
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

	return getClusterProfileCmd
}

func createClusterProfileCommand() *cobra.Command {
	createClusterProfileCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE a cluster profile with all specified configurations in GoCD [https://api.gocd.org/current/#create-a-cluster-profile]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE:    createClusterProfile,
	}

	return createClusterProfileCmd
}

func updateClusterProfileCommand() *cobra.Command {
	updateClusterProfileCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE a cluster profile with all specified configurations in GoCD [https://api.gocd.org/current/#update-a-cluster-profile]",
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

			clusterProfileFetched, err := client.GetClusterProfile(commonCfg.ID)
			if err != nil && !strings.Contains(err.Error(), "404") {
				return err
			}

			if create {
				if reflect.DeepEqual(clusterProfileFetched, gocd.CommonConfig{}) {
					return createClusterProfile(cmd, nil)
				}
			}

			if len(commonCfg.ETAG) == 0 {
				commonCfg.ETAG = clusterProfileFetched.ETAG
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "cluster-profile", clusterProfileFetched.ID)

			existing, err := diffCfg.String(clusterProfileFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			response, err := client.UpdateClusterProfile(commonCfg)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("cluster profile %s updated successfully", commonCfg.ID)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	updateClusterProfileCmd.PersistentFlags().BoolVarP(&create, "create", "", false,
		"if a cluster profile by this name doesn't already exist, run create")

	return updateClusterProfileCmd
}

func deleteClusterProfileCommand() *cobra.Command {
	deleteClusterProfileCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE a specific cluster profile present in GoCD [https://api.gocd.org/current/#delete-a-cluster-profile]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			clusterProfile := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete cluster profile '%s' [y/n]", clusterProfile)

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
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetClusterProfiles()
				if err != nil {
					return err
				}

				var clusterProfiles []string

				for _, commonConfig := range response.ClusterProfilesConfig {
					clusterProfiles = append(clusterProfiles, commonConfig.ID)
				}

				if err = cliRenderer.Render(strings.Join(clusterProfiles, "\n")); err != nil {
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

func createClusterProfile(cmd *cobra.Command, _ []string) error {
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

	response, err := client.CreateClusterProfile(commonCfg)
	if err != nil {
		return err
	}

	if err = cliRenderer.Render(fmt.Sprintf("cluster profile %s created successfully", commonCfg.Name)); err != nil {
		return err
	}

	return cliRenderer.Render(response)
}
