package cmd

import (
	"fmt"

	"github.com/nikhilsbhat/gocd-cli/pkg/query"
	"github.com/spf13/cobra"
)

var (
	goCDEnableMaintenance  bool
	goCDDisableMaintenance bool
)

func registerMaintenanceCommand() *cobra.Command {
	maintenanceCommand := &cobra.Command{
		Use:   "maintenance",
		Short: "Command to operate on maintenance modes in GoCD [https://api.gocd.org/current/#maintenance-mode]",
		Long: `Command leverages GoCD environments apis' [https://api.gocd.org/current/#maintenance-mode] to 
ENABLE/DISABLE/GET maintenance mode information from GoCD`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	maintenanceCommand.SetUsageTemplate(getUsageTemplate())

	registerMaintenanceFlags(maintenanceCommand)

	maintenanceCommand.AddCommand(enableOrDisableMaintenanceCommand())
	maintenanceCommand.AddCommand(getMaintenanceCommand())

	for _, command := range maintenanceCommand.Commands() {
		command.SilenceUsage = true
	}

	return maintenanceCommand
}

func getMaintenanceCommand() *cobra.Command {
	getMaintenanceCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a maintenance mode information from GoCD [https://api.gocd.org/current/#get-maintenance-mode-info]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli maintenance get`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetMaintenanceModeInfo()
			if err != nil {
				return err
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := query.SetQuery(response, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debugf(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(response)
		},
	}

	return getMaintenanceCmd
}

func enableOrDisableMaintenanceCommand() *cobra.Command {
	enableDisableMaintenanceModeCmd := &cobra.Command{
		Use: "action",
		Short: `Command to ENABLE/DISABLE maintenance mode in GoCD, 
              [https://api.gocd.org/current/#enable-maintenance-mode,https://api.gocd.org/current/#disable-maintenance-mode]`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli maintenance --enable/--disable`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var action string
			if goCDEnableMaintenance {
				action = "enabling"
				if err := client.EnableMaintenanceMode(); err != nil {
					return err
				}
			}
			if goCDDisableMaintenance {
				action = "disabling"
				if err := client.DisableMaintenanceMode(); err != nil {
					return err
				}
			}

			return cliRenderer.Render(fmt.Sprintf("%s maintenance mode was successful", action))
		},
	}

	return enableDisableMaintenanceModeCmd
}
