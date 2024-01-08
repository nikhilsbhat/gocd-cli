package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	goYAML "github.com/goccy/go-yaml"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/nikhilsbhat/gocd-sdk-go/pkg/plugin"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

var (
	toCSV                    string
	rawOutput                bool
	goCDPipelineInstance     int
	goCDPipelineMessage      string
	goCDPipelineETAG         string
	goCDPipelineTemplateName string
	goCDPipelines            []string
	goCDPipelineGroups       []string
	goCDEnvironments         []string
	goCDPausePipelineAtStart bool
	goCDPipelinePause        bool
	goCDPipelineUnPause      bool
	numberOfDays             time.Duration
	configRepoNames          []string
)

var defaultGoCDPipelinePatterns = []string{"*.gocd.yaml", "*.gocd.json", "*.gocd.groovy"}

type PipelineVSM struct {
	Pipeline            string   `json:"pipeline,omitempty"             yaml:"pipeline,omitempty"`
	DownstreamPipelines []string `json:"downstream_pipelines,omitempty" yaml:"downstream_pipelines,omitempty"`
	UpstreamPipelines   []string `json:"upstream_pipelines,omitempty"   yaml:"upstream_pipelines,omitempty"`
}

func registerPipelinesCommand() *cobra.Command {
	pipelineCommand := &cobra.Command{
		Use:   "pipeline",
		Short: "Command to operate on pipelines present in GoCD",
		Long: `Command leverages GoCD pipeline apis'
[https://api.gocd.org/current/#pipeline-instances, https://api.gocd.org/current/#pipeline-config, https://api.gocd.org/current/#pipelines] to 
GET/PAUSE/UNPAUSE/UNLOCK/SCHEDULE and comment on a GoCD pipeline`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	pipelineCommand.SetUsageTemplate(getUsageTemplate())

	pipelineCommand.AddCommand(getPipelinesCommand())
	pipelineCommand.AddCommand(getPipelineCommand())
	pipelineCommand.AddCommand(createPipelineCommand())
	pipelineCommand.AddCommand(updatePipelineCommand())
	pipelineCommand.AddCommand(deletePipelineCommand())
	pipelineCommand.AddCommand(getPipelineStateCommand())
	pipelineCommand.AddCommand(getPipelineInstanceCommand())
	pipelineCommand.AddCommand(pauseUnpausePipelineCommand())
	pipelineCommand.AddCommand(schedulePipelineCommand())
	pipelineCommand.AddCommand(commentPipelineCommand())
	pipelineCommand.AddCommand(pipelineExtractTemplateCommand())
	pipelineCommand.AddCommand(listPipelinesCommand())
	pipelineCommand.AddCommand(getPipelineScheduleCommand())
	pipelineCommand.AddCommand(getPipelineHistoryCommand())
	pipelineCommand.AddCommand(getPipelineNotSchedulesCommand())
	pipelineCommand.AddCommand(validatePipelinesCommand())
	pipelineCommand.AddCommand(exportPipelineToConfigRepoFormatCommand())
	pipelineCommand.AddCommand(getPipelineVSMCommand())
	pipelineCommand.AddCommand(getPipelineMapping())
	pipelineCommand.AddCommand(findPipelineFilesCommand())
	pipelineCommand.AddCommand(showPipelineCommand())

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

func getPipelineVSMCommand() *cobra.Command {
	var (
		downStreamPipeline         bool
		upStreamPipeline           bool
		goCDPipelineInstanceNumber []string
	)

	getPipelineVSMCmd := &cobra.Command{
		Use:     "vsm",
		Short:   "Command to GET downstream pipelines of a specified pipeline present in GoCD [https://api.gocd.org/current/#get-pipeline-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline vsm --pipeline animation-movies --pipeline animation-and-action-movies --down-stream --instance animation-movies=14 --yaml"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			vsms := make([]PipelineVSM, 0)

			vsmErrors := make(map[string]string, 0)
			for _, goCDPipeline := range goCDPipelines {
				pipelineHistory, err := client.GetLimitedPipelineRunHistory(goCDPipeline, "10", "0")
				if err != nil {
					return err
				}

				cliLogger.Debugf("run history for pipeline '%s' was fetched successfully", goCDPipeline)

				instance := fmt.Sprintf("%d", pipelineHistory[0].Counter)

				for _, pipelineInstance := range goCDPipelineInstanceNumber {
					filter := strings.Split(pipelineInstance, "=")
					if filter[0] == goCDPipeline {
						cliLogger.Debugf("instance for pipeline '%s' is set to '%s' hence using the same to get VSM", goCDPipeline, filter[1])

						pipelineCounter, err := strconv.Atoi(filter[1])
						if err != nil {
							return err
						}

						if _, err = client.GetPipelineInstance(gocd.PipelineObject{
							Name:    goCDPipeline,
							Counter: pipelineCounter,
						}); err != nil {
							return err
						}

						instance = filter[1]
					}
				}

				response, err := client.GetPipelineVSM(goCDPipeline, instance)
				if err != nil {
					vsmErrors[goCDPipeline] = err.Error()

					continue
				}

				cliLogger.Debugf("VSM details for pipeline '%s' instace '%s' was fetched successfully", goCDPipeline, instance)

				var pipelineStreams []string

				if downStreamPipeline {
					cliLogger.Debugf("since --down-stream is set fetching downstream pipelines")

					pipelineStreams = findDownStreamPipelines(goCDPipeline, response)
				}

				if upStreamPipeline {
					cliLogger.Debugf("since --up-stream is set fetching upstream pipelines")

					pipelineStreams = findUpStreamPipelines(goCDPipeline, response)
				}

				pipelineDependencies, err := parsePipelineConfig(goCDPipeline, pipelineStreams)
				if err != nil {
					return err
				}

				if upStreamPipeline {
					vsms = append(vsms, PipelineVSM{
						Pipeline:          goCDPipeline,
						UpstreamPipelines: pipelineDependencies,
					})
				}

				if downStreamPipeline {
					vsms = append(vsms, PipelineVSM{
						Pipeline:            goCDPipeline,
						DownstreamPipelines: pipelineDependencies,
					})
				}
			}

			if len(vsmErrors) != 0 {
				cliLogger.Errorf("fetching VSM of following pipelines errored")
				for pipeline, vsmError := range vsmErrors {
					cliLogger.Errorf("pipeline '%s': '%s'", pipeline, vsmError)
				}
			}

			if len(toCSV) != 0 {
				return renderVSMtoCSV(vsms, upStreamPipeline)
			}

			return cliRenderer.Render(vsms)
		},
	}

	getPipelineVSMCmd.PersistentFlags().StringSliceVarP(&goCDPipelines, "pipeline", "", nil,
		"name of the pipeline for which the VSM has to be retrieved")
	getPipelineVSMCmd.PersistentFlags().BoolVarP(&downStreamPipeline, "down-stream", "", false,
		"when enabled, will fetch all downstream pipelines of a specified pipeline")
	getPipelineVSMCmd.PersistentFlags().BoolVarP(&upStreamPipeline, "up-stream", "", false,
		"when enabled, will fetch all upstream pipelines of a specified pipeline. (NOTE: flag up-stream is still in experimental phase)")
	getPipelineVSMCmd.PersistentFlags().StringSliceVarP(&goCDPipelineInstanceNumber, "instance", "", nil,
		"instance of the selected pipeline for which the VSM has to be retrieved, the latest VSM number would be picked if not passed. ex: --instance pipeline1=20")
	getPipelineVSMCmd.PersistentFlags().StringVarP(&toCSV, "to-csv", "", "",
		"csv file, to which the output needs to be written")

	getPipelineVSMCmd.MarkFlagsMutuallyExclusive("down-stream", "up-stream")

	if err := getPipelineVSMCmd.MarkPersistentFlagRequired("pipeline"); err != nil {
		cliLogger.Fatalf("%v", err)
	}

	return getPipelineVSMCmd
}

func getPipelineCommand() *cobra.Command {
	getPipelineCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET pipeline config of a specified pipeline present in GoCD [https://api.gocd.org/current/#get-pipeline-config]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline get sample-pipeline --query "[*] | name eq sample-group"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPipelineConfig(args[0])
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

	return getPipelineCmd
}

func getPipelineScheduleCommand() *cobra.Command {
	getPipelineScheduleCmd := &cobra.Command{
		Use:     "last-schedule",
		Short:   "Command to GET last scheduled time of the pipeline present in GoCD [/pipelineHistory.json?pipelineName=nameOfThePipeline]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline last-schedule sample-pipeline`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPipelineSchedules(args[0], "0", "1")
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

			const faultyLength = 2

			if len(response.Groups) == faultyLength {
				if response.Groups[1].History[0].ScheduledDate == "N/A" {
					return nil
				}
			} else {
				if response.Groups[0].History[0].ScheduledDate == "N/A" {
					return nil
				}
			}

			return cliRenderer.Render(response)
		},
	}

	return getPipelineScheduleCmd
}

func getPipelineHistoryCommand() *cobra.Command {
	getPipelineHistoryCmd := &cobra.Command{
		Use:   "history",
		Short: "Command to GET pipeline run history present in GoCD [https://api.gocd.org/current/#get-pipeline-history]",
		Long: `Command leverages GoCD api [https://api.gocd.org/current/#get-pipeline-history] to get the history
This would be an expensive operation especially when you have more pipeline instance to fetch
Prefer invoking this command when GoCD is not serving huge traffic`,
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline history sample-pipeline`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPipelineRunHistory(args[0], "10", delay)
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

	registerPipelineHistoryFlags(getPipelineHistoryCmd)

	return getPipelineHistoryCmd
}

func getPipelineNotSchedulesCommand() *cobra.Command {
	getPipelineNotScheduledCmd := &cobra.Command{
		Use:     "not-scheduled",
		Short:   "Command to GET pipelines not scheduled in last X days from GoCD [/pipelineHistory.json?]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline not-scheduled --time 10`,
		RunE: func(cmd *cobra.Command, args []string) error {
			goCDPipelineNames := make([]string, 0)

			if len(configRepoNames) != 0 {
				cliLogger.Debugf("fetching pipelines from config repo is enabled, hence identifying limiting pipelines to config repo '%v'", configRepoNames)

				for _, configRepoName := range configRepoNames {
					definitions, err := client.GetConfigRepoDefinitions(configRepoName)
					if err != nil {
						return err
					}

					for _, group := range definitions.Groups {
						for _, pipelineName := range group.Pipelines {
							goCDPipelineNames = append(goCDPipelineNames, pipelineName.Name)
						}
					}
				}
			} else {
				cliLogger.Debugf("not limiting config repo while identifying pipelines")

				goCdPipelines, err := client.GetPipelines()
				if err != nil {
					return err
				}

				for _, pipeline := range goCdPipelines.Pipeline {
					pipelineName, err := gocd.GetPipelineName(pipeline.Href)
					if err != nil {
						cliLogger.Errorf("fetching pipeline name from pipline url erored with:, %v", err)

						continue
					}
					goCDPipelineNames = append(goCDPipelineNames, pipelineName)
				}
			}

			pipelineSchedules := make([]gocd.PipelineSchedules, 0)

			for _, pipeline := range goCDPipelineNames {
				cliLogger.Infof("fetching schedules of pipeline '%s'", pipeline)
				response, err := client.GetPipelineSchedules(pipeline, "0", "1")
				if err != nil {
					cliLogger.Errorf("getting schedules for pipline '%s' errored with '%v'", pipeline, err)

					continue
				}

				timeNow := time.Now()
				var timeThen time.Time

				const faultyLength = 2
				if len(response.Groups) == faultyLength {
					if response.Groups[1].History[0].ScheduledDate == "N/A" {
						continue
					}

					timeThen = time.UnixMilli(response.Groups[1].History[0].ScheduledTimestamp).UTC()
				} else {
					if response.Groups[0].History[0].ScheduledDate == "N/A" {
						continue
					}

					timeThen = time.UnixMilli(response.Groups[0].History[0].ScheduledTimestamp).UTC()
				}

				timeDiff := timeNow.Sub(timeThen).Round(1).Hours()
				if timeDiff >= numberOfDays.Hours() {
					pipelineSchedules = append(pipelineSchedules, response)
				}

				time.Sleep(delay)
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := render.SetQuery(pipelineSchedules, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debugf(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(pipelineSchedules)
		},
	}

	registerPipelineHistoryFlags(getPipelineNotScheduledCmd)

	return getPipelineNotScheduledCmd
}

func createPipelineCommand() *cobra.Command {
	createPipelineGroupCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE the pipeline with all specified configuration [https://api.gocd.org/current/#create-a-pipeline]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline create sample-pipeline --from-file sample-pipeline.yaml --log-level debug
// the inputs can be passed either from file using '--from-file' flag or entire content as argument to command`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var pipeline map[string]interface{}
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &pipeline); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &pipeline); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			pipelineConfig := gocd.PipelineConfig{
				Config: pipeline,
			}

			if goCDPausePipelineAtStart {
				pipelineConfig.PausePipeline = true
			}

			if len(goCDPipelineMessage) != 0 {
				pipelineConfig.PauseReason = goCDPipelineMessage
			}

			if _, err = client.CreatePipeline(pipelineConfig); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("pipeline %s created successfully", pipeline["name"]))
		},
	}

	registerPipelineFlags(createPipelineGroupCmd)

	return createPipelineGroupCmd
}

