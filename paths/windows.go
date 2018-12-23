package paths

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Windows is a P implementation that consumes and generates paths suitable for
// Windows systems.
var Windows P

func init() {
	Windows = windowsImpl
}

// windowsReservedNames lists reserved Windows names. Search for PRN in
// https://docs.microsoft.com/en-us/windows/desktop/fileio/naming-a-file
// for details.
var windowsReservedNames = []string{
	"CON", "PRN", "AUX", "NUL",
	"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
	"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
}

func isWindowsReservedName(path string) bool {
	if len(path) == 0 {
		return false
	}
	for _, reserved := range windowsReservedNames {
		if strings.EqualFold(path, reserved) {
			return true
		}
	}
	return false
}

func (im impl) windowsJoin(elem []string) string {
	for i, e := range elem {
		if e != "" {
			return im.windowsJoinNonEmpty(elem[i:])
		}
	}
	return ""
}

func (im impl) windowsJoinNonEmpty(elem []string) string {
	if len(elem[0]) == 2 && elem[0][1] == ':' {
		// First element is drive letter without terminating slash.
		// Keep path relative to current directory on that drive.
		// Skip empty elements.
		i := 1
		for ; i < len(elem); i++ {
			if elem[i] != "" {
				break
			}
		}
		return im.Clean(elem[0] + strings.Join(elem[i:], string(im.separator())))
	}
	// The following logic prevents Join from inadvertently creating a
	// UNC path on Windows. Unless the first element is a UNC path, Join
	// shouldn't create a UNC path. See golang.org/issue/9167.
	p := im.Clean(strings.Join(elem, string(im.separator())))
	if !im.isUNC(p) {
		return p
	}
	// p == UNC only allowed when the first element is a UNC path.
	head := im.Clean(elem[0])
	if im.isUNC(head) {
		return p
	}
	// head + tail == UNC, but joining two non-UNC paths should not result
	// in a UNC path. Undo creation of UNC path.
	tail := im.Clean(strings.Join(elem[1:], string(im.separator())))
	if head[len(head)-1] == im.separator() {
		return head + tail
	}
	return head + string(im.separator()) + tail
}

func (im impl) windowsToURL(path string) *url.URL {
	// We always produce the "standard" convention for file: URLs on Windows
	// as understood by the Windows shell and Internet Explorer:
	// https://blogs.msdn.microsoft.com/ie/2006/12/06/file-uris-in-windows/

	if !im.IsAbs(path) {
		// A relative path becomes a schemeless URL.
		return &url.URL{
			Path: im.toSlash(path),
		}
	}

	if im.isUNC(path) {
		// For UNC paths, the hostname part goes in the hostname part of
		// the file URI:
		//  \\foo\bar\baz
		//  file://foo/bar/baz
		volLen := im.volumeNameLen(path)
		return &url.URL{
			Scheme: "file",
			Host:   path[2:volLen],
			Path:   path[volLen:],
		}
	}

	// Absolute non-UNC (i.e. drive-lettered) paths become file: URLs with
	// an empty hostname part.
	return &url.URL{
		Scheme: "file",
		Path:   im.toSlash(path),
	}
}

func (im impl) windowsFromURL(u *url.URL) (string, error) {
	// There are lots of weird legacy forms of file: URI that are accepted
	// on Windows systems. This is not a faithful bug-compatible implementation,
	// but we do try to honor some of the more common misguided forms.

	switch u.Scheme {
	case "":
		// Relative URL
		return im.Clean(u.Path), nil
	case "file":
		// We'll tolerate and ignore query and fragment parts
		if u.User != nil {
			return "", errors.New("user portion not allowed in file: URLs")
		}
		if u.Host != "" {
			// If a hostname is present then it's a UNC path.
			return im.Clean(fmt.Sprintf("\\%s%s", u.Host, u.Path)), nil
		}
		// Otherwise the first part of the path ought to be a drive letter
		// followed by an unescaped colon.
		p := u.Path
		if len(p) >= 3 && p[2] == '|' {
			// Legacy Netscape form with pipe instead of colon
			p = fmt.Sprintf("%s:%s", p[:2], p[3:])
		}
		if len(p) < 4 || p[0] != '/' || p[2] != ':' {
			return "", errors.New("local file: URLs must begin with a drive letter and then a colon in the path portion")
		}
		return im.Clean(p[1:]), nil // trim off leading slash
	default:
		return "", errors.New("file: is the only allowed URL scheme")
	}
}

func (im impl) isUNC(path string) bool {
	return im.volumeNameLen(path) > 2
}
