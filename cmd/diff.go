package cmd

import (
	"fmt"
	"os"
)

func CheckDiffAndAllow(oldData, newData string) error {
	hasDiff, diff, err := diffCfg.Diff(oldData, newData)
	if err != nil {
		return err
	}

	if !hasDiff {
		cliLogger.Info("no changes to the input file, nothing to update, quitting")
		os.Exit(0)
	}

	fmt.Printf("%s\n", diff)

	contains, option := cliShellReadConfig.Reader()
	if !contains {
		cliLogger.Fatalln(inputValidationFailureMessage)
	}

	if option.Short == "n" {
		cliLogger.Warn(optingOutMessage)

		os.Exit(0)
	}

	return nil
}
