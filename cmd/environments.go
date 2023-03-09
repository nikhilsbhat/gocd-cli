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

func getEnvironmentsCommand() *cobra.Command {
	environmentCommand := &cobra.Command{
		Use:   "environment",
		Short: "Command to operate on environments present in GoCD [https://api.gocd.org/current/#environment-config]",
		Long: `Command leverages GoCD environments apis' [https://api.gocd.org/current/#environment-config] to 
GET/CREATE/UPDATE/PATCH/DELETE GoCD environments`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	environmentCommand.SetUsageTemplate(getUsageTemplate())

	environmentCommand.AddCommand(getGetEnvironmentsCommand())
	environmentCommand.AddCommand(getGetEnvironmentCommand())
	environmentCommand.AddCommand(createEnvironmentCommand())
	environmentCommand.AddCommand(updateEnvironmentCommand())
	environmentCommand.AddCommand(patchEnvironmentCommand())
	environmentCommand.AddCommand(deleteEnvironmentCommand())

	return environmentCommand
}

func getGetEnvironmentsCommand() *cobra.Command {
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

			return cliRenderer.Render(response)
		},
	}

	return getEnvironmentsCmd
}

func getGetEnvironmentCommand() *cobra.Command {
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

			return cliRenderer.Render(response)
		},
	}

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
			case utils.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			case utils.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			if err = client.CreateEnvironment(envs); err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("environment %s created successfully", envs.Name)); err != nil {
				return err
			}

			return nil
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
			case utils.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			case utils.FileTypeJSON:
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

			if err = cliRenderer.Render(env); err != nil {
				return err
			}

			return nil
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
			case utils.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &envs); err != nil {
					return err
				}
			case utils.FileTypeJSON:
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

			if err = cliRenderer.Render(env); err != nil {
				return err
			}

			return nil
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

			if err := cliRenderer.Render(fmt.Sprintf("environment deleted: %s", args[0])); err != nil {
				return err
			}

			return nil
		},
	}

	return deleteEnvironmentCmd
}
