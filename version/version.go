// Package version powers the versioning of gocd-cli.
package version

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	// Version specifies the version of the application and cannot be changed by end user.
	Version string

	// Env tells end user that what variant (here we use the name of the git branch to make it simple)
	// of application is he using.
	Env string

	// BuildDate of the app.
	BuildDate string
	// GoVersion represents golang version used.
	GoVersion string
	// Platform is the combination of OS and Architecture for which the binary is built for.
	Platform string
	// Revision represents the git revision used to build the current version of app.
	Revision string
)

// BuildInfo represents version of utility.
type BuildInfo struct {
	Version     string
	Revision    string
	Environment string
	BuildDate   string
	GoVersion   string
	Platform    string
}

// GetBuildInfo return the version and other build info of the application.
func GetBuildInfo() BuildInfo {
	if strings.ToLower(Env) != "production" {
		Env = "alfa"
	}

	return BuildInfo{
		Version:     Version,
		Revision:    Revision,
		Environment: Env,
		Platform:    Platform,
		BuildDate:   BuildDate,
		GoVersion:   GoVersion,
	}
}

func AppVersion(cmd *cobra.Command, args []string) error {
	buildInfo, err := json.Marshal(GetBuildInfo())
	if err != nil {
		log.Fatalf("fetching version of helm-images failed with: %v", err)
	}

	writer := bufio.NewWriter(os.Stdout)
	versionInfo := strings.Join([]string{"images version", string(buildInfo)}, ": ")
	_, err = writer.Write([]byte(versionInfo))
	if err != nil {
		log.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			log.Fatalln(err)
		}
	}(writer)

	return nil
}
