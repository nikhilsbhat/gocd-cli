package cmd

import (
	"github.com/spf13/cobra"
)

func registerIHaveCommand() *cobra.Command {
	var entity string

	iCanCmd := &cobra.Command{
		Use:     "i-have",
		Short:   "Command to check the permissions that the current user has",
		Long:    `Command leverages GoCD permissions api [https://api.gocd.org/current/#permissions] to show all permission that current user has in GoCD`,
		PreRunE: setCLIClient,
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			query := map[string]string{"type": entity}

			response, err := client.GetPermissions(query)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	iCanCmd.SetUsageTemplate(getUsageTemplate())

	for _, command := range iCanCmd.Commands() {
		command.SilenceUsage = true
	}

	iCanCmd.PersistentFlags().StringVarP(&entity, "entity", "e", "",
		"type of GoCD entity to filter the permissions")

	return iCanCmd
}
