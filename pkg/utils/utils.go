package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
)

// ReadConfig holds the necessary inputs that is required by Reader.
type ReadConfig struct {
	ShellName    string    `json:"shell_name,omitempty" yaml:"shell_name,omitempty"`
	ShellMessage string    `json:"shell_message,omitempty" yaml:"shell_message,omitempty"`
	InputOptions []Options `json:"inputs,omitempty" yaml:"inputs,omitempty"`
	logger       *logrus.Logger
}

// Options that should be considered while configuring the shell reader.
// ex: yes/no (name) (y/n) (short)
// To get about value the options would look like: []Options{{Name: yes,Short: y}, {Name: no,Short: n}}.
type Options struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	Short string `json:"short,omitempty" yaml:"short,omitempty"`
}

// Option implements method that helps in identifying if input is present in predefined list.
type Option []Options

// Reader reads the inputs from the shell input and matches against the inputs set.
// This would help one in designing the CLI commands that interactively takes input from end user ang validate them.
// For example: taking inputs such as yes or no from end user.
func (cfg *ReadConfig) Reader() (bool, Options) {
	shellReader := bufio.NewReader(os.Stdin)
	flattenedInputs := make([]string, 0)

	funk.ForEach(cfg.InputOptions, func(inputOption Options) {
		flattenedInputs = append(flattenedInputs, inputOption.Name, inputOption.Short)
	})

	for {
		fmt.Printf("$%s>> ", cfg.ShellName)
		fmt.Printf("%s: ", cfg.ShellMessage)

		inputStringRaw, err := shellReader.ReadString('\n')
		if err != nil {
			return false, Options{}
		}

		inputString := strings.TrimSpace(inputStringRaw)
		inputLength := len(inputString)

		switch {
		case inputLength == 0:
			cfg.logger.Warnf("did not get any input, please pass one of the valid inputs '%s'", flattenedInputs)

			return false, Options{}
		default:
			cfg.logger.Debug("inputs are matching the required count, proceeding further")
		}

		inputs := getArrayOfInputs(inputString)

		if contains, option := Option(cfg.InputOptions).Contains(inputs[0]); contains {
			return true, option
		}

		cfg.logger.Errorf("the options should be one of '%s'", funk.FlattenDeep(cfg.InputOptions))
	}
}

// Contains checks if user passed input is part of predefined Options.
// If yes returns true else returns false.
func (inputOptions Option) Contains(input string) (bool, Options) {
	var option Options

	return funk.Contains(inputOptions, func(inputOption Options) bool {
		if inputOption.Name == input || inputOption.Short == input {
			option = inputOption

			return true
		}

		return false
	}), option
}

func getArrayOfInputs(inputs string) []string {
	inputs = strings.TrimSuffix(inputs, "\n")

	return strings.Fields(inputs)
}

// NewReadConfig returns new instance of ReadConfig.
func NewReadConfig(name, message string, options []Options, logger *logrus.Logger) *ReadConfig {
	return &ReadConfig{
		ShellName:    name,
		ShellMessage: message,
		InputOptions: options,
		logger:       logger,
	}
}
