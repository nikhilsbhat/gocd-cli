package main

import (
	"log"

	"github.com/nikhilsbhat/gocd-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	commands := cmd.SetGoCDCliCommands()
	err := doc.GenMarkdownTree(commands, "docs")
	if err != nil {
		log.Fatal(err)
	}
}
