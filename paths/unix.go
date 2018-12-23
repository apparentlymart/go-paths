package paths

import (
	"errors"
	"net/url"
	"strings"
)

// Unix is a P implementation that consumes and generates paths suitable for
// Unix systems.
var Unix P

func init() {
	Unix = unixImpl
}

func (im impl) unixJoin(elem []string) string {
	for i, e := range elem {
		if e != "" {
			return im.Clean(strings.Join(elem[i:], string(im.separator())))
		}
	}
	return ""
}

func (im impl) unixToURL(path string) *url.URL {
	u := &url.URL{
		Path: im.Clean(path),
	}
	// For an absolute path we also set the scheme.
	if im.IsAbs(path) {
		u.Scheme = "file"
	}
	return u
}

func (im impl) unixFromURL(u *url.URL) (string, error) {
	switch u.Scheme {
	case "":
		// Relative URL
		return im.Clean(u.Path), nil
	case "file":
		if u.Host != "" && !strings.EqualFold(u.Host, "localhost") {
			return "", errors.New("only local file: URLs are allowed")
		}
		if u.User != nil {
			return "", errors.New("user portion not allowed in file: URLs")
		}
		// We'll tolerate and ignore query and fragment parts
		return im.Clean(u.Path), nil
	default:
		return "", errors.New("file: is the only allowed URL scheme")
	}
}
