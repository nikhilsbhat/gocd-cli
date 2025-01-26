package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/query"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

func registerPipelineGroupsCommand() *cobra.Command {
	pipelineGroupCommand := &cobra.Command{
		Use:   "pipeline-group",
		Short: "Command to operate on pipeline groups present in GoCD [https://api.gocd.org/current/#pipeline-group-config]",
		Long: `Command leverages GoCD pipeline group config apis' [https://api.gocd.org/current/#pipeline-group-config] to 
GET/CREATE/UPDATE/DELETE and list GoCD pipeline groups`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	pipelineGroupCommand.SetUsageTemplate(getUsageTemplate())

	pipelineGroupCommand.AddCommand(getPipelineGroupsCommand())
	pipelineGroupCommand.AddCommand(getPipelineGroupCommand())
	pipelineGroupCommand.AddCommand(createPipelineGroupCommand())
	pipelineGroupCommand.AddCommand(updatePipelineGroupCommand())
	pipelineGroupCommand.AddCommand(deletePipelineGroupCommand())
	pipelineGroupCommand.AddCommand(listPipelineGroupsCommand())

	for _, command := range pipelineGroupCommand.Commands() {
		command.SilenceUsage = true
	}

	return pipelineGroupCommand
}

func getPipelineGroupsCommand() *cobra.Command {
	getPipelineGroupsCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all pipeline groups present in GoCD [https://api.gocd.org/current/#get-all-pipeline-groups]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline-group get-all --query "[*] | name eq sample-group"
// should return only one pipeline group 'sample-group', this is as good as running 'gocd-cli pipeline-group get movies'`,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetPipelineGroups()
				if err != nil {
					return err
				}

				if dangling {
					response = funk.Filter(response, func(group gocd.PipelineGroup) bool {
						return len(group.Pipelines) == 0
					}).([]gocd.PipelineGroup)
				}

				if len(jsonQuery) != 0 {
					cliLogger.Debugf(queryEnabledMessage, jsonQuery)

					baseQuery, err := query.SetQuery(response, jsonQuery)
					if err != nil {
						return err
					}

					cliLogger.Debug(baseQuery.Print())

					return cliRenderer.Render(baseQuery.RunQuery())
				}

				if err = cliRenderer.Render(response); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	registerDanglingFlags(getPipelineGroupsCmd)

	return getPipelineGroupsCmd
}

func getPipelineGroupCommand() *cobra.Command {
	getPipelineGroupCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a specific pipeline group present in GoCD [https://api.gocd.org/current/#get-a-pipeline-group]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline-group get movies --query "pipelines.[*] | name" -o yaml
// should return only the list of pipeline names based on the query`,
		RunE: func(_ *cobra.Command, args []string) error {
			for {
				response, err := client.GetPipelineGroup(args[0])
				if err != nil {
					return err
				}

				if len(jsonQuery) != 0 {
					cliLogger.Debugf(queryEnabledMessage, jsonQuery)

					baseQuery, err := query.SetQuery(response, jsonQuery)
					if err != nil {
						return err
					}

					cliLogger.Debug(baseQuery.Print())

					return cliRenderer.Render(baseQuery.RunQuery())
				}

				if err = cliRenderer.Render(response); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	return getPipelineGroupCmd
}

func createPipelineGroupCommand() *cobra.Command {
	createPipelineGroupCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE the pipeline group with all specified configuration [https://api.gocd.org/current/#create-a-pipeline-group]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline-group create movies --from-file pipeline-group-movies.yaml --log-level debug
// the inputs can be passed either from file using '--from-file' flag or entire content as argument to command`,
		RunE: createPipelineGroup,
	}

	return createPipelineGroupCmd
}

func updatePipelineGroupCommand() *cobra.Command {
	updatePipelineGroupCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE the pipeline group with the latest specified configuration [https://api.gocd.org/current/#update-a-pipeline-group]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline-group update movies --from-file pipeline-group-movies.yaml --log-level debug
// the inputs can be passed either from file using '--from-file' flag or entire content as argument to command`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var ppGroup gocd.PipelineGroup
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &ppGroup); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &ppGroup); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			pipelineGroupFetched, err := client.GetPipelineGroup(ppGroup.Name)
			if err != nil && !strings.Contains(err.Error(), "404") {
				return err
			}

			if create {
				if reflect.DeepEqual(pipelineGroupFetched, gocd.PipelineGroup{}) {
					return createPipelineGroup(cmd, nil)
				}
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "pipeline-group", ppGroup.Name)

			existing, err := diffCfg.String(pipelineGroupFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			env, err := client.UpdatePipelineGroup(ppGroup)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("pipeline group %s updated successfully", ppGroup.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(env)
		},
	}

	updatePipelineGroupCmd.PersistentFlags().BoolVarP(&create, "create", "", false,
		"if a pipeline group by this name doesn't already exist, run create")

	return updatePipelineGroupCmd
}

func deletePipelineGroupCommand() *cobra.Command {
	deletePipelineGroupCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE the specified pipeline group from GoCD [https://api.gocd.org/current/#delete-a-pipeline-group]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline-group delete movies`,
		RunE: func(_ *cobra.Command, args []string) error {
			pipelineGroupName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete pipeline-group '%s' [y/n]", pipelineGroupName)

			if !cliCfg.Yes {
				contains, option := cliShellReadConfig.Reader()
				if !contains {
					cliLogger.Fatalln(inputValidationFailureMessage)
				}

				if option.Short == "n" {
					cliLogger.Warn(optingOutMessage)

					os.Exit(0)
				}
			}

			if err := client.DeletePipelineGroup(pipelineGroupName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("pipeline group deleted: %s", pipelineGroupName))
		},
	}

	return deletePipelineGroupCmd
}

func listPipelineGroupsCommand() *cobra.Command {
	listPipelineCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all the pipeline group present in GoCD [https://api.gocd.org/current/#get-all-pipeline-groups]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline-group list`,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetPipelineGroups()
				if err != nil {
					return err
				}

				if dangling {
					response = funk.Filter(response, func(group gocd.PipelineGroup) bool {
						return len(group.Pipelines) == 0
					}).([]gocd.PipelineGroup)
				}

				var pipelineGroups []string

				for _, environment := range response {
					pipelineGroups = append(pipelineGroups, environment.Name)
				}

				if err = cliRenderer.Render(strings.Join(pipelineGroups, "\n")); err != nil {
					return err
				}

				if !cliCfg.Watch {
					break
				}

				time.Sleep(cliCfg.WatchInterval)
			}

			return nil
		},
	}

	registerDanglingFlags(listPipelineCmd)

	return listPipelineCmd
}

func createPipelineGroup(cmd *cobra.Command, _ []string) error {
	var ppGroup gocd.PipelineGroup

	object, err := readObject(cmd)
	if err != nil {
		return err
	}

	switch objType := object.CheckFileType(cliLogger); objType {
	case content.FileTypeYAML:
		if err = yaml.Unmarshal([]byte(object), &ppGroup); err != nil {
			return err
		}
	case content.FileTypeJSON:
		if err = json.Unmarshal([]byte(object), &ppGroup); err != nil {
			return err
		}
	default:
		return &errors.UnknownObjectTypeError{Name: objType}
	}

	if err = client.CreatePipelineGroup(ppGroup); err != nil {
		return err
	}

	return cliRenderer.Render(fmt.Sprintf("pipeline group %s created successfully", ppGroup.Name))
}
