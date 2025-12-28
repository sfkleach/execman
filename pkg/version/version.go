package version

// Version information
const (
	Version   = "0.1.0"
	BuildDate = "TBD"
	GitCommit = "TBD"
)

// GetVersion returns the current version
func GetVersion() string {
	return "v" + Version
}

// GetFullVersion returns version with build info
func GetFullVersion() string {
	return Version + " (build: " + BuildDate + ", commit: " + GitCommit + ")"
}
