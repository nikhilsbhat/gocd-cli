package cmd

import (
	"github.com/spf13/cobra"
)

func registerJobsCommand() *cobra.Command {
	jobsCommand := &cobra.Command{
		Use:   "job",
		Short: "Command to operate on jobs present in GoCD",
		Long: `Command leverages GoCD job apis'
[https://api.gocd.org/current/#scheduled-jobs] to
GET/SCHEDULE jobs of specific pipelines present GoCD`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	jobsCommand.SetUsageTemplate(getUsageTemplate())

	jobsCommand.AddCommand(getScheduledJobsCommand())

	for _, command := range jobsCommand.Commands() {
		command.SilenceUsage = true
	}

	return jobsCommand
}

func getScheduledJobsCommand() *cobra.Command {
	getScheduledJobsCmd := &cobra.Command{
		Use:     "scheduled",
		Short:   "Command to GET a list of scheduled jobs in GoCD [https://api.gocd.org/current/#scheduled-jobs]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli job scheduled"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetScheduledJobs()
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerPipelineFlags(getScheduledJobsCmd)

	return getScheduledJobsCmd
}
