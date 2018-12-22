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
package paths
