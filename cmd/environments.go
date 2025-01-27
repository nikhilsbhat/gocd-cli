package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/query"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

var (
	dangling             bool
	fetchPipeline        bool
	environmentVariables []string
)

func registerEnvironmentsCommand() *cobra.Command {
	environmentCommand := &cobra.Command{
		Use:   "environment",
		Short: "Command to operate on environments present in GoCD [https://api.gocd.org/current/#environment-config]",
		Long: `Command leverages GoCD environments apis' [https://api.gocd.org/current/#environment-config] to 
GET/CREATE/UPDATE/PATCH/DELETE and list GoCD environments`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	environmentCommand.SetUsageTemplate(getUsageTemplate())

	environmentCommand.AddCommand(getEnvironmentsCommand())
	environmentCommand.AddCommand(getEnvironmentCommand())
	environmentCommand.AddCommand(getEnvironmentMapping())
	environmentCommand.AddCommand(createEnvironmentCommand())
	environmentCommand.AddCommand(updateEnvironmentCommand())
	environmentCommand.AddCommand(patchEnvironmentCommand())
	environmentCommand.AddCommand(deleteEnvironmentCommand())
	environmentCommand.AddCommand(listEnvironmentsCommand())

	for _, command := range environmentCommand.Commands() {
		command.SilenceUsage = true
	}

	return environmentCommand
}

func getEnvironmentsCommand() *cobra.Command {
	getEnvironmentsCmd := &cobra.Command{
		Use:   "get-all",
		Short: "Command to GET all the environments present in GoCD [https://api.gocd.org/current/#get-all-environments]",
		Example: `gocd-cli environment get-all -o yaml
gocd-cli environment get-all --env-var ENVIRONMENT_VAR_1 --env-var ENVIRONMENT_VAR_2 -o yaml
gocd-cli environment get-all --pipelines -o yaml`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			response, err := client.GetEnvironments()
			if err != nil {
				return err
			}

			if fetchPipeline {
				pipelineName := make([]string, 0)
				for _, environment := range response {
					for _, name := range environment.Pipelines {
						pipelineName = append(pipelineName, name.Name)
					}
				}

				return cliRenderer.Render(pipelineName)
			}

			if len(environmentVariables) != 0 {
				envVars := make([]gocd.EnvVars, 0)
				for _, environment := range response {
					for _, envVar := range environment.EnvVars {
						for _, environmentVariable := range environmentVariables {
							if funk.Contains(envVar.Name, environmentVariable) {
								envVars = append(envVars, envVar)
							}
						}
					}
				}

				return cliRenderer.Render(envVars)
			}

			if dangling {
				response = funk.Filter(response, func(environment gocd.Environment) bool {
					return len(environment.Pipelines) == 0
				}).([]gocd.Environment)
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

			return cliRenderer.Render(response)
		},
	}

	getEnvironmentsCmd.PersistentFlags().BoolVarP(&fetchPipeline, "pipelines", "", false,
		"when set fetches the pipeline alone")
	getEnvironmentsCmd.PersistentFlags().StringSliceVarP(&environmentVariables, "env-var", "", nil,
		"list of environment variables to fetch from the GoCD environment")

	registerDanglingFlags(getEnvironmentsCmd)

	return getEnvironmentsCmd
}

func getEnvironmentCommand() *cobra.Command {
	getEnvironmentCmd := &cobra.Command{
		Use:   "get",
		Short: "Command to GET a specific environments present in GoCD [https://api.gocd.org/current/#get-environment-config]",
		Example: `gocd-cli environment get gocd_environment_1
gocd-cli environment get gocd_environment_1 --env-var ENVIRONMENT_VAR_1 --env-var ENVIRONMENT_VAR_2 -o yaml
gocd-cli environment get gocd_environment_1 --pipelines -o yaml`,
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := client.GetEnvironment(args[0])
			if err != nil {
				return err
			}

			if fetchPipeline {
				pipelineName := make([]string, 0)

				for _, name := range response.Pipelines {
					pipelineName = append(pipelineName, name.Name)
				}

				return cliRenderer.Render(pipelineName)
			}

			if len(environmentVariables) != 0 {
				envVars := make([]gocd.EnvVars, 0)
				for _, envVar := range response.EnvVars {
					for _, environmentVariable := range environmentVariables {
						if funk.Contains(envVar.Name, environmentVariable) {
							envVars = append(envVars, envVar)
						}
					}
				}

				return cliRenderer.Render(envVars)
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

			return cliRenderer.Render(response)
		},
	}

	getEnvironmentCmd.PersistentFlags().BoolVarP(&fetchPipeline, "pipelines", "", false,
		"when set fetches the pipeline alone")
	getEnvironmentCmd.PersistentFlags().StringSliceVarP(&environmentVariables, "env-var", "", nil,
		"list of environment variables to fetch from the GoCD environment")

	return getEnvironmentCmd
}

func getEnvironmentMapping() *cobra.Command {
	getEnvironmentMappingCmd := &cobra.Command{
		Use:     "get-mappings",
		Short:   "Command to Identify the given environment is part of which config-repo of GoCD",
		Example: "gocd-cli environment get-mappings --environment production -o yaml",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			environmentNames, err := client.GetEnvironmentsMerged(goCDEnvironments)
			if err != nil {
				return err
			}

			cliLogger.Debugf("GoCD environment mapping information was fetched successfully")

			environmentMappings := make([]map[string]string, 0)

			for _, environmentName := range environmentNames {
				environmentMappings = append(environmentMappings, getOriginType(map[string]string{"name": environmentName.Name}, environmentName.Origins))
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := query.SetQuery(environmentMappings, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debug(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(environmentMappings)
		},
	}

	getEnvironmentMappingCmd.PersistentFlags().StringSliceVarP(&goCDEnvironments, "environment", "", nil,
		"name of the environment for which mappings to be fetched")

	return getEnvironmentMappingCmd
}

func createEnvironmentCommand() *cobra.Command {
	createEnvironmentCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE the environment with all specified configuration [https://api.gocd.org/current/#create-an-environment]",
		Example: `gocd-cli environment create gocd_environment_1 --from-file gocd_environment_1.yaml`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE:    createEnvironment,
	}

	return createEnvironmentCmd
}

func updateEnvironmentCommand() *cobra.Command {
	updateEnvironmentCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE the environment with the latest specified configuration [https://api.gocd.org/current/#update-an-environment]",
		Example: `gocd-cli environment update gocd_environment_1 --from-file gocd_environment_1.yaml`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var envs gocd.Environment
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			environmentFetched, err := client.GetEnvironment(envs.Name)
			if err != nil && !strings.Contains(err.Error(), "404") {
				return err
			}

			if create {
				if reflect.DeepEqual(environmentFetched, gocd.Environment{}) {
					return createEnvironment(cmd, nil)
				}
			}

			if len(envs.ETAG) == 0 {
				envs.ETAG = environmentFetched.ETAG
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "environment", envs.Name)

			existing, err := diffCfg.String(environmentFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			env, err := client.UpdateEnvironment(envs)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("environment %s updated successfully", envs.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(env)
		},
	}

	updateEnvironmentCmd.PersistentFlags().BoolVarP(&create, "create", "", false,
		"if a environment by this name doesn't already exist, run create")

	return updateEnvironmentCmd
}

func patchEnvironmentCommand() *cobra.Command {
	patchEnvironmentCmd := &cobra.Command{
		Use:     "patch",
		Short:   "Command to PATCH the environment with the latest specified configuration [https://api.gocd.org/current/#patch-an-environment]",
		Example: `gocd-cli environment patch gocd_environment_1 --from-file gocd_environment_1.yaml`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var envs gocd.Environment
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(patchMessage, "environment", envs.Name)

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

			env, err := client.PatchEnvironment(envs)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("environment %s patched successfully", envs.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(env)
		},
	}

	return patchEnvironmentCmd
}

func deleteEnvironmentCommand() *cobra.Command {
	deleteEnvironmentCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE the specified environment from GoCD [https://api.gocd.org/current/#delete-an-environment]",
		Example: `gocd-cli environment delete gocd_environment_1`,
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			environmentName := args[0]

			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete gocd environment '%s' [y/n]", environmentName)

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

			if err := client.DeleteEnvironment(environmentName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("environment deleted: %s", environmentName))
		},
	}

	return deleteEnvironmentCmd
}

func listEnvironmentsCommand() *cobra.Command {
	listEnvironmentsCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all the environments present in GoCD [https://api.gocd.org/current/#get-all-environments]",
		Example: `gocd-cli environment list -o yaml`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			response, err := client.GetEnvironments()
			if err != nil {
				return err
			}

			if dangling {
				agents, err := client.GetAgents()
				if err != nil {
					return err
				}

				response = funk.Filter(response, func(environment gocd.Environment) bool {
					agentsFound := funk.Find(agents, func(agent gocd.Agent) bool {
						return funk.Find(agent.Environments, environment).(bool)
					}).(bool)

					return len(environment.Pipelines) == 0 || agentsFound
				}).([]gocd.Environment)
			}

			var goCdEnvironments []string

			for _, environment := range response {
				goCdEnvironments = append(goCdEnvironments, environment.Name)
			}

			return cliRenderer.Render(strings.Join(goCdEnvironments, "\n"))
		},
	}

	registerDanglingFlags(listEnvironmentsCmd)

	return listEnvironmentsCmd
}

func getOriginType(mappings map[string]string, origins []gocd.EnvironmentOrigin) map[string]string {
	originTypeConfigRepo := 2
	originTypeGoCD := 1

	switch len(origins) {
	case originTypeConfigRepo:
		mappings["origin_type"] = origins[1].Type
		mappings["origin"] = origins[1].ID
	case originTypeGoCD:
		mappings["origin_type"] = origins[0].Type
	}

	return mappings
}

func createEnvironment(cmd *cobra.Command, _ []string) error {
	var envs gocd.Environment

	object, err := readObject(cmd)
	if err != nil {
		return err
	}

	switch objType := object.CheckFileType(cliLogger); objType {
	case content.FileTypeYAML:
		if err = yaml.Unmarshal([]byte(object), &envs); err != nil {
			return err
		}
	case content.FileTypeJSON:
		if err = json.Unmarshal([]byte(object), &envs); err != nil {
			return err
		}
	default:
		return &errors.UnknownObjectTypeError{Name: objType}
	}

	if err = client.CreateEnvironment(envs); err != nil {
		return err
	}

	return cliRenderer.Render(fmt.Sprintf("environment %s created successfully", envs.Name))
}
