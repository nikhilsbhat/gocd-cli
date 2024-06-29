## gocd-cli roles delete

Command to DELETE a specific role present in GoCD [https://api.gocd.org/current/#delete-a-role]

```
gocd-cli roles delete [flags]
```

### Examples

```
gocd-cli role delete sample-config
gocd-cli role delete sample-config -y
```

### Options

```
  -h, --help   help for delete
```

### Options inherited from parent commands

```
      --api-log-level string     log level for GoCD API calls, this sets log level to [https://pkg.go.dev/github.com/go-resty/resty/v2#Client.SetLogger],log levels supported by [https://github.com/sirupsen/logrus] will work (default "info")
      --api-retry-count int      number to times to retry when api calls fails, the value passed here would be set to GoCD sdk client
      --api-retry-interval int   time interval to wait before making subsequent API calls following API call failures (in seconds) (default 5)
  -t, --auth-token string        token to authenticate with GoCD server, should not be co-used with basic auth (username/password)
      --ca-file-path string      path to file containing CA cert used to authenticate GoCD server, if you have one
      --from-file string         file containing configurations of objects that needs to be created in GoCD, config-repo/pipeline-group/environment and etc.
  -l, --log-level string         log level for GoCD cli, log levels supported by [https://github.com/sirupsen/logrus] will work (default "info")
      --no-color                 enable this to Render output with no color
  -o, --output string            the format to which the output should be rendered to, it should be one of yaml|json|table|csv, if nothing specified it sets to default
  -p, --password string          password to authenticate with GoCD server
  -q, --query string             query to filter the results, ex: '.material.attributes.type | id eq git'. this uses library gojsonq beneath
                                 more queries can be found here https://github.com/thedevsaddam/gojsonq/wiki/Queries
      --server-url string        GoCD server URL base path (default "http://localhost:8153/go")
      --skip-cache-config        if enabled locally save auth configs would not be used to authenticate GoCD server (path: $HOME/.gocd/auth_config.yaml)
      --to-file string           file to which the output needs to be written
  -u, --username string          username to authenticate with GoCD server
  -y, --yes                      when enabled, end user confirmation would be skipped
```

### SEE ALSO

* [gocd-cli roles](gocd-cli_roles.md)	 - Command to operate on roles present in GoCD [https://api.gocd.org/current/#roles]

###### Auto generated by spf13/cobra on 29-Jun-2024