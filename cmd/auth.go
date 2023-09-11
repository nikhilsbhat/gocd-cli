package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	goCdCacheDirName       = ".gocd"
	goCdAuthConfigFileName = "auth_config.yaml"
)

func registerAuthConfigCommand() *cobra.Command {
	registerAuthConfigCmd := &cobra.Command{
		Use:   "auth-config",
		Short: "Command to store/remove the authorization configuration to be used by the cli",
		Long: `Using the auth config commands, one can cache the authorization configuration onto a file so it can be used by further calls made using this utility.
Also, the cached authentication configurations can be erased using the same`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	registerAuthConfigCmd.SetUsageTemplate(getUsageTemplate())

	registerAuthConfigCmd.AddCommand(getAuthStoreCommand())
	registerAuthConfigCmd.AddCommand(getAuthEraseCommand())

	for _, command := range registerAuthConfigCmd.Commands() {
		command.SilenceUsage = true
	}

	return registerAuthConfigCmd
}

func getAuthStoreCommand() *cobra.Command {
	authStoreCmd := &cobra.Command{
		Use:     "store",
		Short:   "Command to cache the GoCD authorization configuration to be used by the cli",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliLogger.Debug("saving authorisation config to cache, so that it can be reused next")
			home, err := os.UserHomeDir()
			if err != nil {
				cliLogger.Errorf("fetching user's home directory errored with '%v'", err)

				return err
			}

			authConfigDir := filepath.Join(home, goCdCacheDirName)
			const filePermission = 0o644
			if err = os.Mkdir(authConfigDir, filePermission); os.IsNotExist(err) {
				cliLogger.Errorf("creating directory '%s' errored with '%v'", authConfigDir, err)

				return err
			}

			authFile := filepath.Join(authConfigDir, goCdAuthConfigFileName)
			authConfigFile, err := os.Create(authFile)
			if err != nil {
				cliLogger.Errorf("creating authfile '%s' errored with '%v'", authFile, err)

				return err
			}

			cliLogger.Infof("authorisation config would be saved under %s", authConfigFile.Name())

			cfgYAML, err := yaml.Marshal(cliCfg)
			if err != nil {
				cliLogger.Errorf("serializing auth config data to yaml errored with '%s'", err)

				return err
			}

			//nolint:mirror
			if _, err = authConfigFile.WriteString(string(cfgYAML)); err != nil {
				cliLogger.Errorf("writing auth config data to file '%s' errored with '%s'", authConfigFile.Name(), err)

				return err
			}

			cliLogger.Infof("authorisation config was successfully saved under %s", authConfigFile.Name())

			return nil
		},
	}

	return authStoreCmd
}

func getAuthEraseCommand() *cobra.Command {
	authEraseCmd := &cobra.Command{
		Use:     "remove",
		Short:   "Command to remove the cached GoCD authorization configuration that is used by the cli.",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			home, err := os.UserHomeDir()
			if err != nil {
				cliLogger.Errorf("fetching user's home directory errored with '%v'", err)

				return err
			}

			authConfigFile := filepath.Join(home, goCdCacheDirName, goCdAuthConfigFileName)

			cliLogger.Infof("authorisation config saved in '%s' would be cleaned", authConfigFile)

			if err = os.RemoveAll(authConfigFile); err != nil {
				cliLogger.Errorf("cleaning authorisation config saved in '%s' errored with '%v'", authConfigFile, err)

				return err
			}

			cliLogger.Infof("authorisation config saved in '%s' was cleaned successfully", authConfigFile)

			return nil
		},
	}

	return authEraseCmd
}

func checkForConfig() (bool, string, error) {
	cliLogger.Debug("searching for authorisation configuration in cache")
	home, err := os.UserHomeDir()
	if err != nil {
		return false, "", err
	}

	configPath := filepath.Join(home, goCdCacheDirName, goCdAuthConfigFileName)

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		cliLogger.Debug("no authorisation configuration found in cache")

		return false, "", nil
	}

	return true, configPath, nil
}
