package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	cliCfg    Config
	jsonQuery string
)

func registerGlobalFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cliCfg.URL, "server-url", "", "http://localhost:8153/go",
		"GoCD server URL base path")
	cmd.PersistentFlags().StringVarP(&cliCfg.Auth.UserName, "username", "u", "",
		"username to authenticate with GoCD server")
	cmd.PersistentFlags().StringVarP(&cliCfg.Auth.Password, "password", "p", "",
		"password to authenticate with GoCD server")
	cmd.PersistentFlags().StringVarP(&cliCfg.Auth.BearerToken, "auth-token", "t", "",
		"token to authenticate with GoCD server, should not be co-used with basic auth (username/password)")
	cmd.PersistentFlags().StringVarP(&cliCfg.CaPath, "ca-file-path", "", "",
		"path to file containing CA cert used to authenticate GoCD server, if you have one")
	cmd.PersistentFlags().StringVarP(&cliCfg.LogLevel, "log-level", "l", "info",
		"log level for gocd cli, log levels supported by [https://github.com/sirupsen/logrus] will work")
	cmd.PersistentFlags().BoolVarP(&cliCfg.JSON, "json", "", false,
		"enable this to Render output in JSON format")
	cmd.PersistentFlags().BoolVarP(&cliCfg.YAML, "yaml", "", false,
		"enable this to Render output in YAML format")
	cmd.PersistentFlags().BoolVarP(&cliCfg.YAML, "no-color", "", false,
		"enable this to Render output in YAML format")
	cmd.PersistentFlags().BoolVarP(&cliCfg.skipCacheConfig, "skip-cache-config", "", false,
		"if enabled locally save auth configs would not be used to authenticate GoCD server (path: $HOME/.gocd/auth_config.yaml)")
	cmd.PersistentFlags().StringVarP(&cliCfg.FromFile, "from-file", "", "",
		"file containing configurations of objects that needs to be created in GoCD, config-repo/pipeline-group/environment and etc.")
	cmd.PersistentFlags().StringVarP(&cliCfg.ToFile, "to-file", "", "",
		"file to which the output needs to be written to")
	cmd.PersistentFlags().StringVarP(&jsonQuery, "query", "q", "",
		`query to filter the results, ex: '.material.attributes.type | id eq git'. this uses library gojsonq beneath
more queries can be found here https://github.com/thedevsaddam/gojsonq/wiki/Queries`)
}

func registerEncryptionFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&cipherKey, "cipher-key", "", "",
		"cipher key value used for decryption, the key should same which is used by GoCD server for encryption")
	cmd.PersistentFlags().StringVarP(&cipherKeyPath, "cipher-key-path", "", "",
		"path to cipher key value used for decryption, the key should same which is used by GoCD server for encryption")
}

func registerConfigRepoDefinitionsFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&pipelines, "pipelines", "", false,
		"set this flag to get only the pipelines from the config-repo")
	cmd.PersistentFlags().BoolVarP(&environments, "environments", "", false,
		"set this flag to get only the environments from the config-repo")
}

func registerConfigRepoPreflightFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&configRepoPreflightObj.pluginID, "plugin-id", "i", "",
		"GoCD's config-repo plugin ID against which the pipelines has to be validated")
	cmd.PersistentFlags().StringSliceVarP(&configRepoPreflightObj.pipelineFiles, "pipeline-file", "f", nil,
		"GoCD pipeline files that should be considered for config-repo preflight checks")
	cmd.PersistentFlags().StringVarP(&configRepoPreflightObj.pipelineDir, "pipeline-dir", "", "",
		"path to directory that potentially contains the pipeline configuration file")
	cmd.PersistentFlags().StringVarP(&configRepoPreflightObj.pipelineExtRegex, "regex", "", "*.gocd.yaml",
		"regex to be used while identifying the pipeline files under the directory which was passed in pipeline-dir, "+
			"should be co-used with --pipeline-dir")
	cmd.PersistentFlags().BoolVarP(&configRepoPreflightObj.groovy, "groovy", "", false,
		"setting this would set 'plugin-id' to 'cd.go.contrib.plugins.configrepo.groovy'")
	cmd.PersistentFlags().BoolVarP(&configRepoPreflightObj.json, "json", "", false,
		"setting this would set 'plugin-id' to 'json.config.plugin'")
	cmd.PersistentFlags().BoolVarP(&configRepoPreflightObj.yaml, "yaml", "", false,
		"setting this would set 'plugin-id' to 'yaml.config.plugin'")
}

