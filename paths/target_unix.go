// If the following build tags change, remember to also update the documentation
// for the "Target" variable, which lists which OSes we consider to be "Unix"

// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package paths

func init() {
	Target = Unix
}
