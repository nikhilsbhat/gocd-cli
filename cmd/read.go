package cmd

import (
	"io"
	"os"

	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/spf13/cobra"
)

func readObject(cmd *cobra.Command) (render.Object, error) {
	var obj render.Object

	if len(cliCfg.FromFile) != 0 {
		cliLogger.Debug("reading configuration object from file since --from-file is enabled")

		data, err := os.ReadFile(cliCfg.FromFile)
		if err != nil {
			return obj, err
		}

		obj = render.Object(data)
	} else {
		cliLogger.Debug("reading configuration object from stdin")

		stdIn := cmd.InOrStdin()

		in, err := io.ReadAll(stdIn)
		if err != nil {
			return obj, err
		}

		obj = render.Object(in)
	}

	return obj, nil
}
