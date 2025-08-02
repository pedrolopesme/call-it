package version

import "runtime"

var (
	// Version is the current version of the application.
	// This will be set during build time using ldflags.
	Version = "dev"
	
	// BuildTime is when the binary was built.
	// This will be set during build time using ldflags.
	BuildTime = "unknown"
	
	// GitCommit is the git commit hash.
	// This will be set during build time using ldflags.
	GitCommit = "unknown"
)

// Info returns version information about the application.
type Info struct {
	Version   string `json:"version"`
	BuildTime string `json:"build_time"`
	GitCommit string `json:"git_commit"`
	GoVersion string `json:"go_version"`
	Platform  string `json:"platform"`
}

// Get returns the version information.
func Get() Info {
	return Info{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}
}

// String returns a formatted version string.
func String() string {
	info := Get()
	return info.Version
}