package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nikhilsbhat/common/renderer"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/utils"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

var (
	client                 gocd.GoCd
	cliRenderer            renderer.Config
	cliShellReadConfig     *utils.ReadConfig
	supportedOutputFormats = []string{"yaml", "json", "csv", "table"}
)

func setCLIClient(_ *cobra.Command, _ []string) error {
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
		cliCfg.APILogLevel,
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

	if !cliCfg.validateOutputFormats() {
		supportedOutputFormatsString := strings.Join(supportedOutputFormats, "|")
		cliLogger.Errorf("unsupported output format '%s', the value should be one of %s",
			cliCfg.OutputFormat, supportedOutputFormatsString)

		return &errors.CLIError{
			Message: fmt.Sprintf("unsupported output format '%s', the value should be one of %s",
				cliCfg.OutputFormat, supportedOutputFormatsString),
		}
	}

	cliCfg.setOutputFormats()

	cliRenderer = renderer.GetRenderer(writer, cliLogger, cliCfg.NoColor, cliCfg.yaml, cliCfg.json, cliCfg.csv, cliCfg.table)

	inputOptions := []utils.Options{{Name: "yes", Short: "y"}, {Name: "no", Short: "n"}}
	cliShellReadConfig = utils.NewReadConfig("gocd-cli", "", inputOptions, cliLogger)

	return nil
}

func (cfg *Config) validateOutputFormats() bool {
	if len(cfg.OutputFormat) == 0 {
		return true
	}

	return funk.Contains(supportedOutputFormats, strings.ToLower(cfg.OutputFormat))
}

func (cfg *Config) setOutputFormats() {
	switch strings.ToLower(cfg.OutputFormat) {
	case "yaml":
		cfg.yaml = true
	case "json":
		cfg.json = true
	case "csv":
		cfg.csv = true
	case "table":
		cfg.table = true
	default:
	}
}
