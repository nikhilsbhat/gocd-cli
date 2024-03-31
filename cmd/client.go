package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/nikhilsbhat/common/diff"
	"github.com/nikhilsbhat/common/renderer"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/utils"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

var (
	client                 gocd.GoCd
	cliRenderer            renderer.Config
	cliShellReadConfig     *utils.ReadConfig
	diffCfg                *diff.Config
	supportedOutputFormats = []string{"yaml", "json", "csv", "table"}
)

func setCLIClient(_ *cobra.Command, _ []string) error {
	SetLogger(cliCfg.LogLevel)

	if localConfig, localConfigPath, err := checkForConfig(); err != nil {
		return err
	} else if localConfig && !cliCfg.skipCacheConfig {
		cliLogger.Debugf("found authorization configuration in cache, loading config from %s", localConfigPath)
		if yamlConfig, err := os.ReadFile(localConfigPath); err != nil {
			return err
		} else if err := yaml.Unmarshal(yamlConfig, &cliCfg); err != nil {
			return err
		}
		cliLogger.Debug("authorization configuration loaded from cache successfully")
	}

	if len(cliCfg.CaPath) != 0 {
		cliLogger.Debug("CA based auth is enabled, hence reading CA from the path")

		if caAbs, err := filepath.Abs(cliCfg.CaPath); err != nil {
			return err
		} else if caContent, err := os.ReadFile(caAbs); err != nil {
			return err
		} else {
			client = gocd.NewClient(cliCfg.URL, cliCfg.Auth, cliCfg.APILogLevel, caContent)
		}
	} else {
		client = gocd.NewClient(cliCfg.URL, cliCfg.Auth, cliCfg.APILogLevel, nil)
	}

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
		errMsg := fmt.Sprintf("unsupported output format '%s', the value should be one of %s", cliCfg.OutputFormat, supportedOutputFormatsString)
		cliLogger.Errorf(errMsg)

		return &errors.CLIError{Message: errMsg}
	}

	diffCfg = diff.NewDiff(cliCfg.OutputFormat, cliCfg.NoColor, cliLogger)

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

func (cfg *Config) GetOutputFormat() string {
	switch {
	case cfg.yaml:
		return "yaml"
	case cfg.json:
		return "json"
	case cfg.table:
		return "table"
	case cfg.csv:
		return "csv"
	default:
		return ""
	}
}
