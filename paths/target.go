package paths

// Target is a path implementation that matches the current compilation target.
//
// For example, if GOOS=windows then this is equivalent to "Windows". On most
// other platforms it is equivalent to "Unix".
//
// To avoid potential incorrect behavior when new supported operating systems
// are added to Go, Target is nil for any platforms this package does not
// yet know about. For robust support for local paths on all target Go
// platforms, use the path/filepath package directly.
//
// In this version, the following operating systems are supported:
//
//     Unix:    aix darwin dragonfly freebsd linux netbsd openbsd solaris
//     Windows: windows
var Target P

// TargetRecognized returns true only if the target OS (GOOS) is recognized
// by this package. In other words, it returns true if the variable Target
// is safe to use.
func TargetRecognized() bool {
	return Target != nil
}

// ForGOOS returns the P implementation that this package would select for
// variable "Target" on the given GOOS-style operating system name. Returns
// nil if there is no implementation for the given OS.
func ForGOOS(goos string) P {
	switch goos {
	case "aix", "darwin", "dragonfly", "freebsd", "linux", "netbsd", "openbsd", "solaris":
		return Unix
	case "windows":
		return Windows
	default:
		return nil
	}
}
