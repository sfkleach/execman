package version

// Version information.
// These variables are set via ldflags during build.
var (
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

// GetVersion returns the current version.
func GetVersion() string {
	return "v" + Version
}

// GetFullVersion returns version with build info.
func GetFullVersion() string {
	return Version + " (build: " + BuildDate + ", commit: " + GitCommit + ")"
}
