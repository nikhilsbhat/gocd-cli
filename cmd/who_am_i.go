package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func registerWhoAmICommand() *cobra.Command {
	var raw bool

	whoCmd := &cobra.Command{
		Use:     "who-am-i",
		Short:   "Command to check which user being used by GoCD Command line interface",
		Long:    `Command leverages GoCD current user api [https://api.gocd.org/current/#current-user] to GET current user from GoCD`,
		PreRunE: setCLIClient,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			response, err := client.GetCurrentUser()
			if err != nil {
				return err
			}

			if !raw {
				fmt.Printf("user: %s\n", response.Name)

				return nil
			}

			return cliRenderer.Render(response)
		},
	}

	whoCmd.SetUsageTemplate(getUsageTemplate())

	for _, command := range whoCmd.Commands() {
		command.SilenceUsage = true
	}

	whoCmd.PersistentFlags().BoolVarP(&raw, "raw", "", false,
		"enabling this would print the raw response")

	return whoCmd
}
