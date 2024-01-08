package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nikhilsbhat/gocd-cli/pkg/errors"
	"github.com/nikhilsbhat/gocd-cli/pkg/render"
	"github.com/nikhilsbhat/gocd-sdk-go"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"gopkg.in/yaml.v3"
)

type configRepoPreflight struct {
	pipelineFiles    []string
	pipelineDir      string
	pipelineExtRegex string
}

var (
	pipelines, environments bool
	configRepoPreflightObj  configRepoPreflight
	queryEnabledMessage     = "since --query is passed, applying query '%v' to the output"
)

func registerConfigRepoCommand() *cobra.Command {
	configRepoCommand := &cobra.Command{
		Use:   "configrepo",
		Short: "Command to operate on configrepo present in GoCD [https://api.gocd.org/current/#config-repo]",
		Long: `Command leverages GoCD config repo apis' [https://api.gocd.org/current/#config-repo] to 
GET/CREATE/UPDATE/DELETE and trigger update on the same`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	configRepoCommand.SetUsageTemplate(getUsageTemplate())

	configRepoCommand.AddCommand(getConfigRepoTriggerUpdateCommand())
	configRepoCommand.AddCommand(getConfigRepoStatusCommand())
	configRepoCommand.AddCommand(getConfigReposCommand())
	configRepoCommand.AddCommand(getConfigRepoCommand())
	configRepoCommand.AddCommand(getCreateConfigRepoCommand())
	configRepoCommand.AddCommand(getUpdateConfigRepoCommand())
	configRepoCommand.AddCommand(getDeleteConfigRepoCommand())
	configRepoCommand.AddCommand(listConfigReposCommand())
	configRepoCommand.AddCommand(getConfigRepoPreflightCheckCommand())
	configRepoCommand.AddCommand(getConfigReposDefinitionsCommand())
	configRepoCommand.AddCommand(getFailedConfigReposCommand())

	for _, command := range configRepoCommand.Commands() {
		command.SilenceUsage = true
	}

	return configRepoCommand
}

func getConfigReposCommand() *cobra.Command {
	configGetCommand := &cobra.Command{
		Use:     "get-all",
		Short:   "Command to GET all config-repo information present in GoCD [https://api.gocd.org/current/#get-all-config-repos]",
		Example: "gocd-cli configrepo get-all",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetConfigRepos()
			if err != nil {
				return err
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := render.SetQuery(response, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debugf(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(response)
		},
	}

	configGetCommand.SetUsageTemplate(getUsageTemplate())

	return configGetCommand
}

func getFailedConfigReposCommand() *cobra.Command {
	var (
		detailed         bool
		failedConfigRepo bool
		getLastModified  bool
	)

	type configRepoLastModified struct {
		LastModified float64 `json:"lastModified,omitempty" yaml:"lastModified,omitempty"`
		Name         string  `json:"name,omitempty" yaml:"name,omitempty"`
		URL          string  `json:"url,omitempty" yaml:"url,omitempty"`
	}

	configGetCommand := &cobra.Command{
		Use: "get-internal",
		Short: `Command to GET all config repo information present in GoCD using internal api [/api/internal/config_repos]
Do not use this command unless you know what you are doing with it`,
		Example: "gocd-cli configrepo get-internal --failed --detailed --yaml",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetConfigReposInternal()
			if err != nil {
				return err
			}

			var repos interface{}

			if failedConfigRepo {
				response = funk.Filter(response, func(configRepo gocd.ConfigRepo) bool {
					return len(configRepo.ConfigRepoParseInfo.Error) != 0
				}).([]gocd.ConfigRepo)
			}

			if detailed {
				repos = response
			} else {
				names := make([]string, 0)
				for _, configRepo := range response {
					names = append(names, configRepo.ID)
				}
				repos = names
			}

			if getLastModified {
				configRepo := make([]configRepoLastModified, 0)
				for _, cfgRepo := range response {
					modificationTime := cfgRepo.ConfigRepoParseInfo.LatestParsedModification["modified_time"]
					if modificationTime != nil {
						modifiedDate := modificationTime.(string)
						configRepo = append(configRepo, configRepoLastModified{
							LastModified: lastUpdatedCommit(modifiedDate),
							Name:         cfgRepo.ID,
							URL:          cfgRepo.Material.Attributes.URL,
						})
					} else {
						cliLogger.Debugf("looks like config repo '%s' was never parsed, check the status of it", cfgRepo.ID)
					}
				}

				return cliRenderer.Render(configRepo)
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := render.SetQuery(repos, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debugf(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(repos)
		},
	}

	configGetCommand.PersistentFlags().BoolVarP(&detailed, "detailed", "", false,
		"when enabled prints the config-repo information in detailed")
	configGetCommand.PersistentFlags().BoolVarP(&failedConfigRepo, "failed", "", false,
		"when enabled, fetches only the failed config repositories")
	configGetCommand.PersistentFlags().BoolVarP(&getLastModified, "last-modified", "", false,
		"list config repo with last modified in number of days")
	configGetCommand.MarkFlagsMutuallyExclusive("last-modified", "detailed")

	configGetCommand.SetUsageTemplate(getUsageTemplate())

	return configGetCommand
}

func getConfigReposDefinitionsCommand() *cobra.Command {
	getConfigReposDefinitionsCmd := &cobra.Command{
		Use:   "get-definitions",
		Short: "Command to GET config-repo definitions present in GoCD [https://api.gocd.org/current/#definitions-defined-in-config-repo]",
		Example: `gocd-cli configrepo get-definitions sample-repo --yaml
gocd-cli configrepo get-definitions sample-repo --yaml --pipelines #should print only pipeline names`,
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetConfigRepoDefinitions(args[0])
			if err != nil {
				return err
			}

			var output interface{}
			output = response

			envsNPipelines := make([]string, 0)

			if environments {
				for _, env := range response.Environments {
					envsNPipelines = append(envsNPipelines, env.Name)
				}
				output = envsNPipelines
			}

			if pipelines {
				for _, group := range response.Groups {
					for _, pipeline := range group.Pipelines {
						envsNPipelines = append(envsNPipelines, pipeline.Name)
					}
				}
				output = envsNPipelines
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := render.SetQuery(output, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debugf(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(output)
		},
	}

	registerConfigRepoDefinitionsFlags(getConfigReposDefinitionsCmd)

	getConfigReposDefinitionsCmd.SetUsageTemplate(getUsageTemplate())

	return getConfigReposDefinitionsCmd
}

func getConfigRepoCommand() *cobra.Command {
	configGetCommand := &cobra.Command{
		Use:     "get",
		Short:   "Command to GET the config-repo information with a specified ID present in GoCD [https://api.gocd.org/current/#get-a-config-repo]",
		Example: "gocd-cli configrepo get helm-images",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetConfigRepo(args[0])
			if err != nil {
				return err
			}

			if len(jsonQuery) != 0 {
				cliLogger.Debugf(queryEnabledMessage, jsonQuery)

				baseQuery, err := render.SetQuery(response, jsonQuery)
				if err != nil {
					return err
				}

				cliLogger.Debugf(baseQuery.Print())

				return cliRenderer.Render(baseQuery.RunQuery())
			}

			return cliRenderer.Render(response)
		},
	}

	configGetCommand.SetUsageTemplate(getUsageTemplate())

	return configGetCommand
}

func getCreateConfigRepoCommand() *cobra.Command {
	configCreateStatusCommand := &cobra.Command{
		Use:     "create",
		Short:   "Command to CREATE the config-repo with specified configuration [https://api.gocd.org/current/#create-a-config-repo]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var configRepo gocd.ConfigRepo
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &configRepo); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &configRepo); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			if err = client.CreateConfigRepo(configRepo); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("config repo %s created successfully", configRepo.ID))
		},
	}

	configCreateStatusCommand.SetUsageTemplate(getUsageTemplate())

	return configCreateStatusCommand
}

func getUpdateConfigRepoCommand() *cobra.Command {
	configCreateStatusCommand := &cobra.Command{
		Use:     "update",
		Short:   "Command to UPDATE the config-repo present in GoCD [https://api.gocd.org/current/#update-config-repo]",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var configRepo gocd.ConfigRepo
			object, err := readObject(cmd)
			if err != nil {
				return err
			}

			switch objType := object.CheckFileType(cliLogger); objType {
			case render.FileTypeYAML:
				if err = yaml.Unmarshal([]byte(object), &configRepo); err != nil {
					return err
				}
			case render.FileTypeJSON:
				if err = json.Unmarshal([]byte(object), &configRepo); err != nil {
					return err
				}
			default:
				return &errors.UnknownObjectTypeError{Name: objType}
			}

			response, err := client.UpdateConfigRepo(configRepo)
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	configCreateStatusCommand.SetUsageTemplate(getUsageTemplate())

	return configCreateStatusCommand
}

func getDeleteConfigRepoCommand() *cobra.Command {
	configUpdateStatusCommand := &cobra.Command{
		Use:     "delete",
		Short:   "Command to DELETE the specified config-repo [https://api.gocd.org/current/#delete-a-config-repo]",
		Example: "gocd-cli configrepo delete helm-images",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			configRepoName := args[0]
			cliShellReadConfig.ShellMessage = fmt.Sprintf("do you want to delete config repo '%s'", configRepoName)

			if !cliCfg.Yes {
				contains, option := cliShellReadConfig.Reader()
				if !contains {
					cliLogger.Fatalln("user input validation failed, cannot proceed further")
				}

				if option.Short == "n" {
					cliLogger.Warn("not proceeding further since 'no' was opted")

					os.Exit(0)
				}
			}

			if err := client.DeleteConfigRepo(configRepoName); err != nil {
				return err
			}

			return cliRenderer.Render(fmt.Sprintf("config repo deleted: %s", configRepoName))
		},
	}

	return configUpdateStatusCommand
}

func listConfigReposCommand() *cobra.Command {
	listConfigReposCmd := &cobra.Command{
		Use:     "list",
		Short:   "Command to LIST all configuration repository present in GoCD [https://api.gocd.org/current/#get-all-config-repos]",
		Example: "gocd-cli configrepo list",
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.GetConfigRepos()
			if err != nil {
				return err
			}

			var configRepos []string

			for _, configRepo := range response {
				configRepos = append(configRepos, configRepo.ID)
			}

			return cliRenderer.Render(strings.Join(configRepos, "\n"))
		},
	}

	return listConfigReposCmd
}

func getConfigRepoStatusCommand() *cobra.Command {
	configStatusCommand := &cobra.Command{
		Use:     "status",
		Short:   "Command to GET the status of config-repo update operation [https://api.gocd.org/current/#status-of-config-repository-update]",
		Example: "gocd-cli configrepo status helm-images",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.ConfigRepoStatus(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	configStatusCommand.SetUsageTemplate(getUsageTemplate())

	return configStatusCommand
}

func getConfigRepoTriggerUpdateCommand() *cobra.Command {
	configTriggerUpdateCommand := &cobra.Command{
		Use:     "trigger-update",
		Short:   "Command to TRIGGER the update for config-repo to get latest revisions [https://api.gocd.org/current/#trigger-update-of-config-repository]",
		Example: "gocd-cli configrepo trigger-update helm-images",
		Args:    cobra.RangeArgs(1, 1),
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			response, err := client.ConfigRepoTriggerUpdate(args[0])
			if err != nil {
				return err
			}

			return cliRenderer.Render(response)
		},
	}

	configTriggerUpdateCommand.SetUsageTemplate(getUsageTemplate())

	return configTriggerUpdateCommand
}

func getConfigRepoPreflightCheckCommand() *cobra.Command {
	var goCDConfigRepoName string

	configTriggerUpdateCommand := &cobra.Command{
		Use:     "preflight-check",
		Short:   "Command to PREFLIGHT check the config repo configurations [https://api.gocd.org/current/#preflight-check-of-config-repo-configurations]",
		Example: `gocd-cli configrepo preflight-check -f path/to/pipeline1.gocd.yaml -f path/to/pipeline2.gocd.yaml --repo-name helm-images --yaml`,
		Args:    cobra.NoArgs,
		PreRunE: setCLIClient,
		RunE: func(cmd *cobra.Command, args []string) error {
			var pipelineFiles []gocd.PipelineFiles
			var pathAndPattern []string

			if len(configRepoPreflightObj.pipelineFiles) != 0 {
				for _, pipelineFile := range configRepoPreflightObj.pipelineFiles {
					file, err := client.GetPipelineFiles(pipelineFile)
					if err != nil {
						cliLogger.Errorf("fetching pipeline %s errored with: %v", pipelineFile, err)

						return err
					}
					pipelineFiles = append(pipelineFiles, file...)
				}
			} else {
				if len(configRepoPreflightObj.pipelineExtRegex) == 0 {
					return &errors.ConfigRepoError{Message: "pipeline file regex not passed, make sure to set --regex if --pipeline-dir is set"}
				}

				pathAndPattern = append(pathAndPattern, configRepoPreflightObj.pipelineDir)
				pathAndPattern = append(pathAndPattern, configRepoPreflightObj.pipelineExtRegex)
				file, err := client.GetPipelineFiles(pathAndPattern[0], pathAndPattern[1])
				if err != nil {
					cliLogger.Errorf("fetching pipeline using regex errored with: %v", err)

					return err
				}

				pipelineFiles = append(pipelineFiles, file...)
			}

			pipelineMap := client.SetPipelineFiles(pipelineFiles)

			response, err := client.ConfigRepoPreflightCheck(pipelineMap, goCdPluginObj.getPluginID(), goCDConfigRepoName)
			if err != nil {
				cliLogger.Errorf("preflight checks errored with: %v", err)

				return err
			}

			return cliRenderer.Render(response)
		},
	}

	configTriggerUpdateCommand.SetUsageTemplate(getUsageTemplate())
	registerConfigRepoPreflightFlags(configTriggerUpdateCommand)

	configTriggerUpdateCommand.PersistentFlags().StringVarP(&goCDConfigRepoName, "repo-name", "", "",
		"name of the config repo present in GoCD against which the pipeline has to be validated")

	return configTriggerUpdateCommand
}

func (cfg *goCdPlugin) getPluginID() string {
	if cfg.json {
		return "json.config.plugin"
	}

	if cfg.yaml {
		return "yaml.config.plugin"
	}

	if cfg.groovy {
		return "cd.go.contrib.plugins.configrepo.groovy"
	}

	return cfg.pluginID
}

func lastUpdatedCommit(date string) float64 {
	const hoursInADay = 24

	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalln(err)
	}

	tm, err := time.ParseInLocation(time.RFC3339, date, loc)
	if err != nil {
		log.Fatalln(err)
	}

	parsedTime := tm.In(loc)

	timeNow := time.Now().In(loc)

	diff := timeNow.Sub(parsedTime).Hours() / hoursInADay

	return diff
}
