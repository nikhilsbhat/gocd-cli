package cmd

import "github.com/spf13/cobra"

func registerStageCommand() *cobra.Command {
	jobsCommand := &cobra.Command{
		Use:   "stage",
		Short: "Command to operate on stages of a pipeline present in GoCD",
		Long: `Command leverages GoCD job apis'
[https://api.gocd.org/current/#stage-instancess] to
CANCEL/RUN stage of specific pipeline present GoCD`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	jobsCommand.SetUsageTemplate(getUsageTemplate())

	jobsCommand.AddCommand(getCancelStageCommand())
	jobsCommand.AddCommand(getRunStageCommand())

	for _, command := range jobsCommand.Commands() {
		command.SilenceUsage = true
	}

	return jobsCommand
}

func getCancelStageCommand() *cobra.Command {
	getCancelStageCmd := &cobra.Command{
		Use:     "cancel",
		Short:   "Command to cancel specific stage of a pipeline present in GoCD [https://api.gocd.org/current/#run-failed-jobs]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli stage cancel --pipeline myPipeline --pipeline-counter 2 --stage myStage --stage-counter 3`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.CancelStage(stageConfig)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerJobsNStageFlags(getCancelStageCmd)

	return getCancelStageCmd
}

func getRunStageCommand() *cobra.Command {
	getScheduledJobsCmd := &cobra.Command{
		Use:     "run",
		Short:   "Command to run a stage from a selected pipeline present in GoCD [https://api.gocd.org/current/#run-selected-jobs]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli stage run --pipeline myPipeline --pipeline-counter 2 --stage myStage --stage-counter 3 --job job1 --job job2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.RunStage(stageConfig)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerJobsNStageFlags(getScheduledJobsCmd)

	return getScheduledJobsCmd
}
