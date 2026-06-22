package version

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// These variables are intended to be set via ldflags:
// go build -ldflags "-X cronx/pkg/version.GitVersion=v1.0.0 -X cronx/pkg/version.GitCommit=abc123"
var (
	GitVersion   = "v0.0.0-dev"
	GitCommit    = "unknown"
	GitTreeState = "unknown" // dirty / clean
	BuildDate    = "unknown" // ISO8601, $(date -u +'%Y-%m-%dT%H:%M:%SZ')
)

// Info holds versioning information.
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// Get returns an Info struct with all version information.
func Get() Info {
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a readable version string.
func String() string {
	v := Get()
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
