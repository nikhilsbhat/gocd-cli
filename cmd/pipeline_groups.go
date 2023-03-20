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

func registerPipelineGroupsCommand() *cobra.Command {
	pipelineGroupCommand := &cobra.Command{
		Use:   "pipeline-group",
		Short: "Command to operate on pipeline groups present in GoCD [https://api.gocd.org/current/#pipeline-group-config]",
		Long: `Command leverages GoCD pipeline group config apis' [https://api.gocd.org/current/#pipeline-group-config] to 
GET/CREATE/UPDATE/DELETE and list GoCD pipeline groups`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
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
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPipelineGroups()
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

	return getPipelineGroupsCmd
}

func getPipelineGroupCommand() *cobra.Command {
	getPipelineGroupCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a specific pipeline group present in GoCD [https://api.gocd.org/current/#get-a-pipeline-group]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline-group get movies --query "pipelines.[*] | name" --yaml
// should return only the list of pipeline names based on the query`,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPipelineGroup(args[0])
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
		RunE: func(cmd *cobra.Command, args []string) error {
			var ppGroup gocd.PipelineGroup
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &ppGroup); err != nil {
					return err
				}
			case render.FileTypeJSON:
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
		},
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
		RunE: func(cmd *cobra.Command, args []string) error {
			var ppGroup gocd.PipelineGroup
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &ppGroup); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &ppGroup); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			fmt.Println(ppGroup)

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

	return updatePipelineGroupCmd
}

func deletePipelineGroupCommand() *cobra.Command {
	deletePipelineGroupCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE the specified pipeline group from GoCD [https://api.gocd.org/current/#delete-a-pipeline-group]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		Example: `gocd-cli pipeline-group delete movies`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeletePipelineGroup(args[0]); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("pipeline group deleted: %s", args[0]))
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
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetPipelineGroups()
			if err != nil {
				return err
			}

			var pipelineGroups []string

			for _, environment := range response {
				pipelineGroups = append(pipelineGroups, environment.Name)
			}

			return cliRenderer.Render(strings.Join(pipelineGroups, "\n"))
		},
	}

	return listPipelineCmd
}
