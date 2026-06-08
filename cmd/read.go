package cmd

import (
	"encoding/json"
	"io"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/nikhilsbhat/common/content"
	clierrors "github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/spf13/cobra"
)

func readObject(cmd *cobra.Command) (content.Object, error) {
	var obj content.Object

	if len(cliCfg.FromFile) != 0 {
		cliLogger.Debug("reading configuration object from file since --from-file is enabled")

		data, err := os.ReadFile(cliCfg.FromFile)
		if err != nil {
			return obj, err
		}

		obj = content.Object(data)
	} else {
		cliLogger.Debug("reading configuration object from stdin")

		stdIn := cmd.InOrStdin()

		in, err := io.ReadAll(stdIn)
		if err != nil {
			return obj, err
		}

		obj = content.Object(in)
	}

	return obj, nil
}

func unmarshalObject(object content.Object, target any) error {
	switch objType := object.CheckFileType(cliLogger); objType {
	case content.FileTypeYAML:
		if err := yaml.UnmarshalWithOptions([]byte(object), target, yaml.AllowDuplicateMapKey()); err != nil {
			return err
		}
	case content.FileTypeJSON:
		if err := json.Unmarshal([]byte(object), target); err != nil {
			return err
		}
	default:
		return &clierrors.UnknownObjectTypeError{Name: objType}
	}

	return nil
}

func readConfigObject[T any](cmd *cobra.Command) (T, content.Object, error) {
	var config T

	object, err := readObject(cmd)
	if err != nil {
		return config, object, err
	}

	if err = unmarshalObject(object, &config); err != nil {
		return config, object, err
	}

	return config, object, nil
}
