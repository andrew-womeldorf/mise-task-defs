package buildinfo

// Version gets set at build time with the linker tool.
//
// https://go.dev/wiki/GcToolchainTricks#including-build-information-in-the-executable
//
//	go build -ldflags "-X pkg.utils.buildinfo.Version=foo"
var (
	Commit  = "none"
	Date    = "unknown"
	Version = "local"
)
