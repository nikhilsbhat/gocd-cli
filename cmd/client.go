package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-cli/pkg/utils"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	client      gocd.GoCd
	cliRenderer utils.Renderer
)

func setCLIClient(cmd *cobra.Command, args []string) error {
	var caContent []byte

	SetLogger(cliCfg.LogLevel)

	localConfig, localConfigPath, err := checkForConfig()
	if err != nil {
		return err
	}

	if localConfig && !cliCfg.skipCacheConfig {
		cliLogger.Debugf("found authorisation configuration in cache, loading config from %s", localConfigPath)
		yamlConfig, err := os.ReadFile(localConfigPath)
		if err != nil {
			return err
		}
		if err = yaml.Unmarshal(yamlConfig, &cliCfg); err != nil {
			return err
		}
		cliLogger.Debug("authorisation configuration loaded from cache successfully")
	}

	if len(cliCfg.CaPath) != 0 {
		cliLogger.Debug("CA based auth is enabled, hence reading ca from the path")
		caAbs, err := filepath.Abs(cliCfg.CaPath)
		if err != nil {
			return err
		}

		caContent, err = os.ReadFile(caAbs)
		if err != nil {
			log.Fatal(err)
		}
	}

	if cliCfg.saveConfig {
		cliLogger.Debug("--save-config is enabled, hence saving authorisation configuration")
		if err = cliCfg.saveAuthConfig(); err != nil {
			return err
		}
	}

	goCDClient := gocd.NewClient(
		cliCfg.URL,
		cliCfg.Auth,
		cliCfg.LogLevel,
		caContent,
	)

	client = goCDClient

	writer := os.Stdout
	if len(cliCfg.ToFile) != 0 {
		cliLogger.Debugf("--to-file is opted, output would be saved under a file '%s'", cliCfg.ToFile)
		filePTR, err := os.Create(cliCfg.ToFile)
		if err != nil {
			return err
		}
		writer = filePTR
	}

	cliRenderer = utils.GetRenderer(writer, cliLogger, cliCfg.YAML, cliCfg.JSON)

	return nil
}

func (cfg *Config) saveAuthConfig() error {
	cliLogger.Debug("saving authorisation config to cache, so that it can be reused next")
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	authConfigDir := filepath.Join(home, ".gocd")
	const filePermission = 0o644
	if err = os.Mkdir(authConfigDir, filePermission); os.IsNotExist(err) {
		return err
	}

	authConfigFile, err := os.Create(filepath.Join(authConfigDir, "auth_config.yaml"))
	if err != nil {
		return err
	}

	cliLogger.Debugf("authorisation config would be saved under %s", authConfigFile.Name())

	cfgYAML, err := yaml.Marshal(cliCfg)
	if err != nil {
		return err
	}

	if _, err = authConfigFile.WriteString(string(cfgYAML)); err != nil {
		return err
	}

	return nil
}

func checkForConfig() (bool, string, error) {
	cliLogger.Debug("searching for authorisation configuration in cache")
	home, err := os.UserHomeDir()
	if err != nil {
		return false, "", err
	}

	authConfigDir := filepath.Join(home, ".gocd")
	configPath := filepath.Join(authConfigDir, "auth_config.yaml")

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		cliLogger.Debug("no authorisation configuration found in cache")

		return false, "", nil
	}

	return true, configPath, nil
}
