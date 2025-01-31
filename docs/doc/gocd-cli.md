## gocd-cli

Command line interface for GoCD

### Synopsis

Command line interface for GoCD that helps in interacting with GoCD CI/CD server

```
gocd-cli [flags]
```

### Options

```
      --api-log-level string      log level for GoCD API calls, this sets log level to [https://pkg.go.dev/github.com/go-resty/resty/v2#Client.SetLogger],log levels supported by [https://github.com/sirupsen/logrus] will work (default "info")
      --api-retry-count int       number to times to retry when api calls fails, the value passed here would be set to GoCD sdk client
      --api-retry-interval int    time interval to wait before making subsequent API calls following API call failures (in seconds) (default 5)
  -t, --auth-token string         token to authenticate with GoCD server, should not be co-used with basic auth (username/password)
      --ca-file-path string       path to file containing CA cert used to authenticate GoCD server, if you have one
      --from-file string          file containing configurations of objects that needs to be created in GoCD, config-repo/pipeline-group/environment and etc.
  -h, --help                      help for gocd-cli
  -l, --log-level string          log level for GoCD cli, log levels supported by [https://github.com/sirupsen/logrus] will work (default "info")
      --no-auth                   enabling this will disable authentication when connecting to the GoCD server
      --no-color                  enable this to Render output with no color
  -o, --output string             the format to which the output should be rendered to, it should be one of yaml|json|table|csv, if nothing specified it sets to default
  -p, --password string           password to authenticate with GoCD server
      --profile string            set the profile when managing multiple GoCD, ex: default, central etc (default "default")
  -q, --query string              query to filter the results, ex: '.material.attributes.type | id eq git'. this uses library gojsonq beneath
                                  more queries can be found here https://github.com/thedevsaddam/gojsonq/wiki/Queries
      --server-url string         GoCD server URL base path (default "http://localhost:8153/go")
      --skip-cache-config         if enabled locally save auth configs would not be used to authenticate GoCD server (path: $HOME/.gocd/auth_config.yaml)
      --to-file string            file to which the output needs to be written
  -u, --username string           username to authenticate with GoCD server
  -w, --watch                     enable this to monitor resources continuously, applicable only if supported by the command
      --watch-interval duration   time interval between each watch cycle (default 5s)
  -y, --yes                       when enabled, end user confirmation would be skipped
```

### SEE ALSO

* [gocd-cli agents](gocd-cli_agents.md)	 - Command to operate on agents present in GoCD [https://api.gocd.org/current/#agents]
* [gocd-cli artifact](gocd-cli_artifact.md)	 - Command to operate on artifacts store/config present in GoCD
* [gocd-cli auth-config](gocd-cli_auth-config.md)	 - Command to store/remove the authorization configuration to be used by the cli
* [gocd-cli authorization](gocd-cli_authorization.md)	 - Command to operate on authorization-configuration present in GoCD [https://api.gocd.org/current/#authorization-configuration]
* [gocd-cli backup](gocd-cli_backup.md)	 - Command to operate on backup in GoCD [https://api.gocd.org/current/#backups]
* [gocd-cli cluster-profile](gocd-cli_cluster-profile.md)	 - Command to operate on cluster-profile present in GoCD [https://api.gocd.org/current/#cluster-profiles]
* [gocd-cli configrepo](gocd-cli_configrepo.md)	 - Command to operate on configrepo present in GoCD [https://api.gocd.org/current/#config-repo]
* [gocd-cli elastic-agent-profile](gocd-cli_elastic-agent-profile.md)	 - Command to operate on elastic-agent-profile in GoCD [https://api.gocd.org/current/#elastic-agent-profiles]
* [gocd-cli encryption](gocd-cli_encryption.md)	 - Command to encrypt/decrypt plain text value [https://api.gocd.org/current/#encryption]
* [gocd-cli environment](gocd-cli_environment.md)	 - Command to operate on environments present in GoCD [https://api.gocd.org/current/#environment-config]
* [gocd-cli i-have](gocd-cli_i-have.md)	 - Command to check the permissions that the current user has
* [gocd-cli job](gocd-cli_job.md)	 - Command to operate on jobs present in GoCD
* [gocd-cli maintenance](gocd-cli_maintenance.md)	 - Command to operate on maintenance modes in GoCD [https://api.gocd.org/current/#maintenance-mode]
* [gocd-cli materials](gocd-cli_materials.md)	 - Command to operate on materials present in GoCD [https://api.gocd.org/current/#get-all-materials]
* [gocd-cli pipeline](gocd-cli_pipeline.md)	 - Command to operate on pipelines present in GoCD
* [gocd-cli pipeline-group](gocd-cli_pipeline-group.md)	 - Command to operate on pipeline groups present in GoCD [https://api.gocd.org/current/#pipeline-group-config]
* [gocd-cli plugin](gocd-cli_plugin.md)	 - Command to operate on plugins present in GoCD
* [gocd-cli roles](gocd-cli_roles.md)	 - Command to operate on roles present in GoCD [https://api.gocd.org/current/#roles]
* [gocd-cli server](gocd-cli_server.md)	 - Command to operate on GoCD server health status
* [gocd-cli server-config](gocd-cli_server-config.md)	 - Command to operate on GoCD server's configurations
* [gocd-cli stage](gocd-cli_stage.md)	 - Command to operate on stages of a pipeline present in GoCD
* [gocd-cli user](gocd-cli_user.md)	 - Command to operate on users in GoCD [https://api.gocd.org/current/#users]
* [gocd-cli version](gocd-cli_version.md)	 - Command to fetch the version of gocd-cli installed
* [gocd-cli who-am-i](gocd-cli_who-am-i.md)	 - Command to check which user being used by GoCD Command line interface

###### Auto generated by spf13/cobra on 27-Jan-2025
