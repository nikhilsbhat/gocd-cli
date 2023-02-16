package cmd

import (
	"io"
	"os"

	"github.com/nikhilsbhat/gocd-cli/pkg/utils"
	"github.com/spf13/cobra"
)

func readObject(cmd *cobra.Command) (utils.Object, error) {
	var obj utils.Object
	if len(cliCfg.FromFile) != 0 {
		cliLogger.Debug("reading configuration object from file since --from-file is enabled")
		data, err := os.ReadFile(cliCfg.FromFile)
		if err != nil {
			return obj, err
		}
		obj = utils.Object(data)
	} else {
		cliLogger.Debug("reading configuration object from stdin")
		stdIn := cmd.InOrStdin()
		in, err := io.ReadAll(stdIn)
		if err != nil {
			return obj, err
		}
		obj = utils.Object(in)
	}

	return obj, nil
}
