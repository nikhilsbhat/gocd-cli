package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ghodss/yaml"
)

var cliWriter = bufio.NewWriter(os.Stdout)

func render(value interface{}) error {
	if cliCfg.JSON {
		if err := toJSON(value); err != nil {
			return err
		}

		return nil
	}

	if cliCfg.YAML {
		if err := toYAML(value); err != nil {
			return err
		}

		return nil
	}

	cliLogger.Debug("no format was specified for rendering output to defaults")

	fmt.Printf("%v\n", value)

	return nil
}

func toYAML(value interface{}) error {
	cliLogger.Debug("rendering output in yaml format since --yaml is enabled")
	valueYAML, err := yaml.Marshal(value)
	if err != nil {
		return err
	}

	yamlString := strings.Join([]string{"---", string(valueYAML)}, "\n")

	_, err = cliWriter.Write([]byte(yamlString))
	if err != nil {
		cliLogger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			cliLogger.Fatalln(err)
		}
	}(cliWriter)

	return nil
}

func toJSON(value interface{}) error {
	cliLogger.Debug("rendering output in json format since --json is enabled")
	valueJSON, err := json.MarshalIndent(value, " ", " ")
	if err != nil {
		return err
	}

	jsonString := strings.Join([]string{string(valueJSON)}, "\n")

	_, err = cliWriter.Write([]byte(jsonString))
	if err != nil {
		cliLogger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			cliLogger.Fatalln(err)
		}
	}(cliWriter)

	return nil
}
