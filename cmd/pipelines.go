package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	goCDPipelineInstance int
	goCDPipelineName     string
	goCDPipelineMessage  string
	goCDPipelinePause    bool
	goCDPipelineUnPause  bool
)

func registerPipelinesCommand() *cobra.Command {
	pipelineCommand := &cobra.Command{
		Use:   "pipeline",
		Short: "Command to operate on pipelines present in GoCD ",
		Long: `Command leverages GoCD pipeline group config apis' 
[https://api.gocd.org/current/#pipeline-instances, https://api.gocd.org/current/#pipeline-config, https://api.gocd.org/current/#pipelines] to 
GET/PAUSE/UNPAUSE/UNLOCK/SCHEDULE and comment on a GoCD pipeline`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	pipelineCommand.SetUsageTemplate(getUsageTemplate())

	registerPipelineFlags(pipelineCommand)

	pipelineCommand.AddCommand(getPipelinesCommand())
	pipelineCommand.AddCommand(getPipelineStateCommand())
	pipelineCommand.AddCommand(getPipelineInstanceCommand())
	pipelineCommand.AddCommand(pauseUnpausePipelineCommand())
	pipelineCommand.AddCommand(schedulePipelineCommand())
	pipelineCommand.AddCommand(commentPipelineCommand())
	pipelineCommand.AddCommand(listPipelinesCommand())

	for _, command := range pipelineCommand.Commands() {
		command.SilenceUsage = true
	}

	return pipelineCommand
}

func getPipelinesCommand() *cobra.Command {
	getPipelinesCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all pipelines present in GoCD [https://api.gocd.org/current/#get-feed-of-all-stages-in-a-pipeline]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline get-all --query "[*] | name eq sample-group"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPipelines()
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

	return getPipelinesCmd
}

func getPipelineStateCommand() *cobra.Command {
	getPipelineStateCmd := &cobra.Command{
		Use:     "status",
		Short:   "Command to GET status of a specific pipeline present in GoCD [https://api.gocd.org/current/#get-pipeline-status]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline status sample-pipeline`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPipelineState(args[0])
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

	return getPipelineStateCmd
}

func getPipelineInstanceCommand() *cobra.Command {
	getPipelineInstanceCmd := &cobra.Command{
		Use:     "instance",
		Short:   "Command to GET instance of a specific pipeline present in GoCD [https://api.gocd.org/current/#get-pipeline-instance]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline instance sample-pipeline --instance 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pipelineObject := gocd.PipelineObject{
				Name:    args[0],
				Counter: goCDPipelineInstance,
			}

			response, err := client.GetPipelineInstance(pipelineObject)
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

	return getPipelineInstanceCmd
}

func pauseUnpausePipelineCommand() *cobra.Command {
	pauseUnpausePipelineCmd := &cobra.Command{
		Use: "action",
		Short: "Command to PAUSE/UNPAUSE a specific pipeline present in GoCD " +
			"[https://api.gocd.org/current/#pause-a-pipeline,https://api.gocd.org/current/#unpause-a-pipeline]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline action sample-pipeline --pause/--un-pause`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var action string
			if goCDPipelinePause {
				action = "pausing"
				if err := client.PipelinePause(args[0], goCDPipelineMessage); err != nil {
					return err
				}
			}
			if goCDPipelineUnPause {
				action = "unpausing"
				if err := client.PipelineUnPause(args[0]); err != nil {
					return err
				}
			}

			return cliRenderer.Render(fmt.Sprintf("%s pipeline '%s' was successful", action, args[0]))
		},
	}

	return pauseUnpausePipelineCmd
}

func schedulePipelineCommand() *cobra.Command {
	schedulePipelineCmd := &cobra.Command{
		Use:     "schedule",
		Short:   "Command to SCHEDULE a specific pipeline present in GoCD [https://api.gocd.org/current/#scheduling-pipelines]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline schedule --name sample --from-file schedule-config.yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var schedule gocd.Schedule
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &schedule); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &schedule); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			if err = client.SchedulePipeline(goCDPipelineName, schedule); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("pipeline '%s' scheduled successfully", goCDPipelineName))
		},
	}

	return schedulePipelineCmd
}

func commentPipelineCommand() *cobra.Command {
	commentOnPipelineCmd := &cobra.Command{
		Use:     "comment",
		Short:   "Command to COMMENT on a specific pipeline instance present in GoCD [https://api.gocd.org/current/#comment-on-pipeline-instance]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline comment --message "message to be commented"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pipelineObject := gocd.PipelineObject{
				Name:    goCDPipelineName,
				Counter: goCDPipelineInstance,
				Message: goCDPipelineMessage,
			}

			if err := client.CommentOnPipeline(pipelineObject); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("commented on pipeline '%s' successfully", goCDPipelineName))
		},
	}

	return commentOnPipelineCmd
}

func listPipelinesCommand() *cobra.Command {
	listPipelinesCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all the pipelines present in GoCD [https://api.gocd.org/current/#get-feed-of-all-stages-in-a-pipeline]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPipelines()
			if err != nil {
				return err
			}

			var pipelines []string

			for _, pipeline := range response.Pipeline {
				pipelineName, err := gocd.GetPipelineName(pipeline.Href)
				if err != nil {
					cliLogger.Errorf("fetching pipeline name from pipline url erored with:, %v", err)
				} else {
					pipelines = append(pipelines, pipelineName)
				}
			}

			return cliRenderer.Render(strings.Join(pipelines, "\n"))
		},
	}

	return listPipelinesCmd
}
