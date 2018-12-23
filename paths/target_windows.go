// If the following build tags change, remember to also update the documentation
// for the "Target" variable, which lists which OSes we consider to be "Windows"

// +build windows

package paths

func init() {
	Target = Windows
}
