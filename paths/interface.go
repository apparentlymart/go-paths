package paths

import "net/url"

// P is the main interface of this package, with each implementation providing
// path manipulation functionality for a particular target OS.
type P interface {
	Base(path string) string
	Clean(path string) string
	Dir(path string) string
	Ext(path string) string
	IsAbs(path string) bool
	Join(elems ...string) string
	Rel(basepath, targpath string) (string, error)
	Split(path string) (dir, file string)
	VolumeName(path string) string
	ToURL(path string) *url.URL
	FromURL(u *url.URL) (string, error)
}
