// Package paths implements utility routines for manipulating filename paths
// in a way compatible with different operating systems that may or may not
// be the target of the calling program.
//
// This is, essentially, a portable subset of the go path/filepath package
// repackaged so that the routines for each OS are available regardless of
// the compilation target of the program, which can be useful if for example
// a program is preparing paths to be sent to a remote system using a
// different operating system.
//
// This package includes only functions that do not require access to a real
// filesystem. However, paths returned by methods of paths.Target can be used
// with the path/filepath package functions for that additional functionality.
//
// OS-Specific Paths
//
// The main entrypoints to this package are the variables Unix and Windows,
// which have methods with the same names and behaviors as those in the
// path/filepath package, and behave as that package would on the respective
// GOOS.
//
// The variable Target refers to either Unix or Windows depending on GOOS,
// and thus it effectively provides aliases for a subset of the functions
// of path/filepath.
//
// Slash Paths
//
// The variable Slash is a wrapper around the "path" package for handling
// slash-based paths as seen in URLs.
//
// Other Implementations
//
// Interface P is the type of all of the different path implementations in this
// package. Other packages can potentially implement this interface; although
// it has a large number of methods, the methodset will not change until a
// hypothetical major version 2 of this package.
package paths
