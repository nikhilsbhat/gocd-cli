## gocd-cli configrepo

Command to operate on configrepo present in GoCD [https://api.gocd.org/current/#config-repo]

### Synopsis

Command leverages GoCD config repo apis' [https://api.gocd.org/current/#config-repo] to 
GET/CREATE/UPDATE/DELETE and trigger update on the same

```
gocd-cli configrepo [flags]
```

### Options

```
  -h, --help   help for configrepo
```

### Options inherited from parent commands

```
      --api-log-level string      log level for GoCD API calls, this sets log level to [https://pkg.go.dev/github.com/go-resty/resty/v2#Client.SetLogger],log levels supported by [https://github.com/sirupsen/logrus] will work (default "info")
      --api-retry-count int       number to times to retry when api calls fails, the value passed here would be set to GoCD sdk client
      --api-retry-interval int    time interval to wait before making subsequent API calls following API call failures (in seconds) (default 5)
  -t, --auth-token string         token to authenticate with GoCD server, should not be co-used with basic auth (username/password)
      --ca-file-path string       path to file containing CA cert used to authenticate GoCD server, if you have one
      --from-file string          file containing configurations of objects that needs to be created in GoCD, config-repo/pipeline-group/environment and etc.
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

* [gocd-cli](gocd-cli.md)	 - Command line interface for GoCD
* [gocd-cli configrepo create](gocd-cli_configrepo_create.md)	 - Command to CREATE the config-repo with specified configuration [https://api.gocd.org/current/#create-a-config-repo]
* [gocd-cli configrepo delete](gocd-cli_configrepo_delete.md)	 - Command to DELETE the specified config-repo [https://api.gocd.org/current/#delete-a-config-repo]
* [gocd-cli configrepo get](gocd-cli_configrepo_get.md)	 - Command to GET the config-repo information with a specified ID present in GoCD [https://api.gocd.org/current/#get-a-config-repo]
* [gocd-cli configrepo get-all](gocd-cli_configrepo_get-all.md)	 - Command to GET all config-repo information present in GoCD [https://api.gocd.org/current/#get-all-config-repos]
* [gocd-cli configrepo get-definitions](gocd-cli_configrepo_get-definitions.md)	 - Command to GET config-repo definitions present in GoCD [https://api.gocd.org/current/#definitions-defined-in-config-repo]
* [gocd-cli configrepo get-internal](gocd-cli_configrepo_get-internal.md)	 - Command to GET all config repo information present in GoCD using internal api [/api/internal/config_repos]
Do not use this command unless you know what you are doing with it
* [gocd-cli configrepo list](gocd-cli_configrepo_list.md)	 - Command to LIST all configuration repository present in GoCD [https://api.gocd.org/current/#get-all-config-repos]
* [gocd-cli configrepo preflight-check](gocd-cli_configrepo_preflight-check.md)	 - Command to PREFLIGHT check the config repo configurations [https://api.gocd.org/current/#preflight-check-of-config-repo-configurations]
* [gocd-cli configrepo status](gocd-cli_configrepo_status.md)	 - Command to GET the status of config-repo update operation [https://api.gocd.org/current/#status-of-config-repository-update]
* [gocd-cli configrepo trigger-update](gocd-cli_configrepo_trigger-update.md)	 - Command to TRIGGER the update for config-repo to get latest revisions [https://api.gocd.org/current/#trigger-update-of-config-repository]
* [gocd-cli configrepo update](gocd-cli_configrepo_update.md)	 - Command to UPDATE the config-repo present in GoCD [https://api.gocd.org/current/#update-config-repo]

###### Auto generated by spf13/cobra on 3-Nov-2024
