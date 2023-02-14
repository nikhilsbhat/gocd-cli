package cmd

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

const (
	fileTypeYAML    = "yaml"
	fileTypeJSON    = "json"
	fileTypeUnknown = "unknown"
)

func IsJSON(content string) bool {
	var js interface{}

	return json.Unmarshal([]byte(content), &js) == nil
}

func IsYAML(content string) bool {
	var js interface{}

	return yaml.Unmarshal([]byte(content), &js) == nil
}

func (obj Object) CheckFileType() string {
	cliLogger.Debug("identifying the input file type, only YAML/JSON is allowed")
	if IsJSON(string(obj)) {
		cliLogger.Debug("input file type identified as JSON")

		return fileTypeJSON
	}

	if IsYAML(string(obj)) {
		cliLogger.Debug("input file type identified as YAML")

		return fileTypeYAML
	}

	cliLogger.Debug("input file type identified as UNKNOWN")

	return fileTypeUnknown
}
