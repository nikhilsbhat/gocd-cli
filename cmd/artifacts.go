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

func registerArtifactCommand() *cobra.Command {
	agentsCommand := &cobra.Command{
		Use:   "artifact",
		Short: "Command to operate on artifacts store/config present in GoCD",
		Long: `Command leverages GoCD agents apis' [https://api.gocd.org/current/#artifacts-config, https://api.gocd.org/current/#artifact-store] to 
GET/CREATE/UPDATE/DELETE GoCD artifact`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	agentsCommand.SetUsageTemplate(getUsageTemplate())

	agentsCommand.AddCommand(getArtifactStoreCommand())
	agentsCommand.AddCommand(getArtifactStoresCommand())
	agentsCommand.AddCommand(getArtifactConfigCommand())
	agentsCommand.AddCommand(createArtifactStoreCommand())
	agentsCommand.AddCommand(updateArtifactStoreCommand())
	agentsCommand.AddCommand(updateArtifactConfigCommand())
	agentsCommand.AddCommand(deleteArtifactStoreCommand())
	agentsCommand.AddCommand(listArtifactsStoreCommand())
	agentsCommand.AddCommand(killTaskCommand())
	agentsCommand.AddCommand(getJobRunHistoryCommand())

	for _, command := range agentsCommand.Commands() {
		command.SilenceUsage = true
	}

	return agentsCommand
}

func getArtifactStoreCommand() *cobra.Command {
	getArtifactStoreCmd := &cobra.Command{
		Use:     "get-store",
		Short:   "Command to GET an specific artifact store in GoCD [https://api.gocd.org/current/#get-an-artifact-store]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetArtifactStore(args[0])
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

	return getArtifactStoreCmd
}

func getArtifactStoresCommand() *cobra.Command {
	getArtifactStoresCmd := &cobra.Command{
		Use:     "get-all-stores",
		Short:   "Command to GET all the artifact stores present in GoCD [https://api.gocd.org/current/#get-all-artifact-stores]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetArtifactStores()
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

	return getArtifactStoresCmd
}

func getArtifactConfigCommand() *cobra.Command {
	getArtifactsConfigCmd := &cobra.Command{
		Use:     "get-config",
		Short:   "Command to GET a configured artifact configuration GoCD [https://api.gocd.org/current/#get-artifacts-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetArtifactConfig()
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

	return getArtifactsConfigCmd
}

func createArtifactStoreCommand() *cobra.Command {
	getArtifactsStoreCmd := &cobra.Command{
		Use:     "create-store",
		Short:   "Command to CREATE an artifact store with all specified configurations in GoCD [https://api.gocd.org/current/#create-an-artifact-store]",
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

			response, err := client.CreateArtifactStore(commonCfg)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("artifact store %s created successfully", commonCfg.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return getArtifactsStoreCmd
}

func updateArtifactStoreCommand() *cobra.Command {
	updateArtifactsStoreCmd := &cobra.Command{
		Use:     "update-store",
		Short:   "Command to UPDATE an artifact store with all specified configurations in GoCD [https://api.gocd.org/current/#update-an-artifact-store]",
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

			response, err := client.UpdateArtifactStore(commonCfg)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("artifact store %s updated successfully", commonCfg.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return updateArtifactsStoreCmd
}

func updateArtifactConfigCommand() *cobra.Command {
	updateArtifactsConfigCmd := &cobra.Command{
		Use:     "update-config",
		Short:   "Command to UPDATE artifact config specified configurations in GoCD [https://api.gocd.org/current/#update-artifacts-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var artifactInfo gocd.ArtifactInfo
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &artifactInfo); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &artifactInfo); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			response, err := client.UpdateArtifactConfig(artifactInfo)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render("artifact config updated successfully"); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return updateArtifactsConfigCmd
}

func deleteArtifactStoreCommand() *cobra.Command {
	deleteArtifactsStoreCmd := &cobra.Command{
		Use:     "delete-store",
		Short:   "Command to DELETE a specific artifact store present in GoCD [https://api.gocd.org/current/#delete-an-artifact-store]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			storeName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete store '%s'", storeName)

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

			if err := client.DeleteArtifactStore(storeName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("artifact store '%s' deleted successfully", storeName))
		},
	}

	return deleteArtifactsStoreCmd
}

func listArtifactsStoreCommand() *cobra.Command {
	listArtifactsStoreCmd := &cobra.Command{
		Use:     "list-store",
		Short:   "Command to LIST all artifact stores present in GoCD [https://api.gocd.org/current/#get-all-artifact-stores]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetArtifactStores()
			if err != nil {
				return err
			}

			var artifactStores []string

			for _, commonConfig := range response.CommonConfigs {
				artifactStores = append(artifactStores, commonConfig.ID)
			}

			return cliRenderer.Render(strings.Join(artifactStores, "\n"))
		},
	}

	return listArtifactsStoreCmd
}
