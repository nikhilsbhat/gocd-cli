package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/nikhilsbhat/common/content"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/query"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
)

func registerRolesCommand() *cobra.Command {
	environmentCommand := &cobra.Command{
		Use:   "roles",
		Short: "Command to operate on roles present in GoCD [https://api.gocd.org/current/#roles]",
		Long: `Command leverages GoCD environments apis' [https://api.gocd.org/current/#roles] to 
GET/CREATE/UPDATE/DELETE and list GoCD roles`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	environmentCommand.SetUsageTemplate(getUsageTemplate())

	environmentCommand.AddCommand(getRolesCommand())
	environmentCommand.AddCommand(getRoleCommand())
	environmentCommand.AddCommand(createRoleCommand())
	environmentCommand.AddCommand(updateRoleCommand())
	environmentCommand.AddCommand(deleteRoleCommand())
	environmentCommand.AddCommand(listRolesCommand())

	for _, command := range environmentCommand.Commands() {
		command.SilenceUsage = true
	}

	return environmentCommand
}

func getRolesCommand() *cobra.Command {
	var roleType string

	getRoleCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all the roles present in GoCD [https://api.gocd.org/current/#get-all-roles]",
		Example: "gocd-cli roles get-all",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				var rolesConfig gocd.RolesConfig
				switch len(roleType) {
				case 0:
					response, err := client.GetRoles()
					if err != nil {
						return err
					}

					rolesConfig = response
				default:
					cliLogger.Debugf("fetching roles by type '%s'", roleType)

					response, err := client.GetRolesByType(roleType)
					if err != nil {
						return err
					}

					rolesConfig = response
				}

				if len(jsonQuery) != 0 {
					cliLogger.Debugf(queryEnabledMessage, jsonQuery)

					baseQuery, err := query.SetQuery(rolesConfig, jsonQuery)
					if err != nil {
						return err
					}

					cliLogger.Debug(baseQuery.Print())

					return cliRenderer.Render(baseQuery.RunQuery())
				}

				if err := cliRenderer.Render(rolesConfig); err != nil {
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

	getRoleCmd.PersistentFlags().StringVarP(&roleType, "type", "", "",
		"type of role to be fetched, ex: gocd, plugin")

	return getRoleCmd
}

func getRoleCommand() *cobra.Command {
	getRoleCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET a specific role in GoCD [https://api.gocd.org/current/#get-a-role]",
		Example: "gocd-cli role get sample-config",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			for {
				response, err := client.GetRole(args[0])
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

	return getRoleCmd
}

func createRoleCommand() *cobra.Command {
	createRoleCmd := &cobra.Command{
		Use: "create",
		Short: "Command to CREATE a role with all specified configurations in GoCD " +
			"[https://api.gocd.org/current/#create-a-gocd-role, https://api.gocd.org/current/#create-a-plugin-role]",
		Example: "gocd-cli role create sample-config --from-file sample-config.yaml --log-level debug",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE:    createRole,
	}

	return createRoleCmd
}

func updateRoleCommand() *cobra.Command {
	updateRoleCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE a role with all specified configurations in GoCD [https://api.gocd.org/current/#update-a-role]",
		Example: "gocd-cli role update sample-config --from-file sample-config.yaml --log-level debug",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var roleCfg gocd.Role
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &roleCfg); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &roleCfg); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			rolesFetched, err := client.GetRole(roleCfg.Name)
			if err != nil && !strings.Contains(err.Error(), "404") {
				return err
			}

			if create {
				if reflect.DeepEqual(rolesFetched, gocd.Role{}) {
					return createRole(cmd, nil)
				}
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "role", rolesFetched.Name)

			existing, err := diffCfg.String(rolesFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			response, err := client.UpdateRole(roleCfg)
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(fmt.Sprintf("role %s updated successfully", roleCfg.Name)); err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	updateRoleCmd.PersistentFlags().BoolVarP(&create, "create", "", false,
		"if a role by this name doesn't already exist, run create")

	return updateRoleCmd
}

func deleteRoleCommand() *cobra.Command {
	deleteRoleCmd := &cobra.Command{
		Use:   "delete",
		Short: "Command to DELETE a specific role present in GoCD [https://api.gocd.org/current/#delete-a-role]",
		Example: `gocd-cli role delete sample-config
gocd-cli role delete sample-config -y`,
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			roleName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete role '%s' [y/n]", roleName)

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

			if err := client.DeleteRole(roleName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("role '%s' deleted successfully", roleName))
		},
	}

	return deleteRoleCmd
}

func listRolesCommand() *cobra.Command {
	listRolesCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all roles present in GoCD [https://api.gocd.org/current/#get-all-roles]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetRoles()
				if err != nil {
					return err
				}

				var elasticAgentProfiles []string

				for _, commonConfig := range response.Role {
					elasticAgentProfiles = append(elasticAgentProfiles, commonConfig.Name)
				}

				if err = cliRenderer.Render(strings.Join(elasticAgentProfiles, "\n")); err != nil {
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

	return listRolesCmd
}

func createRole(cmd *cobra.Command, _ []string) error {
	var roleCfg gocd.Role

	object, err := readObject(cmd)
	if err != nil {
		return err
	}

	switch objType := object.CheckFileType(cliLogger); objType {
	case content.FileTypeYAML:
		if err = yaml.Unmarshal([]byte(object), &roleCfg); err != nil {
			return err
		}
	case content.FileTypeJSON:
		if err = json.Unmarshal([]byte(object), &roleCfg); err != nil {
			return err
		}
	default:
		return &errors.UnknownObjectTypeError{Name: objType}
	}

	response, err := client.CreateRole(roleCfg)
	if err != nil {
		return err
	}

	if err = cliRenderer.Render(fmt.Sprintf("role %s created successfully", roleCfg.Name)); err != nil {
		return err
	}

	return cliRenderer.Render(response)
}
