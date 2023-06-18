package cmd

import (
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

var materialFilter []string

func registerMaterialsCommand() *cobra.Command {
	registerAgentProfilesCmd := &cobra.Command{
		Use:   "materials",
		Short: "Command to operate on materials present in GoCD [https://api.gocd.org/current/#get-all-materials]",
		Long: `Command leverages GoCD materials apis' [https://api.gocd.org/current/#get-all-materials] to 
GET/LIST and get USAGE of material present in GoCD (make sure you have appropriate plugin is installed before using this)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	registerAgentProfilesCmd.SetUsageTemplate(getUsageTemplate())

	registerAgentProfilesCmd.AddCommand(getMaterialsCommand())
	registerAgentProfilesCmd.AddCommand(getMaterialUsageCommand())
	registerAgentProfilesCmd.AddCommand(listMaterialsCommand())

	for _, command := range registerAgentProfilesCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerAgentProfilesCmd
}

func getMaterialsCommand() *cobra.Command {
	getMaterialsCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET all materials present in GoCD [https://api.gocd.org/current/#get-all-materials]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetMaterials()
			if err != nil {
				return err
			}

			if len(materialFilter) != 0 {
				response = funk.Filter(response, func(material gocd.Material) bool {
					return funk.Contains(materialFilter, material.Attributes.Name)
				}).([]gocd.Material)
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

	registerMaterialFlags(getMaterialsCmd)

	return getMaterialsCmd
}

func listMaterialsCommand() *cobra.Command {
	getMaterialsCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all materials present in GoCD [https://api.gocd.org/current/#get-all-materials]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetMaterials()
			if err != nil {
				return err
			}

			materials := make([]string, 0)
			for _, material := range response {
				if material.Type == "plugin" {
					continue
				}

				if len(material.Attributes.Name) != 0 {
					materials = append(materials, material.Attributes.Name)
				}
				if len(material.Attributes.URL) != 0 {
					materials = append(materials, material.Attributes.URL)
				}
			}

			return cliRenderer.Render(materials)
		},
	}

	return getMaterialsCmd
}

func getMaterialUsageCommand() *cobra.Command {
	getAgentProfilesUsageCmd := &cobra.Command{
		Use:     "usage",
		Short:   "Command to GET an information about pipelines using the specified material",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetMaterialUsage(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerAgentProfileFlags(getAgentProfilesUsageCmd)

	return getAgentProfilesUsageCmd
}
