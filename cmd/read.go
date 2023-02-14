package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

type Object string

func readObject(cmd *cobra.Command) (Object, error) {
	var obj Object
	if len(cliCfg.FromFile) != 0 {
		cliLogger.Debug("reading configuration object from file since --from-file is enabled")
		data, err := os.ReadFile(cliCfg.FromFile)
		if err != nil {
			return obj, err
		}
		obj = Object(data)
	} else {
		cliLogger.Debug("reading configuration object from stdin")
		stdIn := cmd.InOrStdin()
		in, err := io.ReadAll(stdIn)
		if err != nil {
			return obj, err
		}
		obj = Object(in)
	}

	return obj, nil
}
