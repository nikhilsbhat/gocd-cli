package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/query"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

var (
	agentsDisabled                             bool
	agentName, agentID                         string
	agentIDs                                   []string
	agentEnvironments, agentResources, agentOS []string
)

func registerAgentCommand() *cobra.Command {
	agentsCommand := &cobra.Command{
		Use:   "agents",
		Short: "Command to operate on agents present in GoCD [https://api.gocd.org/current/#agents]",
		Long: `Command leverages GoCD agents apis' [https://api.gocd.org/current/#agents] to 
GET/UPDATE/DELETE GoCD agent also kill task and job run history from an agent`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	agentsCommand.SetUsageTemplate(getUsageTemplate())

	agentsCommand.AddCommand(getAgentsCommand())
	agentsCommand.AddCommand(getAgentCommand())
	agentsCommand.AddCommand(updateAgentCommand())
	agentsCommand.AddCommand(disableAgentCommand())
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
		RunE: func(_ *cobra.Command, _ []string) error {
			response, err := client.GetAgents()
			if err != nil {
				return err
			}

			response = filterAgentsResponse(response)

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := query.SetQuery(response, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debug(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(response)
		},
	}

	registerAgentsFilterFlags(getAgentsCmd)

	return getAgentsCmd
}

func getAgentCommand() *cobra.Command {
	getAgentCmd := &cobra.Command{
		Use:   "get",
		Short: "Command to GET all the agents present in GoCD [https://api.gocd.org/current/#get-one-agent]",
		Example: `gocd-cli agents get --name my-gocd-agent
gocd-cli agents get --id 938d1935-bdca-4728-83d5-e96cbf0a4f8b`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				var response gocd.Agent

				if len(agentID) != 0 {
					agentResponse, err := client.GetAgent(agentID)
					if err != nil {
						return err
					}
					response = agentResponse
				}

				if len(agentName) != 0 {
					agentsResponse, err := client.GetAgents()
					if err != nil {
						return err
					}
					for _, agentResponse := range agentsResponse {
						if agentResponse.Name == agentName {
							response = agentResponse

							break
						}
					}

					if reflect.DeepEqual(response, gocd.Agent{}) {
						cliLogger.Infof("agent with name '%s' does not exists in GoCD", agentName)

						return nil
					}
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

				if err := cliRenderer.Render(response); err != nil {
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

	registerAgentsFlags(getAgentCmd)

	return getAgentCmd
}

func updateAgentCommand() *cobra.Command {
	updateAgentCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE an agent with all specified configuration [https://api.gocd.org/current/#update-an-agent]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var agent gocd.Agent
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &agent); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &agent); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			agentFetched, err := client.GetAgent(agent.ID)
			if err != nil {
				return err
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "agent", agentFetched.ID)

			existing, err := diffCfg.String(agentFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			if err = client.UpdateAgent(agent); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("agent %s updated successfully", agent.ID))
		},
	}

	return updateAgentCmd
}

func disableAgentCommand() *cobra.Command {
	var wait bool

	disableAgentCmd := &cobra.Command{
		Use:     "disable",
		Short:   "Command to DISABLE an agent registered with GoCD server [https://api.gocd.org/current/#update-an-agent]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			if err := client.UpdateAgentBulk(gocd.Agent{UUIDS: agentIDs, ConfigState: "Disabled"}); err != nil {
				return err
			}

			startTime := time.Now()

			if wait {
				cliLogger.Debugf("waiting for agent %s to be completely disabled...", agentID)

				for _, id := range agentIDs { //nolint:varnamelen
					for {
						agentInfo, err := client.GetAgent(id)
						if err != nil {
							return err
						}

						if isAgentDisabled(agentInfo) {
							break
						}

						cliLogger.Infof("looks like there is a pipeline scheduled in the agent, waiting for it to complete")

						if time.Since(startTime) > timeout {
							cliLogger.Errorf("timedout waiting for the agent '%s' to get disabled, skipping disablement", id)
							cliLogger.Errorf("wait time crossed the default timeout of '%s', "+
								"try increasing the timeout value using '--timeout'", timeout.String())

							continue
						}

						time.Sleep(delay)
					}
				}
			}

			return cliRenderer.Render(fmt.Sprintf("agent %s disabled successfully", agentID))
		},
	}

	disableAgentCmd.PersistentFlags().StringSliceVarP(&agentIDs, "ids", "", nil,
		"ids of the agent which has to be disabled")
	disableAgentCmd.PersistentFlags().BoolVarP(&wait, "wait", "", false,
		"enable this if you want to wait until the agent is disabled completely")
	disableAgentCmd.PersistentFlags().DurationVarP(&delay, "delay", "", defaultDelay,
		"time delay between each retries that would be made to get the agent status")
	disableAgentCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "", defaultTimeout,
		"timeout, if the operation is not successful in this specified duration")

	disableAgentCmd.MarkFlagsRequiredTogether("wait", "timeout")

	return disableAgentCmd
}

