package utils

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Object implements method that check for file content type.
type Object string

const (
	FileTypeYAML    = "yaml"
	FileTypeJSON    = "json"
	FileTypeUnknown = "unknown"
)

// IsJSON checks if the passed content of JSON.
func IsJSON(content string) bool {
	var js interface{}

	return json.Unmarshal([]byte(content), &js) == nil
}

// IsYAML checks if the passed content of YAML.
func IsYAML(content string) bool {
	var js interface{}

	return yaml.Unmarshal([]byte(content), &js) == nil
}

// CheckFileType checks the file type of the content passed, it validates for YAML/JSON.
func (obj Object) CheckFileType(log *logrus.Logger) string {
	log.Debug("identifying the input file type, only YAML/JSON is allowed")
	if IsJSON(string(obj)) {
		log.Debug("input file type identified as JSON")

		return FileTypeJSON
	}

	if IsYAML(string(obj)) {
		log.Debug("input file type identified as YAML")

		return FileTypeYAML
	}

	log.Debug("input file type identified as UNKNOWN")

	return FileTypeUnknown
}