func updatePipelineCommand() *cobra.Command {
	updatePipelineGroupCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE the pipeline config with the latest specified configuration [https://api.gocd.org/current/#edit-pipeline-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline update sample-movies --from-file sample-movies.yaml --log-level debug
// the inputs can be passed either from file using '--from-file' flag or entire content as argument to command`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var pipeline map[string]interface{}
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &pipeline); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &pipeline); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			pipelineConfig := gocd.PipelineConfig{
				ETAG:   goCDPipelineETAG,
				Config: pipeline,
			}

			response, err := client.UpdatePipelineConfig(pipelineConfig)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("pipeline %s updated successfully", pipeline["name"])); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerPipelineFlags(updatePipelineGroupCmd)

	return updatePipelineGroupCmd
}

func deletePipelineCommand() *cobra.Command {
	deletePipelineCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE the specified pipeline from GoCD [https://api.gocd.org/current/#delete-a-pipeline]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline delete movies`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pipelineName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete pipeline '%s'", pipelineName)

			if !cliCfg.Yes {
				contains, option := cliShellReadConfig.Reader()
				if !contains {
					cliLogger.Fatalln("user input validation failed, cannot proceed further")
				}

				if option.Short == "n" {
					cliLogger.Warn("not proceeding further since 'no' was opted")

					os.Exit(0)
				}
			}

			if err := client.DeletePipeline(pipelineName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("pipeline deleted: %s", pipelineName))
		},
	}

	return deletePipelineCmd
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

	registerPipelineFlags(getPipelineInstanceCmd)

	return getPipelineInstanceCmd
}