func deleteAgentCommand() *cobra.Command {
	deleteAgentCmd := &cobra.Command{
		Use:   "delete",
		Short: "Command to DELETE a specific agent present in GoCD [https://api.gocd.org/current/#delete-an-agent]",
		Example: `gocd-cli agents delete --name my-gocd-agent
gocd-cli agents delete --id 938d1935-bdca-4728-83d5-e96cbf0a4f8b`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(agentName) != 0 {
				agentsResponse, err := client.GetAgents()
				if err != nil {
					return err
				}
				for _, agentResponse := range agentsResponse {
					if agentResponse.Name == agentName {
						agentID = agentResponse.ID

						break
					}
				}

				if len(agentID) == 0 {
					cliLogger.Errorf("failed to delete agent '%s', as it does not exists in GoCD", agentName)

					return nil
				}
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete agent '%s' [y/n]", agentID)

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

			response, err := client.DeleteAgent(agentID)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	registerAgentsFlags(deleteAgentCmd)

	return deleteAgentCmd
}

func listAgentsCommand() *cobra.Command {
	listAgentsCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all the agents present in GoCD [https://api.gocd.org/current/#get-all-agents]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetAgents()
				if err != nil {
					return err
				}

				response = filterAgentsResponse(response)

				agents := make([]map[string]string, 0)

				for _, agentResponse := range response {
					agent := map[string]string{
						"id":   agentResponse.ID,
						"name": agentResponse.Name,
					}

					agents = append(agents, agent)
				}

				if err = cliRenderer.Render(agents); err != nil {
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

	registerAgentsFilterFlags(listAgentsCmd)

	return listAgentsCmd
}

func getJobRunHistoryCommand() *cobra.Command {
	jobHistoryCmd := &cobra.Command{
		Use:   "job-history",
		Short: "Command to GET information of the jobs that ran on a specific agent present in GoCD [https://api.gocd.org/current/#agent-job-run-history]",
		Example: `gocd-cli agents job-history --name my-gocd-agent
gocd-cli agents job-history --id 938d1935-bdca-4728-83d5-e96cbf0a4f8b`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(agentName) != 0 {
				agentsResponse, err := client.GetAgents()
				if err != nil {
					return err
				}
				for _, agentResponse := range agentsResponse {
					if agentResponse.Name == agentName {
						agentID = agentResponse.ID

						break
					}
				}

				if len(agentID) == 0 {
					cliLogger.Errorf("failed to get job run history from agent '%s', as it does not exists in GoCD", agentName)

					return nil
				}
			}

			response, err := client.GetAgentJobRunHistory(agentID)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response.Jobs)
		},
	}

	registerAgentsFlags(jobHistoryCmd)

	return jobHistoryCmd
}

func killTaskCommand() *cobra.Command {
	killTaskCmd := &cobra.Command{
		Use:   "kill-task",
		Short: "Command to KILL a specific task running on a specific agent present in GoCD [https://api.gocd.org/current/#kill-running-tasks]",
		Args:  cobra.NoArgs,
		Example: `gocd-cli agents kill-task --name my-gocd-agent
gocd-cli agents kill-task --id 938d1935-bdca-4728-83d5-e96cbf0a4f8b`,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			if len(agentName) != 0 {
				agentsResponse, err := client.GetAgents()
				if err != nil {
					return err
				}
				for _, agentResponse := range agentsResponse {
					if agentResponse.Name == agentName {
						agentID = agentResponse.ID

						break
					}
				}

				if len(agentID) == 0 {
					cliLogger.Errorf("failed to kill task from agent '%s', as it does not exists in GoCD", agentName)

					return nil
				}
			}

			if err := client.AgentKillTask(gocd.Agent{ID: agentID}); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("task running on agent with ID %s killed successfully", agentName))
		},
	}

	registerAgentsFlags(killTaskCmd)

	return killTaskCmd
}

func getEnvironmentNames(environments any) []string {
	envs := environments.([]interface{})

	envNames := make([]string, 0)

	for _, value := range envs {
		name := value.(map[string]interface{})
		envNames = append(envNames, name["name"].(string))
	}

	return envNames
}

func filterAgentsResponse(response []gocd.Agent) []gocd.Agent {
	if agentsDisabled {
		response = funk.Filter(response, func(agent gocd.Agent) bool {
			return agent.ConfigState == "Disabled"
		}).([]gocd.Agent)
	}

	if len(agentName) != 0 {
		response = funk.Filter(response, func(agent gocd.Agent) bool {
			return funk.Contains(agent.Name, agentName)
		}).([]gocd.Agent)
	}

	if len(agentOS) != 0 {
		response = funk.Filter(response, func(agent gocd.Agent) bool {
			for _, goCDAgentOS := range agentOS {
				return funk.Contains(agent.OS, goCDAgentOS)
			}

			return false
		}).([]gocd.Agent)
	}

	if len(agentResources) != 0 {
		response = funk.Filter(response, func(agent gocd.Agent) bool {
			for _, resource := range agentResources {
				return funk.Contains(agent.Resources, resource)
			}

			return false
		}).([]gocd.Agent)
	}

	if len(agentEnvironments) != 0 {
		response = funk.Filter(response, func(agent gocd.Agent) bool {
			for _, environment := range agentEnvironments {
				return funk.Contains(getEnvironmentNames(agent.Environments), environment)
			}

			return false
		}).([]gocd.Agent)
	}

	if len(agentOS) != 0 {
		response = funk.Filter(response, func(agent gocd.Agent) bool {
			for _, agOS := range agentOS {
				return funk.Contains(agent.OS, agOS)
			}

			return false
		}).([]gocd.Agent)
	}

	return response
}

func isAgentDisabled(agent gocd.Agent) bool {
	return (agent.BuildState == "Idle" || agent.BuildState == "Unknown") &&
		(agent.CurrentState == "Idle" || agent.CurrentState == "LostContact" || agent.CurrentState == "Unknown" || agent.CurrentState == "Missing") &&
		agent.ConfigState == "Disabled"
}
