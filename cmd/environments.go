package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

var (
	fetchPipeline        bool
	environmentVariables []string
)

func registerEnvironmentsCommand() *cobra.Command {
	environmentCommand := &cobra.Command{
		Use:   "environment",
		Short: "Command to operate on environments present in GoCD [https://api.gocd.org/current/#environment-config]",
		Long: `Command leverages GoCD environments apis' [https://api.gocd.org/current/#environment-config] to 
GET/CREATE/UPDATE/PATCH/DELETE and list GoCD environments`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	environmentCommand.SetUsageTemplate(getUsageTemplate())

	environmentCommand.AddCommand(getEnvironmentsCommand())
	environmentCommand.AddCommand(getEnvironmentCommand())
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
		Use:     "get-all",
		Short:   "Command to GET all the environments present in GoCD [https://api.gocd.org/current/#get-all-environments]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
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

	getEnvironmentsCmd.PersistentFlags().BoolVarP(&fetchPipeline, "pipelines", "", false,
		"when set fetches the pipeline alone")
	getEnvironmentsCmd.PersistentFlags().StringSliceVarP(&environmentVariables, "env-var", "", nil,
		"list of environment variables to fetch from the GoCD environment")

	return getEnvironmentsCmd
}

func getEnvironmentCommand() *cobra.Command {
	getEnvironmentCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a specific environments present in GoCD [https://api.gocd.org/current/#get-environment-config]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
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

	getEnvironmentCmd.PersistentFlags().BoolVarP(&fetchPipeline, "pipelines", "", false,
		"when set fetches the pipeline alone")
	getEnvironmentCmd.PersistentFlags().StringSliceVarP(&environmentVariables, "env-var", "", nil,
		"list of environment variables to fetch from the GoCD environment")

	return getEnvironmentCmd
}

func createEnvironmentCommand() *cobra.Command {
	createEnvironmentCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE the environment with all specified configuration [https://api.gocd.org/current/#create-an-environment]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var envs gocd.Environment
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			case render.FileTypeJSON:
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
		},
	}

	return createEnvironmentCmd
}

func updateEnvironmentCommand() *cobra.Command {
	updateEnvironmentCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE the environment with the latest specified configuration [https://api.gocd.org/current/#update-an-environment]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var envs gocd.Environment
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			fmt.Println(envs)

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

	return updateEnvironmentCmd
}

func patchEnvironmentCommand() *cobra.Command {
	patchEnvironmentCmd := &cobra.Command{
		Use:     "patch",
		Short:   "Command to PATCH the environment with the latest specified configuration [https://api.gocd.org/current/#patch-an-environment]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var envs gocd.Environment
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
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
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteEnvironment(args[0]); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("environment deleted: %s", args[0]))
		},
	}

	return deleteEnvironmentCmd
}

func listEnvironmentsCommand() *cobra.Command {
	listEnvironmentsCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all the environments present in GoCD [https://api.gocd.org/current/#get-all-environments]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetEnvironments()
			if err != nil {
				return err
			}

			var goCdEnvironments []string

			for _, environment := range response {
				goCdEnvironments = append(goCdEnvironments, environment.Name)
			}

			return cliRenderer.Render(strings.Join(goCdEnvironments, "\n"))
		},
	}

	return listEnvironmentsCmd
}
