package render

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/sirupsen/logrus"
)

// Renderer implements methods to render output in JSON/YAML format.
type Renderer struct {
	writer *bufio.Writer
	logger *logrus.Logger
	YAML   bool
	JSON   bool
}

// Render renders the output based on the output format selection (toYAML, toJSON).
// If none is selected it prints as the source.
func (r *Renderer) Render(value interface{}) error {
	if r.JSON {
		if err := r.toJSON(value); err != nil {
			return err
		}

		return nil
	}

	if r.YAML {
		if err := r.toYAML(value); err != nil {
			return err
		}

		return nil
	}

	r.logger.Debug("no format was specified for rendering output to defaults")

	fmt.Printf("%v\n", value)

	return nil
}

func (r *Renderer) toYAML(value interface{}) error {
	r.logger.Debug("rendering output in yaml format since --yaml is enabled")
	valueYAML, err := yaml.Marshal(value)
	if err != nil {
		return err
	}

	yamlString := strings.Join([]string{"---", string(valueYAML)}, "\n")

	_, err = r.writer.Write([]byte(yamlString))
	if err != nil {
		r.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			r.logger.Fatalln(err)
		}
	}(r.writer)

	return nil
}

func (r *Renderer) toJSON(value interface{}) error {
	r.logger.Debug("rendering output in json format since --json is enabled")
	valueJSON, err := json.MarshalIndent(value, "", "     ")
	if err != nil {
		return err
	}

	jsonString := strings.Join([]string{string(valueJSON), "\n"}, "")

	_, err = r.writer.Write([]byte(jsonString))
	if err != nil {
		r.logger.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			r.logger.Fatalln(err)
		}
	}(r.writer)

	return nil
}

// GetRenderer returns the new instance of Renderer.
func GetRenderer(writer io.Writer, log *logrus.Logger, yaml, json bool) Renderer {
	renderer := Renderer{
		logger: log,
		YAML:   yaml,
		JSON:   json,
	}

	if writer == nil {
		renderer.writer = bufio.NewWriter(os.Stdout)
	} else {
		renderer.writer = bufio.NewWriter(writer)
	}

	return renderer
}
