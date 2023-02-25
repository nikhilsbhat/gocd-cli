package cmd

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/utils"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	backupRetry int
	delay       time.Duration
)

func getBackupCommand() *cobra.Command {
	configRepoCommand := &cobra.Command{
		Use:   "backup",
		Short: "Command to operate on backup in GoCD [https://api.gocd.org/current/#backups]",
		Long: `Command leverages GoCD backup apis' [https://api.gocd.org/current/#backups] to 
GET/CREATE/UPDATE/DELETE/SCHEDULE the backup in GoCD server.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}

	configRepoCommand.SetUsageTemplate(getUsageTemplate())

	registerBackupFlags(configRepoCommand)

	configRepoCommand.AddCommand(getBackupConfig())
	configRepoCommand.AddCommand(createOrUpdateBackupConfig())
	configRepoCommand.AddCommand(deleteBackupConfig())
	configRepoCommand.AddCommand(getBackupStats())
	configRepoCommand.AddCommand(scheduleBackup())

	return configRepoCommand
}

func getBackupConfig() *cobra.Command {
	getBackupConfigCommand := &cobra.Command{
		Use:     "get-config",
		Short:   "Command to GET backup config configured in GoCD [https://api.gocd.org/current/#get-backup-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetBackupConfig()
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(response); err != nil {
				return err
			}

			return nil
		},
	}

	return getBackupConfigCommand
}

func createOrUpdateBackupConfig() *cobra.Command {
	createOrUpdateBackupConfigCommand := &cobra.Command{
		Use:     "create-config",
		Short:   "Command to CREATE/UPDATE backup config configured in GoCD [https://api.gocd.org/current/#create-or-update-backup-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var backupConfig gocd.BackupConfig
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case utils.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &backupConfig); err != nil {
					return err
				}
			case utils.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &backupConfig); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			if err = client.CreateOrUpdateBackupConfig(backupConfig); err != nil {
				return err
			}

			cliLogger.Infoln("backup config was created/updated successfully")

			return nil
		},
	}

	return createOrUpdateBackupConfigCommand
}

func deleteBackupConfig() *cobra.Command {
	deleteBackupConfigCommand := &cobra.Command{
		Use:     "delete-config",
		Short:   "Command to DELETE backup config configured in GoCD [https://api.gocd.org/current/#delete-backup-config]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := client.DeleteBackupConfig(); err != nil {
				return err
			}

			cliLogger.Infoln("backup config was deleted successfully")

			return nil
		},
	}

	return deleteBackupConfigCommand
}

func getBackupStats() *cobra.Command {
	getBackupStatsCommand := &cobra.Command{
		Use:     "stats",
		Short:   "Command to GET stats of the specific backup taken in GoCD [https://api.gocd.org/current/#get-backup]",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetBackup(args[0])
			if err != nil {
				return err
			}

			if err = cliRenderer.Render(response); err != nil {
				return err
			}

			return nil
		},
	}

	return getBackupStatsCommand
}

func scheduleBackup() *cobra.Command {
	scheduleBackupCommand := &cobra.Command{
		Use:     "schedule",
		Short:   "Command to SCHEDULE backups in GoCD [https://api.gocd.org/current/#schedule-backup]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.ScheduleBackup()
			if err != nil {
				return err
			}

			retryAfter, err := strconv.Atoi(response["RetryAfter"])
			if err != nil {
				return err
			}

			cliLogger.Debugf("retry after is set by GoCD server to %d, so retrying after %d seconds", retryAfter, retryAfter)

			backupID := response["BackUpID"]
			cliLogger.Debugf("fetching information of the backup id: '%s'", backupID)

			currentRetryCount := 0
			var latestBackupStatus string
			for {
				if currentRetryCount > backupRetry {
					cliLogger.Fatalf("maximum retry count of '%d' crossed with current count '%d', still backup is not ready yet with status '%s'. Exitting",
						backupRetry, currentRetryCount, latestBackupStatus)
				}

				response, err := client.GetBackup(backupID)
				if err != nil {
					return err
				}

				retryRemaining := backupRetry - currentRetryCount
				if response.Status == "IN_PROGRESS" {
					cliLogger.Infof("the backup stats is still in IN_PROGRESS status, retrying... '%d' more to go", retryRemaining)
				}

				if response.Status == "COMPLETED" {
					cliLogger.Debug("the backup is complete, printing backup stats")

					if err = cliRenderer.Render(response); err != nil {
						return err
					}

					break
				}

				if response.Status != "COMPLETED" && response.Status != "IN_PROGRESS" {
					cliLogger.Errorf("looks like backup status is neither IN_PROGRESS nor COMPLETED rather it is %s", response.Status)

					return err
				}

				latestBackupStatus = response.Status
				time.Sleep(delay)
				currentRetryCount++
			}

			return nil
		},
	}

	return scheduleBackupCommand
}
