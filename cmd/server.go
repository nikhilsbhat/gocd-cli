package cmd

import (
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/spf13/cobra"
)

func registerServerCommand() *cobra.Command {
	serverCommand := &cobra.Command{
		Use:   "server",
		Short: "Command to operate on GoCD server health status",
		Long: `Command leverages GoCD health apis' [https://api.gocd.org/current/#server-health-messages, https://api.gocd.org/current/#server-health] to 
GET GoCD server's health and health messages'`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	serverCommand.SetUsageTemplate(getUsageTemplate())

	serverCommand.AddCommand(getHealthCommand())
	serverCommand.AddCommand(getHealthMessagesCommand())

	for _, command := range serverCommand.Commands() {
		command.SilenceUsage = true
	}

	return serverCommand
}

func getHealthCommand() *cobra.Command {
	getHealthCmd := &cobra.Command{
		Use:     "health",
		Short:   "Command to monitor if the GoCD server is up and running [https://api.gocd.org/current/#check-server-health]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetServerHealth()
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return getHealthCmd
}

func getHealthMessagesCommand() *cobra.Command {
	getHealthCmd := &cobra.Command{
		Use: "health-messages",
		Short: "Command to see any errors and warnings generated by the GoCD server as part of " +
			"its normal routine [https://api.gocd.org/current/#get-server-health-messages]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetServerHealthMessages()
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

	return getHealthCmd
}
