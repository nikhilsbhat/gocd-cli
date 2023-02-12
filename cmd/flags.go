package cmd

import (
	"github.com/spf13/cobra"
)

var cliCfg Config

func registerGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cliCfg.url, "server-url", "", "http://localhost:8153/go",
		"GoCD server URL base path defaults to (http://localhost:8153/go)")
	cmd.PersistentFlags().StringVarP(&cliCfg.auth.UserName, "username", "u", "",
		"username to authenticate with GoCD server")
	cmd.PersistentFlags().StringVarP(&cliCfg.auth.Password, "password", "p", "",
		"password to authenticate with GoCD server")
	cmd.PersistentFlags().StringVarP(&cliCfg.auth.BearerToken, "auth-token", "t", "",
		"token to authenticate with GoCD server, should not be co-used with basic auth (username/epassword)")
	cmd.PersistentFlags().StringVarP(&cliCfg.caPath, "ca-file-path", "", "",
		"path to file containing CA cert used to authenticate GoCD server, if you have one")
	cmd.PersistentFlags().StringVarP(&cliCfg.loglevel, "log-level", "l", "info",
		"log level for gocd cli (defaults to info), log levels supported by [https://github.com/sirupsen/logrus] will work")
	cmd.PersistentFlags().BoolVarP(&cliCfg.json, "to-json", "", false,
		"enable this to render output in JSON format")
	cmd.PersistentFlags().BoolVarP(&cliCfg.yaml, "to-yaml", "", false,
		"enable this to render output in YAML format")
}

func registerEncryptionFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cipherKey, "cipher-key", "", "",
		"cipher key value used for decryption, the key should same which is used by GoCD server for encryption")
}
