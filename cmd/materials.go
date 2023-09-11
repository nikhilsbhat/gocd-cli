package cmd

import (
	"strconv"
	"strings"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

var (
	materialNames   []string
	materialFilters []string
	materialID      string
	materialFailed  bool
	fetchID         bool
)

func registerMaterialsCommand() *cobra.Command {
	registerAgentProfilesCmd := &cobra.Command{
		Use:   "materials",
		Short: "Command to operate on materials present in GoCD [https://api.gocd.org/current/#get-all-materials]",
		Long: `Command leverages GoCD materials apis' [https://api.gocd.org/current/#get-all-materials] to 
GET/LIST and get USAGE of material present in GoCD (make sure you have appropriate plugin is installed before using this)`,
		Example: "gocd-cli materials [sub-command] [arg] [--flags]",
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
	registerAgentProfilesCmd.AddCommand(triggerMaterialUpdateCommand())

	for _, command := range registerAgentProfilesCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerAgentProfilesCmd
}

func getMaterialsCommand() *cobra.Command {
	getMaterialsCmd := &cobra.Command{
		Use:   "get",
		Short: "Command to GET all materials present in GoCD [https://api.gocd.org/current/#get-all-materials]",
		Example: `gocd-cli materials get --filter type=git
gocd-cli materials get --failed
gocd-cli materials get --names`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetMaterials()
			if err != nil {
				return err
			}

			if materialFailed {
				response = funk.Filter(response, func(material gocd.Material) bool {
					return len(material.Messages) != 0
				}).([]gocd.Material)
			}

			if len(materialNames) != 0 {
				response = funk.Filter(response, func(material gocd.Material) bool {
					return funk.Contains(materialNames, material.Config.Attributes.URL)
				}).([]gocd.Material)
			}

			if len(materialFilters) != 0 {
				for _, materialFilter := range materialFilters {
					filter := strings.Split(materialFilter, "=")

					response = funk.Filter(response, func(material gocd.Material) bool {
						switch strings.ToLower(filter[0]) {
						case "url":
							if funk.Contains(material.Config.Attributes.URL, filter[1]) {
								return true
							}

							return false
						case "type":
							if funk.Contains(material.Config.Type, filter[1]) {
								return true
							}

							return false
						case "can_update":
							boolValue, _ := strconv.ParseBool(strings.ToLower(filter[1]))

							return material.CanTriggerUpdate == boolValue
						case "auto_update":
							boolValue, _ := strconv.ParseBool(strings.ToLower(filter[1]))

							return material.Config.Attributes.AutoUpdate == boolValue
						}

						return false
					}).([]gocd.Material)
				}
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
		Example: "gocd-cli materials list --yaml (only lists materials that has name or URL)",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetMaterials()
			if err != nil {
				return err
			}

			materials := make([]string, 0)
			for _, material := range response {
				if material.Config.Type == "plugin" {
					continue
				}

				if len(material.Config.Attributes.Name) != 0 {
					materials = append(materials, material.Config.Attributes.Name)
				}
				if len(material.Config.Attributes.URL) != 0 {
					materials = append(materials, material.Config.Attributes.URL)
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
		Example: "gocd-cli materials usage https://github.com/nikhilsbhat/helm-drift.git --fetch-id",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			materialID = args[0]

			if fetchID {
				materials, err := client.GetMaterials()
				if err != nil {
					return err
				}

				materials = funk.Filter(materials, func(material gocd.Material) bool {
					return funk.Contains(material.Config.Attributes.URL, args[0])
				}).([]gocd.Material)

				if len(materials) == 0 {
					return &errors.MaterialError{Message: "no material found with the specified URL/Name"}
				}

				materialID = materials[0].Config.Fingerprint
			}

			response, err := client.GetMaterialUsage(materialID)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerRawFlags(getAgentProfilesUsageCmd)
	getAgentProfilesUsageCmd.PersistentFlags().BoolVarP(&fetchID, "fetch-id", "", false,
		"when enabled tries to fetch the ID of the material to get the usages. Do not set this flag if ID is passed")

	return getAgentProfilesUsageCmd
}

func triggerMaterialUpdateCommand() *cobra.Command {
	triggerMaterialUpdateCmd := &cobra.Command{
		Use:     "trigger-update",
		Short:   "Command to trigger update on the specified material",
		Example: "gocd-cli materials trigger-update https://github.com/nikhilsbhat/helm-drift.git --fetch-id",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			materialID = args[0]

			if fetchID {
				materials, err := client.GetMaterials()
				if err != nil {
					return err
				}

				materials = funk.Filter(materials, func(material gocd.Material) bool {
					return funk.Contains(material.Config.Attributes.URL, args[0])
				}).([]gocd.Material)

				if len(materials) == 0 {
					return &errors.MaterialError{Message: "no material found with the specified URL/Name"}
				}

				materialID = materials[0].Config.Fingerprint
			}

			response, err := client.MaterialTriggerUpdate(materialID)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	triggerMaterialUpdateCmd.PersistentFlags().BoolVarP(&fetchID, "fetch-id", "", false,
		"when enabled tries to fetch the ID of the material to get trigger the updates of them. Do not set this flag if ID is passed")

	return triggerMaterialUpdateCmd
}
