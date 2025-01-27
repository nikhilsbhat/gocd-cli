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
	"gopkg.in/yaml.v3"
)

func registerAuthorizationConfigCommand() *cobra.Command {
	registerAuthorizationConfigCmd := &cobra.Command{
		Use:   "authorization",
		Short: "Command to operate on authorization-configuration present in GoCD [https://api.gocd.org/current/#authorization-configuration]",
		Long: `Command leverages GoCD authorization-configuration apis' [https://api.gocd.org/current/#authorization-configuration] to 
GET/CREATE/UPDATE/DELETE cluster profiles present in GoCD`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	registerAuthorizationConfigCmd.SetUsageTemplate(getUsageTemplate())

	registerAuthorizationConfigCmd.AddCommand(getAuthConfigCommand())
	registerAuthorizationConfigCmd.AddCommand(getAuthConfigsCommand())
	registerAuthorizationConfigCmd.AddCommand(getCreateAuthConfigCommand())
	registerAuthorizationConfigCmd.AddCommand(getUpdateAuthConfigCommand())
	registerAuthorizationConfigCmd.AddCommand(getDeleteAuthConfigCommand())
	registerAuthorizationConfigCmd.AddCommand(listAuthConfigsCommand())

	for _, command := range registerAuthorizationConfigCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerAuthorizationConfigCmd
}

func getAuthConfigCommand() *cobra.Command {
	authConfigGetCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a authorization configuration with all specified configurations in GoCD [https://api.gocd.org/current/#get-an-authorization-configuration]",
		Example: "gocd-cli authorization get ldap -o yaml",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			for {
				response, err := client.GetAuthConfig(args[0])
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

	return authConfigGetCmd
}

func getAuthConfigsCommand() *cobra.Command {
	authConfigsGetCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all authorization configurations present in GoCD [https://api.gocd.org/current/#get-all-authorization-configurations]",
		Example: "gocd-cli authorization get-all",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			response, err := client.GetAuthConfigs()
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

			return cliRenderer.Render(response)
		},
	}

	return authConfigsGetCmd
}

func getCreateAuthConfigCommand() *cobra.Command {
	createAuthConfigCommand := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE the authorization configuration with specified configuration [https://api.gocd.org/current/#create-an-authorization-configuration]",
		Example: `gocd-cli authorization create --from-file config-repo.yaml -o yaml`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE:    createAuthConfig,
	}

	createAuthConfigCommand.SetUsageTemplate(getUsageTemplate())

	return createAuthConfigCommand
}

func getUpdateAuthConfigCommand() *cobra.Command {
	authConfigUpdateCommand := &cobra.Command{
		Use:   "update",
		Short: "Command to UPDATE the authorization configuration present in GoCD [https://api.gocd.org/current/#update-an-authorization-configuration]",
		Example: `gocd-cli authorization update --from-file config-repo.yaml -o yaml
gocd-cli authorization update --from-file config-repo.yaml --create -o yaml
gocd-cli authorization update --from-file config-repo.yaml --create -o yaml -y`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var authConfig gocd.CommonConfig
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &authConfig); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &authConfig); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			authConfigFetched, err := client.GetAuthConfig(authConfig.ID)
			if err != nil && !strings.Contains(err.Error(), "404") {
				return err
			}

			if create {
				if reflect.DeepEqual(authConfigFetched, gocd.CommonConfig{}) {
					return createAuthConfig(cmd, nil)
				}
			}

			if len(authConfig.ETAG) == 0 {
				authConfig.ETAG = authConfigFetched.ETAG
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "authorization", authConfigFetched.ID)

			existing, err := diffCfg.String(authConfigFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			response, err := client.UpdateAuthConfig(authConfig)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	authConfigUpdateCommand.SetUsageTemplate(getUsageTemplate())
	authConfigUpdateCommand.PersistentFlags().BoolVarP(&create, "create", "", false,
		"if a config repo by this name doesn't already exist, run create")

	return authConfigUpdateCommand
}

func getDeleteAuthConfigCommand() *cobra.Command {
	deleteAuthConfigCmd := &cobra.Command{
		Use:   "delete",
		Short: "Command to DELETE the specified authorization configuration present in GoCD [https://api.gocd.org/current/#delete-an-authorization-configuration]",
		Example: `gocd-cli authorization delete helm-images
gocd-cli authorization delete helm-images -y`,
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			AuthConfigName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete authorization-configuration repo '%s' [y/n]", AuthConfigName)

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

			if err := client.DeleteAuthConfig(AuthConfigName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("authorization configuration deleted: %s", AuthConfigName))
		},
	}

	return deleteAuthConfigCmd
}

func listAuthConfigsCommand() *cobra.Command {
	listAuthConfigsCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all authorization configurations present in GoCD [https://api.gocd.org/current/#get-all-authorization-configurations]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetAuthConfigs()
				if err != nil {
					return err
				}

				var authConfigs []string

				for _, commonConfig := range response {
					authConfigs = append(authConfigs, commonConfig.ID)
				}

				if err = cliRenderer.Render(strings.Join(authConfigs, "\n")); err != nil {
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

	return listAuthConfigsCmd
}

func createAuthConfig(cmd *cobra.Command, _ []string) error {
	var authConfig gocd.CommonConfig

	object, err := readObject(cmd)
	if err != nil {
		return err
	}

	fmt.Println(object.String())

	switch objType := object.CheckFileType(cliLogger); objType {
	case content.FileTypeYAML:
		if err = yaml.Unmarshal([]byte(object), &authConfig); err != nil {
			return err
		}
	case content.FileTypeJSON:
		if err = json.Unmarshal([]byte(object), &authConfig); err != nil {
			return err
		}
	default:
		return &errors.UnknownObjectTypeError{Name: objType}
	}

	if _, err = client.CreateAuthConfig(authConfig); err != nil {
		return err
	}

	return cliRenderer.Render(fmt.Sprintf("config repo %s created successfully", authConfig.ID))
}
