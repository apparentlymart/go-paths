// If the following build tags change, remember to also update the documentation
// for the "Target" variable and the ForGOOS function, both of which also list
// which OSes we consider to be "Windows".

// +build windows

package paths

func init() {
	Target = Windows
}
