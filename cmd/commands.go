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
	APILogLevel     string    `yaml:"-"`
	FromFile        string    `yaml:"-"`
	ToFile          string    `yaml:"-"`
	skipCacheConfig bool
}

func SetGoCDCliCommands() *cobra.Command {
	return getGoCDCliCommands()
}

// Add an entry in below function to register new command.
func getGoCDCliCommands() *cobra.Command {
	command := new(cliCommands)
	command.commands = append(command.commands, registerEncryptionCommand())
	command.commands = append(command.commands, registerVersionCommand())
	command.commands = append(command.commands, registerConfigRepoCommand())
	command.commands = append(command.commands, registerBackupCommand())
	command.commands = append(command.commands, registerUsersCommand())
	command.commands = append(command.commands, registerEnvironmentsCommand())
	command.commands = append(command.commands, registerPluginsCommand())
	command.commands = append(command.commands, registerClusterProfilesCommand())
	command.commands = append(command.commands, registerAgentProfilesCommand())
	command.commands = append(command.commands, registerAgentCommand())
	command.commands = append(command.commands, registerServerCommand())
	command.commands = append(command.commands, registerPipelineGroupsCommand())
	command.commands = append(command.commands, registerPipelinesCommand())
	command.commands = append(command.commands, registerMaintenanceCommand())
	command.commands = append(command.commands, registerJobsCommand())
	command.commands = append(command.commands, registerStageCommand())
	command.commands = append(command.commands, registerArtifactCommand())
	command.commands = append(command.commands, registerAuthConfigCommand())
	command.commands = append(command.commands, registerMaterialsCommand())
	command.commands = append(command.commands, registerWhoAmICommand())
	command.commands = append(command.commands, registerServerConfigCommand())
	command.commands = append(command.commands, registerIHaveCommand())

	return command.prepareCommands()
}

func (c *cliCommands) prepareCommands() *cobra.Command {
	rootCmd := getRootCommand()
	for _, cmnd := range c.commands {
		rootCmd.AddCommand(cmnd)
	}

	rootCmd.SilenceErrors = true
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
			return cmd.Usage()
		},
	}
	rootCommand.SetUsageTemplate(getUsageTemplate())

	return rootCommand
}

func registerVersionCommand() *cobra.Command {
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
