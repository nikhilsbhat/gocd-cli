package cmd

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/nikhilsbhat/gocd-sdk-go"
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
	Version     string `json:"version,omitempty"`
	Revision    string `json:"revision,omitempty"`
	Environment string `json:"environment,omitempty"`
	BuildDate   string `json:"buildDate,omitempty"`
	GoVersion   string `json:"goVersion,omitempty"`
	Platform    string `json:"platform,omitempty"`
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

func AppVersion(_ *cobra.Command, _ []string) error {
	buildInfo, err := json.Marshal(GetBuildInfo())
	if err != nil {
		log.Fatalf("fetching version of GoCD cli failed with: %v\n", err)
	}

	writer := bufio.NewWriter(os.Stdout)

	var serverVersionInfo string

	cliLogger.Debug("a call to GoCD server would be made to collect server version")

	if serverVersion, _ := client.GetVersionInfo(); !reflect.DeepEqual(serverVersion, gocd.ServerVersion{}) {
		cliLogger.Debug("got an update from GoCD server about server version")

		serverVersionJSON, err := json.Marshal(serverVersion)
		if err != nil {
			cliLogger.Errorf("fetching version of GoCD server failed with %v\n", err)
		}

		serverVersionInfo = strings.Join([]string{"server version", string(serverVersionJSON), "\n"}, ": ")

		//nolint:mirror
		if _, err = writer.Write([]byte(serverVersionInfo)); err != nil {
			log.Fatalln(err)
		}
	}

	cliVersionInfo := strings.Join([]string{"client version", string(buildInfo), "\n"}, ": ")

	//nolint:mirror
	if _, err = writer.Write([]byte(cliVersionInfo)); err != nil {
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