func pauseUnpausePipelineCommand() *cobra.Command {
	pauseUnpausePipelineCmd := &cobra.Command{
		Use: "action",
		Short: `Command to PAUSE/UNPAUSE a specific pipeline present in GoCD,
              [https://api.gocd.org/current/#pause-a-pipeline,https://api.gocd.org/current/#unpause-a-pipeline]`,
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

	registerPipelineFlags(pauseUnpausePipelineCmd)

	return pauseUnpausePipelineCmd
}

func schedulePipelineCommand() *cobra.Command {
	schedulePipelineCmd := &cobra.Command{
		Use:     "schedule",
		Short:   "Command to SCHEDULE a specific pipeline present in GoCD [https://api.gocd.org/current/#scheduling-pipelines]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline schedule sample --from-file schedule-config.yaml`,
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

			if err = client.SchedulePipeline(args[0], schedule); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("pipeline '%s' scheduled successfully", args[0]))
		},
	}

	registerPipelineFlags(schedulePipelineCmd)

	return schedulePipelineCmd
}

func commentPipelineCommand() *cobra.Command {
	commentOnPipelineCmd := &cobra.Command{
		Use:     "comment",
		Short:   "Command to COMMENT on a specific pipeline instance present in GoCD [https://api.gocd.org/current/#comment-on-pipeline-instance]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline comment --message "message to be commented"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pipelineObject := gocd.PipelineObject{
				Name:    args[0],
				Counter: goCDPipelineInstance,
				Message: goCDPipelineMessage,
			}

			if err := client.CommentOnPipeline(pipelineObject); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("commented on pipeline '%s' successfully", args[0]))
		},
	}

	registerPipelineFlags(commentOnPipelineCmd)

	return commentOnPipelineCmd
}

