package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	client      gocd.GoCd
	cliRenderer render.Renderer
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

	cliRenderer = render.GetRenderer(writer, cliLogger, cliCfg.YAML, cliCfg.JSON)

	return nil
}
