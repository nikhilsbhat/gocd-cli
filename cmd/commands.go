package cmd

import (
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
)

type cliCommands struct {
	commands []*cobra.Command
}

// Config holds the information of the cli config.
type Config struct {
	URL             string    `yaml:"url,omitempty"`
	CaPath          string    `yaml:"ca_path,omitempty"`
	Auth            gocd.Auth `yaml:"auth,omitempty"`
	JSON            bool      `yaml:"-"`
	YAML            bool      `yaml:"-"`
	NoColor         bool      `yaml:"-"`
	LogLevel        string    `yaml:"-"`
	FromFile        string    `yaml:"-"`
	ToFile          string    `yaml:"-"`
	saveConfig      bool
	skipCacheConfig bool
}

func setGoCDCliCommands() *cobra.Command {
	return getGoCDCliCommands()
}

// Add an entry in below function to register new command.
func getGoCDCliCommands() *cobra.Command {
	command := new(cliCommands)
	command.commands = append(command.commands, getEncryptionCommand())
	command.commands = append(command.commands, getVersionCommand())
	command.commands = append(command.commands, getConfigRepoCommand())
	command.commands = append(command.commands, getBackupCommand())
	command.commands = append(command.commands, getUsersCommand())
	command.commands = append(command.commands, getEnvironmentsCommand())

	return command.prepareCommands()
}

func (c *cliCommands) prepareCommands() *cobra.Command {
	rootCmd := getRootCommand()
	for _, cmnd := range c.commands {
		rootCmd.AddCommand(cmnd)
	}
	registerGlobalFlags(rootCmd)

	return rootCmd
}

func getRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:     "gocd-cli",
		Short:   "Command line interface for GoCD",
		Long:    `Command line interface for GoCD that helps in interacting with GoCD CI/CD server`,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Usage(); err != nil {
				return err
			}

			return nil
		},
	}
	rootCommand.SetUsageTemplate(getUsageTemplate())

	return rootCommand
}

func getVersionCommand() *cobra.Command {
	versionCommand := &cobra.Command{
		Use:     "version [flags]",
		Short:   "Command to fetch the version of gocd-cli installed",
		Long:    `This will help user to find what version of gocd-cli he/she installed in her machine.`,
		PreRunE: setCLIClient,
		RunE:    AppVersion,
	}
	versionCommand.SetUsageTemplate(getUsageTemplate())

	return versionCommand
}

func getUsageTemplate() string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if gt (len .Aliases) 0}}{{printf "\n" }}
Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}{{printf "\n" }}
Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{printf "\n"}}
Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}{{printf "\n"}}
Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}{{printf "\n"}}
Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}{{printf "\n"}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}
{{if .HasAvailableSubCommands}}{{printf "\n"}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
{{printf "\n"}}`
}
