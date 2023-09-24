package main

import (
	"log"

	"github.com/nikhilsbhat/gocd-cli/cmd"
	"github.com/spf13/cobra/doc"
)

//go:generate go run github.com/nikhilsbhat/gocd-cli/docs
func main() {
	commands := cmd.SetGoCDCliCommands()

	if err := doc.GenMarkdownTree(commands, "doc"); err != nil {
		log.Fatal(err)
	}
}
