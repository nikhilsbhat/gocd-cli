## gocd-cli authorization

Command to operate on authorization-configuration present in GoCD [https://api.gocd.org/current/#authorization-configuration]

### Synopsis

Command leverages GoCD authorization-configuration apis' [https://api.gocd.org/current/#authorization-configuration] to 
GET/CREATE/UPDATE/DELETE cluster profiles present in GoCD

```
gocd-cli authorization [flags]
```

### Options

```
  -h, --help   help for authorization
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
* [gocd-cli authorization create](gocd-cli_authorization_create.md)	 - Command to CREATE the authorization configuration with specified configuration [https://api.gocd.org/current/#create-an-authorization-configuration]
* [gocd-cli authorization delete](gocd-cli_authorization_delete.md)	 - Command to DELETE the specified authorization configuration present in GoCD [https://api.gocd.org/current/#delete-an-authorization-configuration]
* [gocd-cli authorization get](gocd-cli_authorization_get.md)	 - Command to GET a authorization configuration with all specified configurations in GoCD [https://api.gocd.org/current/#get-an-authorization-configuration]
* [gocd-cli authorization get-all](gocd-cli_authorization_get-all.md)	 - Command to GET all authorization configurations present in GoCD [https://api.gocd.org/current/#get-all-authorization-configurations]
* [gocd-cli authorization list](gocd-cli_authorization_list.md)	 - Command to LIST all authorization configurations present in GoCD [https://api.gocd.org/current/#get-all-authorization-configurations]
* [gocd-cli authorization update](gocd-cli_authorization_update.md)	 - Command to UPDATE the authorization configuration present in GoCD [https://api.gocd.org/current/#update-an-authorization-configuration]

###### Auto generated by spf13/cobra on 27-Jan-2025
