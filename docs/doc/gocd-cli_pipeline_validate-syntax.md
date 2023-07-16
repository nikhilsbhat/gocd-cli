## gocd-cli pipeline validate-syntax

Command validate pipeline syntax by running it against appropriate GoCD plugin

```
gocd-cli pipeline validate-syntax [flags]
```

### Examples

```
gocd-cli pipeline validate-syntax --pipeline pipeline1 --pipeline pipeline2
```

### Options

```
      --fetch-version-from-server    if enabled, plugin(auto-detected) version would be fetched from GoCD server
  -h, --help                         help for validate-syntax
      --pipeline strings             list of pipelines for which the syntax has to be validated
      --plugin-download-url string   Auto-detection of the plugin sets the download URL too (Github's release URL); if the URL needs to be set to something else, then it can be set using this
      --plugin-path string           if you prefer managing plugins outside the gocd-cli, the path to already downloaded plugins can be set using this
      --plugin-version string        GoCD plugin version against which the pipeline has to be validated (the plugin type would be auto-detected); if missed, the pipeline would be validated against the latest version of the auto-detected plugin
```

### Options inherited from parent commands

```
  -t, --auth-token string     token to authenticate with GoCD server, should not be co-used with basic auth (username/password)
      --ca-file-path string   path to file containing CA cert used to authenticate GoCD server, if you have one
      --from-file string      file containing configurations of objects that needs to be created in GoCD, config-repo/pipeline-group/environment and etc.
      --json                  enable this to Render output in JSON format
  -l, --log-level string      log level for gocd cli, log levels supported by [https://github.com/sirupsen/logrus] will work (default "info")
      --no-color              enable this to Render output in YAML format
  -p, --password string       password to authenticate with GoCD server
  -q, --query string          query to filter the results, ex: '.material.attributes.type | id eq git'. this uses library gojsonq beneath
                              more queries can be found here https://github.com/thedevsaddam/gojsonq/wiki/Queries
      --server-url string     GoCD server URL base path (default "http://localhost:8153/go")
      --skip-cache-config     if enabled locally save auth configs would not be used to authenticate GoCD server (path: $HOME/.gocd/auth_config.yaml)
      --to-file string        file to which the output needs to be written to
  -u, --username string       username to authenticate with GoCD server
      --yaml                  enable this to Render output in YAML format
```

### SEE ALSO

* [gocd-cli pipeline](gocd-cli_pipeline.md)	 - Command to operate on pipelines present in GoCD

###### Auto generated by spf13/cobra on 6-Jul-2023