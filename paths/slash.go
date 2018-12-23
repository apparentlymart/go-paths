package paths

import (
	"net/url"
	slashpath "path"
)

// Slash is a wrapper around the functionality of the generic slash-separated
// path handling in the standard library "path" package, exposed for convenience
// in programs that operate generically over path implementations or as a
// sort of "default" path implementation when no specific one is suitable.
var Slash P

func init() {
	Slash = slashImpl{}
}

type slashImpl struct{}

func (im slashImpl) Base(path string) string {
	return slashpath.Base(path)
}

func (im slashImpl) Clean(path string) string {
	return slashpath.Clean(path)
}

func (im slashImpl) Dir(path string) string {
	return slashpath.Dir(path)
}

func (im slashImpl) Ext(path string) string {
	return slashpath.Ext(path)
}

func (im slashImpl) IsAbs(path string) bool {
	return slashpath.IsAbs(path)
}

func (im slashImpl) Join(elems ...string) string {
	return slashpath.Join(elems...)
}

func (im slashImpl) Rel(basepath, targpath string) (string, error) {
	return Unix.Rel(basepath, targpath)
}

func (im slashImpl) Split(path string) (string, string) {
	return slashpath.Split(path)
}

func (im slashImpl) VolumeName(path string) string {
	// Slash paths never have volume names
	return ""
}

func (im slashImpl) ToURL(path string) *url.URL {
	return Unix.ToURL(path)
}

func (im slashImpl) FromURL(u *url.URL) (string, error) {
	return Unix.FromURL(u)
}
