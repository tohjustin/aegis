package version

var (
	// Version represents the application semantic version, variable will be replaced at link time after `make` has been run.
	Version = "latest"
	// GitHash represents the application Git SHA-1 hash, variable will be replaced at link time after `make` has been run.
	GitHash = "<UNKNOWN>"
)