func pipelineExtractTemplateCommand() *cobra.Command {
	extractTemplatePipelineCmd := &cobra.Command{
		Use:     "template",
		Short:   "Command to EXTRACT template from specific pipeline instance present in GoCD [https://api.gocd.org/current/#extract-template-from-pipeline]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline template --name sample-pipeline --template-name sample-template`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.ExtractTemplatePipeline(args[0], goCDPipelineTemplateName)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerPipelineFlags(extractTemplatePipelineCmd)

	return extractTemplatePipelineCmd
}

func listPipelinesCommand() *cobra.Command {
	listPipelinesCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all the pipelines present in GoCD [https://api.gocd.org/current/#get-feed-of-all-stages-in-a-pipeline]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline list`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var goCdPipelines []string

			if len(goCDEnvironments) != 0 && len(goCDPipelineGroups) != 0 {
				cliLogger.Fatalf("pipelines cannot be filtered by 'environment' and 'pipeline-group' simultaneously")
			}

			if len(goCDEnvironments) != 0 {
				for _, goCDEnvironment := range goCDEnvironments {
					environment, err := client.GetEnvironment(goCDEnvironment)
					if err != nil {
						cliLogger.Errorf("fetching environment '%s' errored with '%s'", goCDEnvironment, err)
					}

					for _, pipeline := range environment.Pipelines {
						goCdPipelines = append(goCdPipelines, pipeline.Name)
					}
				}
			}

			if len(goCDPipelineGroups) != 0 {
				for _, goCDPipelineGroup := range goCDPipelineGroups {
					pipelineGroups, err := client.GetPipelineGroup(goCDPipelineGroup)
					if err != nil {
						cliLogger.Errorf("fetching pipeline group '%s' errored with '%s'", goCDPipelineGroup, err)
					}

					for _, pipeline := range pipelineGroups.Pipelines {
						goCdPipelines = append(goCdPipelines, pipeline.Name)
					}
				}
			}

			if len(goCDPipelineGroups) == 0 && len(goCDEnvironments) == 0 {
				response, err := client.GetPipelines()
				if err != nil {
					return err
				}

				for _, pipeline := range response.Pipeline {
					pipelineName, err := gocd.GetPipelineName(pipeline.Href)
					if err != nil {
						cliLogger.Errorf("fetching pipeline name from pipline url erored with:, %v", err)
					} else {
						goCdPipelines = append(goCdPipelines, pipelineName)
					}
				}
			}

			return cliRenderer.Render(strings.Join(goCdPipelines, "\n"))
		},
	}

	listPipelinesCmd.PersistentFlags().StringSliceVarP(&goCDPipelineGroups, "pipeline-group", "", nil,
		"pipeline group names from which the pipelines needs to be fetched")
	listPipelinesCmd.PersistentFlags().StringSliceVarP(&goCDEnvironments, "environment", "", nil,
		"GoCD environments from which the pipelines needs to be fetched")

	return listPipelinesCmd
}