const (
	defaultRetryCount    = 30
	defaultDelay         = 5 * time.Second
	defaultInstanceCount = 0
)

func registerBackupFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&backupRetry, "retry", "", defaultRetryCount,
		"number of times to retry to get backup stats when backup status is not ready")
	cmd.PersistentFlags().DurationVarP(&delay, "delay", "", defaultDelay,
		"time delay between each retries that would be made to get backup stats")
}

func registerPipelineFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&goCDPipelineInstance, "instance", "i", defaultInstanceCount,
		"instance number of a pipeline")
	cmd.PersistentFlags().StringVarP(&goCDPipelineName, "name", "n", "",
		"name of the pipeline present in GoCD")
	cmd.PersistentFlags().StringVarP(&goCDPipelineETAG, "etag", "", "",
		"etag used to identify the pipeline config. If you don't have one get it by using pipeline get command")
	cmd.PersistentFlags().StringVarP(&goCDPipelineMessage, "message", "m", "",
		"message to be passed while pausing/unpausing or commenting on pipeline present in GoCD")
	cmd.PersistentFlags().BoolVarP(&goCDPipelinePause, "pause", "", false,
		"enable to pause a pipeline")
	cmd.PersistentFlags().BoolVarP(&goCDPipelineUnPause, "un-pause", "", false,
		"disable to pause a pipeline")
	cmd.PersistentFlags().BoolVarP(&goCDPausePipelineAtStart, "pause-at-start", "", false,
		"enabling this will create the pipeline in the paused state")
	cmd.PersistentFlags().StringVarP(&goCDPipelineTemplateName, "template-name", "", "",
		"name of the template to which the pipeline has to be extracted")
}

func registerMaintenanceFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&goCDEnableMaintenance, "enable", "", false,
		"set this to enable maintenance mode in GoCD")
	cmd.PersistentFlags().BoolVarP(&goCDDisableMaintenance, "disable", "", false,
		"set this to disable maintenance mode in GoCD")
}

func registerRawFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&rawOutput, "raw", "", false,
		"enable this to see the raw output")
}

func registerMaterialFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSliceVarP(&materialFilter, "filter", "", nil,
		"name of the material to filter from the available material in GoCD")
}

func registerPipelineHistoryFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().DurationVarP(&delay, "delay", "", defaultDelay,
		"delay between the calls made to GoCD server to get the pipeline run history in seconds")
	cmd.PersistentFlags().DurationVarP(&numberOfDays, "time", "", defaultDelay,
		"time frame since the pipeline has not run")
}

func registerAgentsFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&agentID, "id", "", "",
		"id of the agent on whom the action is to be performed")
	cmd.PersistentFlags().StringVarP(&agentName, "name", "", "",
		"name of the agent on whom the action is to be performed")
}

func registerAgentsFilterFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringSliceVarP(&agentResources, "resource", "", nil,
		"list of resource names to filter the agents from")
	cmd.PersistentFlags().StringSliceVarP(&agentEnvironments, "environment", "", nil,
		"list of environment names to filter the agents from")
	cmd.PersistentFlags().StringSliceVarP(&agentOS, "os", "", nil,
		"list of operating system names to filter the agents from")
	cmd.PersistentFlags().BoolVarP(&agentsDisabled, "disabled", "", false,
		"when enabled, it fetches only the disabled agents")
	cmd.PersistentFlags().StringVarP(&agentName, "name", "", "",
		"agent's name or pattern to match while filtering the results")
}

func registerJobsNStageFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&stageConfig.Pipeline, "pipeline", "", "",
		"pipeline name from which the jobs/stage to be triggered")
	cmd.PersistentFlags().StringVarP(&stageConfig.Name, "stage", "", "",
		"stage name that should to be operated")
	cmd.PersistentFlags().StringVarP(&stageConfig.PipelineInstance, "pipeline-counter", "", "",
		"instance of the pipeline that should be considered")
	cmd.PersistentFlags().StringVarP(&stageConfig.StageCounter, "stage-counter", "", "",
		"instance of the stage that should be considered")
	cmd.PersistentFlags().StringSliceVarP(&stageConfig.Jobs, "job", "", nil,
		"list of jobs that should be triggered")
}
