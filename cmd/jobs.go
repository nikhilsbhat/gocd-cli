package cmd

import (
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
)

var stageConfig gocd.Stage

func registerJobsCommand() *cobra.Command {
	jobsCommand := &cobra.Command{
		Use:   "job",
		Short: "Command to operate on jobs present in GoCD",
		Long: `Command leverages GoCD job apis'
[https://api.gocd.org/current/#scheduled-jobs] to
SCHEDULE/RUN/RUN-FAILED jobs of specific pipeline present GoCD`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	jobsCommand.SetUsageTemplate(getUsageTemplate())

	jobsCommand.AddCommand(getScheduledJobsCommand())
	jobsCommand.AddCommand(getRunFailedJobsCommand())
	jobsCommand.AddCommand(getRunJobsCommand())

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
		RunE: func(_ *cobra.Command, _ []string) error {
			response, err := client.GetScheduledJobs()
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return getScheduledJobsCmd
}

func getRunFailedJobsCommand() *cobra.Command {
	getScheduledJobsCmd := &cobra.Command{
		Use:     "run-failed",
		Short:   "Command to run failed jobs of specific pipelines in GoCD [https://api.gocd.org/current/#run-failed-jobs]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli job run --pipeline myPipeline --pipeline-counter 2 --stage myStage --stage-counter 3`,
		RunE: func(_ *cobra.Command, _ []string) error {
			response, err := client.RunFailedJobs(stageConfig)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerJobsNStageFlags(getScheduledJobsCmd)

	return getScheduledJobsCmd
}

func getRunJobsCommand() *cobra.Command {
	getScheduledJobsCmd := &cobra.Command{
		Use:     "run",
		Short:   "Command to run list of jobs those are part of selected pipeline in GoCD [https://api.gocd.org/current/#run-selected-jobs]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli job run --pipeline myPipeline --pipeline-counter 2 --stage myStage --stage-counter 3 --job job1 --job job2`,
		RunE: func(_ *cobra.Command, _ []string) error {
			response, err := client.RunJobs(stageConfig)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerJobsNStageFlags(getScheduledJobsCmd)

	return getScheduledJobsCmd
}
