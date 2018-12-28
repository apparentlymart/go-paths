# Go Path Manipulation Utilities

Package `github.com/apparentlymart/go-paths/paths` is a subset of `path/filepath`
that provides access to multiple OS-specific implementations at once, regardless
of the real target of the calling program.

This might be useful in programs that must prepare commands to run on remote
systems with different operating systems, for example.

For more information, see [the package godoc](https://godoc.org/github.com/apparentlymart/go-paths/paths).