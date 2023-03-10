package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/utils"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func registerAgentCommand() *cobra.Command {
	agentsCommand := &cobra.Command{
		Use:   "agents",
		Short: "Command to operate on agents present in GoCD [https://api.gocd.org/current/#agents]",
		Long: `Command leverages GoCD agents apis' [https://api.gocd.org/current/#agents] to 
GET/UPDATE/DELETE GoCD agent also kill task and job run history from an agent`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	agentsCommand.SetUsageTemplate(getUsageTemplate())

	agentsCommand.AddCommand(getAgentsCommand())
	agentsCommand.AddCommand(getAgentCommand())
	agentsCommand.AddCommand(updateAgentCommand())
	agentsCommand.AddCommand(deleteAgentCommand())
	agentsCommand.AddCommand(listAgentsCommand())
	agentsCommand.AddCommand(killTaskCommand())
	agentsCommand.AddCommand(getJobRunHistoryCommand())

	for _, command := range agentsCommand.Commands() {
		command.SilenceUsage = true
	}

	return agentsCommand
}

func getAgentsCommand() *cobra.Command {
	getAgentsCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all the agents present in GoCD [https://api.gocd.org/current/#get-all-agents]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetAgents()
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return getAgentsCmd
}

func getAgentCommand() *cobra.Command {
	getAgentCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET all the agents present in GoCD [https://api.gocd.org/current/#get-one-agent]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetAgent(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return getAgentCmd
}

func updateAgentCommand() *cobra.Command {
	createAgentCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE an agent with all specified configuration [https://api.gocd.org/current/#update-an-agent]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var agent gocd.Agent
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case utils.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &agent); err != nil {
					return err
				}
			case utils.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &agent); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			if err = client.UpdateAgent(agent); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("agent %s updated successfully", agent.ID))
		},
	}

	return createAgentCmd
}

func deleteAgentCommand() *cobra.Command {
	deleteAgentCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE a specific agent present in GoCD [https://api.gocd.org/current/#delete-an-agent]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetAgent(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	return deleteAgentCmd
}

func listAgentsCommand() *cobra.Command {
	listAgentsCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all the agents present in GoCD [https://api.gocd.org/current/#get-all-agents]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			responses, err := client.GetAgents()
			if err != nil {
				return err
			}

			agents := make([]map[string]string, 0)

			for _, response := range responses {
				agent := map[string]string{
					"id":   response.ID,
					"name": response.Name,
				}

				agents = append(agents, agent)
			}

			return cliRenderer.Render(agents)
		},
	}

	return listAgentsCmd
}

func getJobRunHistoryCommand() *cobra.Command {
	jobHistoryCmd := &cobra.Command{
		Use:     "job-history",
		Short:   "Command to GET information of the jobs that ran on a specific agent present in GoCD [https://api.gocd.org/current/#agent-job-run-history]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetAgentJobRunHistory(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response.Jobs)
		},
	}

	return jobHistoryCmd
}

func killTaskCommand() *cobra.Command {
	jobHistoryCmd := &cobra.Command{
		Use:     "kill-task",
		Short:   "Command to KILL a specific task running on a specific agent present in GoCD [https://api.gocd.org/current/#kill-running-tasks]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.AgentKillTask(gocd.Agent{ID: args[0]}); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("task running oon agent with ID %s killed successfully", args[0]))
		},
	}

	return jobHistoryCmd
}
