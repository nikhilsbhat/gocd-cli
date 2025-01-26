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
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func registerUsersCommand() *cobra.Command {
	usersCommand := &cobra.Command{
		Use:   "user",
		Short: "Command to operate on users in GoCD [https://api.gocd.org/current/#users]",
		Long: `Command leverages GoCD users apis' [https://api.gocd.org/current/#users] to 
GET/CREATE/UPDATE/DELETE/BULK-DELETE/BULK-UPDATE the users in GoCD server.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Usage()
		},
	}

	usersCommand.SetUsageTemplate(getUsageTemplate())

	usersCommand.AddCommand(usersGetCommand())
	usersCommand.AddCommand(userGetCommand())
	usersCommand.AddCommand(userCreateCommand())
	usersCommand.AddCommand(userUpdateCommand())
	usersCommand.AddCommand(userDeleteCommand())
	usersCommand.AddCommand(bulkDeleteUsersCommand())
	usersCommand.AddCommand(listUsersCommand())

	for _, command := range usersCommand.Commands() {
		command.SilenceUsage = true
	}

	return usersCommand
}

func usersGetCommand() *cobra.Command {
	getUsersCmd := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all users present in GoCD [https://api.gocd.org/current/#get-all-users]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetUsers()
				if err != nil {
					return err
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

	return getUsersCmd
}

func userGetCommand() *cobra.Command {
	getUserCmd := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET user present in GoCD [https://api.gocd.org/current/#get-one-user]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			for {
				response, err := client.GetUser(args[0])
				if err != nil {
					return err
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

	return getUserCmd
}

func userCreateCommand() *cobra.Command {
	createUserCmd := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE user in GoCD [https://api.gocd.org/current/#create-a-user]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE:    userCreate,
	}

	return createUserCmd
}

func userUpdateCommand() *cobra.Command {
	updateUserCmd := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE user present in GoCD [https://api.gocd.org/current/#update-a-user]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var user gocd.User
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &user); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &user); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			userFetched, err := client.GetUser(user.Name)
			if err != nil && !strings.Contains(err.Error(), "404") {
				return err
			}

			if create {
				if reflect.DeepEqual(userFetched, gocd.User{}) {
					return createRole(cmd, nil)
				}
			}

			cliShellReadConfig.ShellMessage = fmt.Sprintf(updateMessage, "user", user.Name)

			existing, err := diffCfg.String(userFetched)
			if err != nil {
				return err
			}

			if err = cliCfg.CheckDiffAndAllow(existing, object.String()); err != nil {
				return err
			}

			_, err = client.UpdateUser(user)
			if err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("user %s updated successfully", user.Name))
		},
	}

	updateUserCmd.PersistentFlags().BoolVarP(&create, "create", "", false,
		"if a user by this name doesn't already exist, run create")

	return updateUserCmd
}

func userDeleteCommand() *cobra.Command {
	deleteUserCmd := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE user present in GoCD [https://api.gocd.org/current/#delete-a-user]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, args []string) error {
			userName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete user '%s' [y/n]", userName)

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

			if err := client.DeleteUser(userName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("user deleted: %s", userName))
		},
	}

	return deleteUserCmd
}

func bulkDeleteUsersCommand() *cobra.Command {
	bulkDeleteUserCmd := &cobra.Command{
		Use:     "delete-bulk",
		Short:   "Command to BULK-DELETE users present in GoCD [https://api.gocd.org/current/#bulk-delete-users]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var user map[string]interface{}
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case content.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &user); err != nil {
					return err
				}
			case content.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &user); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			if err = client.BulkDeleteUsers(user); err != nil {
				return err
			}

			return cliRenderer.Render("users deleted in bulk")
		},
	}

	return bulkDeleteUserCmd
}

func listUsersCommand() *cobra.Command {
	listPluginsCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all users present in GoCD",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(_ *cobra.Command, _ []string) error {
			for {
				response, err := client.GetUsers()
				if err != nil {
					return err
				}

				var usersList []string

				for _, user := range response {
					usersList = append(usersList, user.Name)
				}

				if err = cliRenderer.Render(strings.Join(usersList, "\n")); err != nil {
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

	listPluginsCmd.SetUsageTemplate(getUsageTemplate())

	return listPluginsCmd
}

func userCreate(cmd *cobra.Command, _ []string) error {
	var user gocd.User

	object, err := readObject(cmd)
	if err != nil {
		return err
	}

	switch objType := object.CheckFileType(cliLogger); objType {
	case content.FileTypeYAML:
		if err = yaml.Unmarshal([]byte(object), &user); err != nil {
			return err
		}
	case content.FileTypeJSON:
		if err = json.Unmarshal([]byte(object), &user); err != nil {
			return err
		}
	default:
		return &errors.UnknownObjectTypeError{Name: objType}
	}

	_, err = client.CreateUser(user)
	if err != nil {
		return err
	}

	return cliRenderer.Render(fmt.Sprintf("user %s created successfully", user.Name))
}
