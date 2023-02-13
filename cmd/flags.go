package cmd

import (
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
	cmd.PersistentFlags().BoolVarP(&cliCfg.JSON, "to-json", "", false,
		"enable this to render output in JSON format")
	cmd.PersistentFlags().BoolVarP(&cliCfg.YAML, "to-yaml", "", false,
		"enable this to render output in YAML format")
	cmd.PersistentFlags().BoolVarP(&cliCfg.saveConfig, "save-config", "", false,
		"enable this to locally save auth configs used to connect GoCD server (path: $HOME/.gocd/auth_config.yaml)")
}

func registerEncryptionFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cipherKey, "cipher-key", "", "",
		"cipher key value used for decryption, the key should same which is used by GoCD server for encryption")
}

func registerConfigRepoFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cipherKey, "from-file", "", "",
		"file containing config repo object that needs to be created")
}