func validatePipelinesCommand() *cobra.Command {
	type pipelineValidate struct {
		pipelines              []string
		pluginVersion          string
		pluginLocalPath        string
		pluginDownloadURL      string
		fetchVersionFromServer bool
	}

	var pipelineValidateObj pipelineValidate

	validatePipelinesCmd := &cobra.Command{
		Use:     "validate-syntax",
		Short:   "Command validate pipeline syntax by running it against appropriate GoCD plugin",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline validate-syntax --pipeline pipeline1 --pipeline pipeline2`,
		RunE: func(cmd *cobra.Command, args []string) error {
			pluginCfg := plugin.NewPluginConfig(
				pipelineValidateObj.pluginVersion,
				pipelineValidateObj.pluginLocalPath,
				pipelineValidateObj.pluginDownloadURL,
				cliCfg.LogLevel,
			)

			success, err := client.ValidatePipelineSyntax(
				pluginCfg,
				pipelineValidateObj.pipelines,
				!pipelineValidateObj.fetchVersionFromServer,
			)
			if err != nil {
				return err
			}

			if !success {
				cliLogger.Error("oops...!! pipeline syntax validation failed")
				os.Exit(1)
			}

			fmt.Println("SUCCESS")

			return nil
		},
	}

	validatePipelinesCmd.PersistentFlags().StringVarP(&pipelineValidateObj.pluginVersion, "plugin-version", "", "",
		"GoCD plugin version against which the pipeline has to be validated (the plugin type would be auto-detected);"+
			" if missed, the pipeline would be validated against the latest version of the auto-detected plugin")
	validatePipelinesCmd.PersistentFlags().StringSliceVarP(&pipelineValidateObj.pipelines, "pipeline", "", nil,
		"list of pipelines for which the syntax has to be validated")
	validatePipelinesCmd.PersistentFlags().StringVarP(&pipelineValidateObj.pluginDownloadURL, "plugin-download-url", "", "",
		"Auto-detection of the plugin sets the download URL too (Github's release URL);"+
			" if the URL needs to be set to something else, then it can be set using this")
	validatePipelinesCmd.PersistentFlags().StringVarP(&pipelineValidateObj.pluginLocalPath, "plugin-path", "", "",
		"if you prefer managing plugins outside the gocd-cli, the path to already downloaded plugins can be set using this")
	validatePipelinesCmd.PersistentFlags().BoolVarP(&pipelineValidateObj.fetchVersionFromServer, "no-fetch-version-from-server", "", false,
		"if enabled, plugin(auto-detected) version would not be fetched from GoCD server")

	return validatePipelinesCmd
}

func exportPipelineToConfigRepoFormatCommand() *cobra.Command {
	var renderToFile bool

	exportPipelineToConfigRepoFormatCmd := &cobra.Command{
		Use: "export-format",
		Short: "Command to export specified pipeline present in GoCD to appropriate config repo format " +
			"[https://api.gocd.org/current/#export-pipeline-config-to-config-repo-format]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline change-config-repo-format pipeline1 --plugin-id yaml.config.plugin`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.ExportPipelineToConfigRepoFormat(args[0], goCdPluginObj.getPluginID())
			if err != nil {
				return err
			}

			if renderToFile {
				cliLogger.Debugf("--render-to-file is enabled, writing exported plugin to file '%s'", response.PipelineFileName)

				file, err := os.Create(response.PipelineFileName)
				if err != nil {
					return err
				}

				//nolint:mirror
				if _, err = file.Write([]byte(response.PipelineContent)); err != nil {
					return err
				}

				cliLogger.Debug("exported plugin was written to file successfully")

				return nil
			}

			if !rawOutput {
				fmt.Printf("%s\n", response.PipelineContent)

				return nil
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

	commonPluginFlags(exportPipelineToConfigRepoFormatCmd)
	exportPipelineToConfigRepoFormatCmd.PersistentFlags().BoolVarP(&rawOutput, "raw", "", false,
		"if enabled, prints response in raw format")
	exportPipelineToConfigRepoFormatCmd.PersistentFlags().BoolVarP(&renderToFile, "render-to-file", "", false,
		"if enabled, the exported pipeline would we written to a file")

	return exportPipelineToConfigRepoFormatCmd
}

func getPipelineMapping() *cobra.Command {
	getPipelineMappingCmd := &cobra.Command{
		Use:     "get-mappings",
		Short:   "Command to Identify the given pipeline is part of which config-repo/environment of GoCD",
		Example: "gocd-cli pipeline get-mappings --pipeline helm-images --yaml",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			pipelineMappings := make([]map[string]string, 0)

			cliLogger.Debugf("fetching all GoCD environment information to identify which environment the selected pipeline is part of ")

			environmentNames, err := client.GetEnvironments()
			if err != nil {
				return err
			}

			cliLogger.Debugf("all GoCD environment information was fetched successfully")

			pipelineErrors := make(map[string]string, 0)

			for _, goCDPipeline := range goCDPipelines {
				var goCDConfigRepoName, goCDEnvironmentName, originGoCD string

				cliLogger.Debugf("fetching pipeline config to identify which config repo this pipeline is part of")
				pipelineConfig, err := client.GetPipelineConfig(goCDPipeline)
				if err != nil {
					pipelineErrors[goCDPipeline] = err.Error()
				}

				cliLogger.Debugf("pipeline config was retrieved successfully")

				originGoCD = "true"
				if pipelineConfig.Origin.Type != "gocd" {
					goCDConfigRepoName = pipelineConfig.Origin.ID
					originGoCD = "false"
				}

				for _, environmentName := range environmentNames {
					for _, pipeline := range environmentName.Pipelines {
						if pipeline.Name == goCDPipeline {
							goCDEnvironmentName = environmentName.Name
						}
					}
				}

				pipelineMappings = append(pipelineMappings, map[string]string{
					"pipeline":    goCDPipeline,
					"config_repo": goCDConfigRepoName,
					"environment": goCDEnvironmentName,
					"origin_gocd": originGoCD,
				})
			}

			if len(pipelineErrors) != 0 {
				cliLogger.Errorf("fetching mappings of following pipelines errored")
				for pipeline, pipelineError := range pipelineErrors {
					cliLogger.Errorf("pipeline '%s': '%s'", pipeline, pipelineError)
				}
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := render.SetQuery(pipelineMappings, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debugf(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(pipelineMappings)
		},
	}

	getPipelineMappingCmd.PersistentFlags().StringSliceVarP(&goCDPipelines, "pipeline", "", nil,
		"name of the pipeline for which the environment and config repo mappings to be fetched")

	return getPipelineMappingCmd
}

func findPipelineFilesCommand() *cobra.Command {
	var (
		goCDPipelinesPath     string
		goCDPipelinesPatterns []string
		absPath               bool
	)

	findPipelineCmd := &cobra.Command{
		Use:     "find",
		Short:   "Command to find all GoCD pipeline files present in a directory (it recursively finds for pipeline files in all sub-directory)",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline find --path /path/to/pipelines --pattern *.gocd.yaml --pattern *.gocd.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliLogger.Debugf("searching GoCD pipelines under '%s'", goCDPipelinesPath)
			cliLogger.Debug("Below is a list of files that may be identified as GoCD's pipeline files. " +
				"A single file may contain multiple pipeline configurations. Use the command 'gocd-cli pipeline show' to see the pipelines in a given file.")

			return filepath.Walk(goCDPipelinesPath, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				for _, goCDPipelinesPattern := range goCDPipelinesPatterns {
					match, err := filepath.Match(goCDPipelinesPattern, info.Name())
					if err != nil {
						cliLogger.Errorf("matching GoCD pipeline file errored with '%s'", err)
					}

					if match {
						if !absPath {
							fmt.Printf("%s\n", info.Name())

							continue
						}

						fmt.Printf("%s\n", path)
					}
				}

				return nil
			})
		},
	}

	findPipelineCmd.PersistentFlags().StringVarP(&goCDPipelinesPath, "path", "f", "",
		"path to search for all GoCD pipeline files")
	findPipelineCmd.PersistentFlags().StringSliceVarP(&goCDPipelinesPatterns, "pattern", "", defaultGoCDPipelinePatterns,
		"list of patterns to match while searching for all GoCD pipeline files")
	findPipelineCmd.PersistentFlags().BoolVarP(&absPath, "absolute-path", "a", false,
		"when enabled prints absolute path of the pipelines")

	if err := findPipelineCmd.MarkPersistentFlagRequired("path"); err != nil {
		cliLogger.Fatalf("%v", err)
	}

	return findPipelineCmd
}

func showPipelineCommand() *cobra.Command {
	var detailed bool

	showPipelinePipelineCmd := &cobra.Command{
		Use:     "show",
		Short:   "Command to analyse pipelines part of a selected pipeline file",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli show --pipeline /path/to/sample.gocd.yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			detailedPipelineNames := make(map[string][]string, 0)
			pipelineNames := make([]string, 0)

			for _, goCDPipeline := range goCDPipelines {
				cliLogger.Debugf("analysing GoCD pipeline file '%s'", goCDPipeline)

				fileData, err := os.ReadFile(goCDPipeline)
				if err != nil {
					cliLogger.Errorf("reading GoCD pipeline file errored with '%s'", err.Error())

					return err
				}

				type Pipelines struct {
					Config map[string]interface{} `json:"pipelines,omitempty" yaml:"pipelines,omitempty"`
				}

				var fileYAML Pipelines

				object := render.Object(fileData)

				switch objType := object.CheckFileType(cliLogger); objType {
				case render.FileTypeYAML:
					if err = goYAML.Unmarshal(fileData, &fileYAML); err != nil {
						cliLogger.Errorf("deserializing yaml GoCD pipeline file errored with '%s'", err.Error())

						return err
					}

					for _, val := range reflect.ValueOf(fileYAML.Config).MapKeys() {
						pipelineNames = append(pipelineNames, val.String())
					}
				case render.FileTypeJSON:
					var fileJSON map[string]interface{}

					if err = json.Unmarshal(fileData, &fileJSON); err != nil {
						cliLogger.Errorf("deserializing json GoCD pipeline file errored with '%s'", err.Error())

						return err
					}

					pipelineNames = append(pipelineNames, fileJSON["name"].(string))
				default:
					cliLogger.Errorf("the command `pipeline show` does not support reading pipeline config file that was passed")

					return &errors.UnknownObjectTypeError{Name: objType}
				}

				detailedPipelineNames[goCDPipeline] = pipelineNames
			}

			if detailed {
				return cliRenderer.Render(detailedPipelineNames)
			}

			return cliRenderer.Render(pipelineNames)
		},
	}

	showPipelinePipelineCmd.PersistentFlags().StringSliceVarP(&goCDPipelines, "pipelines", "f", nil,
		"path to GoCD pipeline config file to identify pipeline names")
	showPipelinePipelineCmd.PersistentFlags().BoolVarP(&detailed, "detailed", "", false,
		"when enabled prints the information in detail")

	if err := showPipelinePipelineCmd.MarkPersistentFlagRequired("pipelines"); err != nil {
		cliLogger.Fatalf("%v", err)
	}

	return showPipelinePipelineCmd
}

func findDownStreamPipelines(pipelineName string, resp gocd.VSM) []string {
	newParents := []string{pipelineName}

	for _, level := range resp.Level {
		for _, node := range level.Nodes {
			for _, newParent := range newParents {
				if funk.Contains(node.Parents, newParent) {
					newParents = append(newParents, node.Name)
				}
			}
		}
	}

	newParents = GetUniqEntries(newParents)

	return newParents
}

func findUpStreamPipelines(pipelineName string, resp gocd.VSM) []string {
	newChilds := []string{pipelineName}

	for _, level := range resp.Level {
		for _, node := range level.Nodes {
			for _, newChild := range newChilds {
				if funk.Contains(node.Dependents, newChild) {
					newChilds = append(newChilds, node.Name)
				}
			}
		}
	}

	newChilds = GetUniqEntries(newChilds)

	return newChilds
}

func GetUniqEntries(slice []string) []string {
	for slc := 0; slc < len(slice); slc++ {
		if Contains(slice[slc+1:], slice[slc]) {
			slice = append(slice[:slc], slice[slc+1:]...)
			slc--
		}
	}

	return slice
}

func Contains(slice []string, image string) bool {
	for _, slc := range slice {
		if slc == image {
			return true
		}
	}

	return false
}

func parsePipelineConfig(pipelineName string, pipelineStreams []string) ([]string, error) {
	var pipelineDependencies []string

	for _, pipelineStream := range pipelineStreams {
		if pipelineStream != pipelineName {
			pipelineConfig, err := client.GetPipelineConfig(pipelineStream)
			if err != nil {
				return nil, err
			}

			cliLogger.Debugf("config of pipeline '%s' was fetched successfully", pipelineStream)
			cliLogger.Debugf("parsing pipeline '%s' to check the VSM mappings", pipelineStream)

			bytes, err := json.Marshal(pipelineConfig.Config)
			if err != nil {
				return nil, err
			}

			pipelineCfg := string(bytes)

			if gjson.Valid(pipelineCfg) {
				for _, material := range gjson.Get(pipelineCfg, "materials").Array() {
					if funk.Contains(gjson.Get(material.String(), "attributes.name").String(), pipelineName) {
						pipelineDependencies = append(pipelineDependencies, pipelineStream)
					}

					if funk.Contains(gjson.Get(material.String(), "attributes.url").String(), pipelineName) {
						pipelineDependencies = append(pipelineDependencies, pipelineStream)
					}

					if funk.Contains(gjson.Get(material.String(), "attributes.pipeline").String(), pipelineName) {
						pipelineDependencies = append(pipelineDependencies, pipelineStream)
					}
				}

				for _, material := range gjson.Get(pipelineCfg, "parameters").Array() {
					if funk.Contains(gjson.Get(material.String(), "name").String(), pipelineName) {
						pipelineDependencies = append(pipelineDependencies, pipelineStream)
					}

					if funk.Contains(gjson.Get(material.String(), "value").String(), pipelineName) {
						pipelineDependencies = append(pipelineDependencies, pipelineStream)
					}
				}

				for _, stage := range gjson.Get(pipelineCfg, "stages").Array() {
					for _, job := range gjson.Get(stage.String(), "jobs").Array() {
						for _, tasks := range gjson.Get(job.String(), "tasks").Array() {
							if gjson.Get(tasks.String(), "type").String() == "fetch" {
								if funk.Contains(gjson.Get(tasks.String(), "attributes.pipeline").String(), pipelineName) {
									cliLogger.Debugf("pipeline '%s' is mapped as dependency for '%s'", pipelineStream, pipelineName)

									pipelineDependencies = append(pipelineDependencies, pipelineStream)
								}
							}
						}
					}
				}

				cliLogger.Debugf("pipeline '%s' parsed successfully", pipelineConfig.Config["name"])
			}
		}
	}

	pipelineDependencies = GetUniqEntries(pipelineDependencies)

	return pipelineDependencies, nil
}

func renderVSMtoCSV(pipelineVSMs []PipelineVSM, upstream bool) error {
	csvWriter, err := cliRenderer.ToCSV(toCSV)
	if err != nil {
		return err
	}

	defer csvWriter.Flush()

	if err = csvWriter.Write([]string{"Pipeline", "Downstream Pipelines"}); err != nil {
		return err
	}

	for _, pipelineVSM := range pipelineVSMs {
		goCdPipelines := pipelineVSM.DownstreamPipelines
		if upstream {
			goCdPipelines = pipelineVSM.UpstreamPipelines
		}

		if err = csvWriter.Write([]string{pipelineVSM.Pipeline, strings.Join(goCdPipelines, " | ")}); err != nil {
			return err
		}
	}

	return err
}
