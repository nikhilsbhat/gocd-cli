package cmd

import (
	"io"
	"os"

	"github.com/nikhilsbhat/common/content"
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
