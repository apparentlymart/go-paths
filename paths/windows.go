package paths

import (
	"net/url"
	"strings"
)

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
	u := &url.URL{
		Path: im.Clean(path),
	}
	// For an absolute path we also set the scheme.
	if im.IsAbs(path) {
		u.Scheme = "file"
	}
	return u
}

func (im impl) windowsFromURL(u *url.URL) (string, error) {
	return "", nil
}

func (im impl) isUNC(path string) bool {
	return im.volumeNameLen(path) > 2
}
