package diff

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/ghodss/yaml"
	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/sirupsen/logrus"
)

const (
	contextLines = 2000000
)

type Config struct {
	NoColor bool
	Format  string
	log     *logrus.Logger
}

func (cfg *Config) Diff(oldData, newData string) (bool, string, error) {
	switch cfg.Format {
	case "yaml":
		cfg.log.Debug("loading diff in yaml format")
	case "json":
		cfg.log.Debug("loading diff in json format")
	default:
		return false, "", &errors.CLIError{Message: fmt.Sprintf("unknown format, cannot calculate diff for the format '%s'", cfg.Format)}
	}

	diffIdentified, err := diff(oldData, newData)
	if err != nil {
		return false, "", err
	}

	if len(diffIdentified) == 0 {
		return false, "", nil
	}

	return true, strings.Join(diffIdentified, "\n"), nil
}

func (cfg *Config) String(input interface{}) (string, error) {
	switch strings.ToLower(cfg.Format) {
	case "yaml":
		out, err := yaml.Marshal(input)
		if err != nil {
			return "", err
		}

		yamlString := strings.Join([]string{"---", string(out)}, "\n")

		return yamlString, nil
	case "json":
		out, err := json.MarshalIndent(input, "", "     ")
		if err != nil {
			return "", err
		}

		return string(out), nil
	default:
		return "", &errors.CLIError{
			Message: fmt.Sprintf("type '%s' is not supported for loading diff", cfg.Format),
		}
	}
}

func (cfg *Config) SetLogger(log *logrus.Logger) {
	cfg.log = log
}

func diff(content1, content2 string) ([]string, error) {
	lines := make([]string, 0)
	diffVal := difflib.UnifiedDiff{
		A:        difflib.SplitLines(content1),
		B:        difflib.SplitLines(content2),
		FromFile: "old",
		ToFile:   "new",
		Context:  contextLines,
	}

	text, err := difflib.GetUnifiedDiffString(diffVal)
	if err != nil {
		return nil, err
	}

	if len(text) == 0 {
		return lines, nil
	}

	lines = strings.Split(text, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "-") {
			lines[i] = color.RedString(line)
		} else if strings.HasPrefix(line, "+") {
			lines[i] = color.GreenString(line)
		}
	}

	return lines, nil
}
