## gocd-cli pipeline

Command to operate on pipelines present in GoCD

### Synopsis

Command leverages GoCD pipeline apis'
[https://api.gocd.org/current/#pipeline-instances, https://api.gocd.org/current/#pipeline-config, https://api.gocd.org/current/#pipelines] to 
GET/PAUSE/UNPAUSE/UNLOCK/SCHEDULE and comment on a GoCD pipeline

```
gocd-cli pipeline [flags]
```

### Options

```
  -h, --help   help for pipeline
```

### Options inherited from parent commands

```
      --api-log-level string   log level for GoCD API calls, this sets log level to [https://pkg.go.dev/github.com/go-resty/resty/v2#Client.SetLogger],log levels supported by [https://github.com/sirupsen/logrus] will work (default "info")
  -t, --auth-token string      token to authenticate with GoCD server, should not be co-used with basic auth (username/password)
      --ca-file-path string    path to file containing CA cert used to authenticate GoCD server, if you have one
      --from-file string       file containing configurations of objects that needs to be created in GoCD, config-repo/pipeline-group/environment and etc.
      --json                   enable this to Render output in JSON format
  -l, --log-level string       log level for GoCD cli, log levels supported by [https://github.com/sirupsen/logrus] will work (default "info")
      --no-color               enable this to Render output in YAML format
  -p, --password string        password to authenticate with GoCD server
  -q, --query string           query to filter the results, ex: '.material.attributes.type | id eq git'. this uses library gojsonq beneath
                               more queries can be found here https://github.com/thedevsaddam/gojsonq/wiki/Queries
      --server-url string      GoCD server URL base path (default "http://localhost:8153/go")
      --skip-cache-config      if enabled locally save auth configs would not be used to authenticate GoCD server (path: $HOME/.gocd/auth_config.yaml)
      --to-file string         file to which the output needs to be written
  -u, --username string        username to authenticate with GoCD server
      --yaml                   enable this to Render output in YAML format
```

### SEE ALSO

* [gocd-cli](gocd-cli.md)	 - Command line interface for GoCD
* [gocd-cli pipeline action](gocd-cli_pipeline_action.md)	 - Command to PAUSE/UNPAUSE a specific pipeline present in GoCD,
              [https://api.gocd.org/current/#pause-a-pipeline,https://api.gocd.org/current/#unpause-a-pipeline]
* [gocd-cli pipeline comment](gocd-cli_pipeline_comment.md)	 - Command to COMMENT on a specific pipeline instance present in GoCD [https://api.gocd.org/current/#comment-on-pipeline-instance]
* [gocd-cli pipeline create](gocd-cli_pipeline_create.md)	 - Command to CREATE the pipeline with all specified configuration [https://api.gocd.org/current/#create-a-pipeline]
* [gocd-cli pipeline delete](gocd-cli_pipeline_delete.md)	 - Command to DELETE the specified pipeline from GoCD [https://api.gocd.org/current/#delete-a-pipeline]
* [gocd-cli pipeline export-format](gocd-cli_pipeline_export-format.md)	 - Command to export specified pipeline present in GoCD to appropriate config repo format [https://api.gocd.org/current/#export-pipeline-config-to-config-repo-format]
* [gocd-cli pipeline find](gocd-cli_pipeline_find.md)	 - Command to find all GoCD pipeline files present in a directory (it recursively finds for pipeline files in all sub-directory)
* [gocd-cli pipeline get](gocd-cli_pipeline_get.md)	 - Command to GET pipeline config of a specified pipeline present in GoCD [https://api.gocd.org/current/#get-pipeline-config]
* [gocd-cli pipeline get-all](gocd-cli_pipeline_get-all.md)	 - Command to GET all pipelines present in GoCD [https://api.gocd.org/current/#get-feed-of-all-stages-in-a-pipeline]
* [gocd-cli pipeline get-mappings](gocd-cli_pipeline_get-mappings.md)	 - Command to Identify the given pipeline is part of which config-repo/environment of GoCD
* [gocd-cli pipeline history](gocd-cli_pipeline_history.md)	 - Command to GET pipeline run history present in GoCD [https://api.gocd.org/current/#get-pipeline-history]
* [gocd-cli pipeline instance](gocd-cli_pipeline_instance.md)	 - Command to GET instance of a specific pipeline present in GoCD [https://api.gocd.org/current/#get-pipeline-instance]
* [gocd-cli pipeline last-schedule](gocd-cli_pipeline_last-schedule.md)	 - Command to GET last scheduled time of the pipeline present in GoCD [/pipelineHistory.json?pipelineName=nameOfThePipeline]
* [gocd-cli pipeline list](gocd-cli_pipeline_list.md)	 - Command to LIST all the pipelines present in GoCD [https://api.gocd.org/current/#get-feed-of-all-stages-in-a-pipeline]
* [gocd-cli pipeline not-scheduled](gocd-cli_pipeline_not-scheduled.md)	 - Command to GET pipelines not scheduled in last X days from GoCD [/pipelineHistory.json?]
* [gocd-cli pipeline schedule](gocd-cli_pipeline_schedule.md)	 - Command to SCHEDULE a specific pipeline present in GoCD [https://api.gocd.org/current/#scheduling-pipelines]
* [gocd-cli pipeline show](gocd-cli_pipeline_show.md)	 - Command to analyse pipelines part of a selected pipeline file
* [gocd-cli pipeline status](gocd-cli_pipeline_status.md)	 - Command to GET status of a specific pipeline present in GoCD [https://api.gocd.org/current/#get-pipeline-status]
* [gocd-cli pipeline template](gocd-cli_pipeline_template.md)	 - Command to EXTRACT template from specific pipeline instance present in GoCD [https://api.gocd.org/current/#extract-template-from-pipeline]
* [gocd-cli pipeline update](gocd-cli_pipeline_update.md)	 - Command to UPDATE the pipeline config with the latest specified configuration [https://api.gocd.org/current/#edit-pipeline-config]
* [gocd-cli pipeline validate-syntax](gocd-cli_pipeline_validate-syntax.md)	 - Command validate pipeline syntax by running it against appropriate GoCD plugin
* [gocd-cli pipeline vsm](gocd-cli_pipeline_vsm.md)	 - Command to GET downstream pipelines of a specified pipeline present in GoCD [https://api.gocd.org/current/#get-pipeline-config]

###### Auto generated by spf13/cobra on 22-Dec-2023
