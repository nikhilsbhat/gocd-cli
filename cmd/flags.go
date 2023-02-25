package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var cliCfg Config

func registerGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cliCfg.URL, "server-url", "", "http://localhost:8153/go",
		"GoCD server URL base path")
	cmd.PersistentFlags().StringVarP(&cliCfg.Auth.UserName, "username", "u", "",
		"username to authenticate with GoCD server")
	cmd.PersistentFlags().StringVarP(&cliCfg.Auth.Password, "password", "p", "",
		"password to authenticate with GoCD server")
	cmd.PersistentFlags().StringVarP(&cliCfg.Auth.BearerToken, "auth-token", "t", "",
		"token to authenticate with GoCD server, should not be co-used with basic auth (username/password)")
	cmd.PersistentFlags().StringVarP(&cliCfg.CaPath, "ca-file-path", "", "",
		"path to file containing CA cert used to authenticate GoCD server, if you have one")
	cmd.PersistentFlags().StringVarP(&cliCfg.LogLevel, "log-level", "l", "info",
		"log level for gocd cli, log levels supported by [https://github.com/sirupsen/logrus] will work")
	cmd.PersistentFlags().BoolVarP(&cliCfg.JSON, "json", "", false,
		"enable this to Render output in JSON format")
	cmd.PersistentFlags().BoolVarP(&cliCfg.YAML, "yaml", "", false,
		"enable this to Render output in YAML format")
	cmd.PersistentFlags().BoolVarP(&cliCfg.YAML, "no-color", "", false,
		"enable this to Render output in YAML format")
	cmd.PersistentFlags().BoolVarP(&cliCfg.saveConfig, "save-config", "", false,
		"enable this to locally save auth configs used to connect GoCD server (path: $HOME/.gocd/auth_config.yaml)")
	cmd.PersistentFlags().BoolVarP(&cliCfg.skipCacheConfig, "skip-cache-config", "", false,
		"if enabled locally save auth configs would not be used to authenticate GoCD server (path: $HOME/.gocd/auth_config.yaml)")
	cmd.PersistentFlags().StringVarP(&cliCfg.FromFile, "from-file", "f", "",
		"file containing configurations of objects that needs to be created in GoCD, config-repo/pipeline-group/environment and etc.")
	cmd.PersistentFlags().StringVarP(&cliCfg.ToFile, "to-file", "", "",
		"file to which the output needs to be written to (this works only if --yaml or --json is enabled)")
}

func registerEncryptionFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cipherKey, "cipher-key", "", "",
		"cipher key value used for decryption, the key should same which is used by GoCD server for encryption")
}

const (
	defaultRetryCount = 30
	defaultDelay      = 5 * time.Second
)

func registerBackupFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&backupRetry, "retry", "", defaultRetryCount,
		"number of times to retry to get backup stats when backup status is not ready")
	cmd.PersistentFlags().DurationVarP(&delay, "delay", "", defaultDelay,
		"time delay between each retries that would be made to get backup stats")
}
