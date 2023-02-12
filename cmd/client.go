package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
)

var client gocd.GoCd

func setGoCDClient(cmd *cobra.Command, args []string) error {
	var caContent []byte

	if len(cliCfg.caPath) != 0 {
		caAbs, err := filepath.Abs(cliCfg.caPath)
		if err != nil {
			return err
		}

		caContent, err = os.ReadFile(caAbs)
		if err != nil {
			log.Fatal(err)
		}
	}

	goCDClient := gocd.NewClient(
		cliCfg.url,
		cliCfg.auth,
		cliCfg.loglevel,
		caContent,
	)

	client = goCDClient

	return nil
}
