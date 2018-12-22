package paths

import (
	"errors"
	"net/url"
	"strings"
)

// impl is am implementation of P that can either support Unix or Windows
// semantics, with many behaviors common to both.
type impl uint8

const windowsImpl impl = '\\'
const unixImpl impl = '/'

func (im impl) Base(path string) string {
	if path == "" {
		return "."
	}
	// Strip trailing slashes.
	for len(path) > 0 && im.isPathSeparator(path[len(path)-1]) {
		path = path[0 : len(path)-1]
	}
	// Throw away volume name
	path = path[len(im.VolumeName(path)):]
	// Find the last element
	i := len(path) - 1
	for i >= 0 && !im.isPathSeparator(path[i]) {
		i--
	}
	if i >= 0 {
		path = path[i+1:]
	}
	// If empty now, it had only slashes.
	if path == "" {
		return string(im.separator())
	}
	return path
}

func (im impl) Clean(path string) string {
	originalPath := path
	volLen := im.volumeNameLen(path)
	path = path[volLen:]
	if path == "" {
		if volLen > 1 && originalPath[1] != ':' {
			// should be UNC
			return im.fromSlash(originalPath)
		}
		return originalPath + "."
	}

	n := len(path)
	if volLen > 2 && n == 1 && im.isPathSeparator(path[0]) {
		// UNC volume name with trailing slash.
		return im.fromSlash(originalPath[:volLen])
	}
	rooted := im.isPathSeparator(path[0])

	// Invariants:
	//	reading from path; r is index of next byte to process.
	//	writing to out; w is index of next byte to write.
	//	dotdot is index in out where .. must stop, either because
	//		it is the leading slash or it is a leading ../../.. prefix.
	out := lazybuf{path: path, volAndPath: originalPath, volLen: volLen}
	r, dotdot := 0, 0
	if rooted {
		out.append(im.separator())
		r, dotdot = 1, 1
	}

	for r < n {
		switch {
		case im.isPathSeparator(path[r]):
			// empty path element
			r++
		case path[r] == '.' && (r+1 == n || im.isPathSeparator(path[r+1])):
			// . element
			r++
		case path[r] == '.' && path[r+1] == '.' && (r+2 == n || im.isPathSeparator(path[r+2])):
			// .. element: remove to last separator
			r += 2
			switch {
			case out.w > dotdot:
				// can backtrack
				out.w--
				for out.w > dotdot && !im.isPathSeparator(out.index(out.w)) {
					out.w--
				}
			case !rooted:
				// cannot backtrack, but not rooted, so append .. element.
				if out.w > 0 {
					out.append(im.separator())
				}
				out.append('.')
				out.append('.')
				dotdot = out.w
			}
		default:
			// real path element.
			// add slash if needed
			if rooted && out.w != 1 || !rooted && out.w != 0 {
				out.append(im.separator())
			}
			// copy element
			for ; r < n && !im.isPathSeparator(path[r]); r++ {
				out.append(path[r])
			}
		}
	}

	// Turn empty string into "."
	if out.w == 0 {
		out.append('.')
	}

	return im.fromSlash(out.string())
}

func (im impl) Dir(path string) string {
	vol := im.VolumeName(path)
	i := len(path) - 1
	for i >= len(vol) && !im.isPathSeparator(path[i]) {
		i--
	}
	dir := im.Clean(path[len(vol) : i+1])
	if dir == "." && len(vol) > 2 {
		// must be Windows UNC
		return vol
	}
	return vol + dir
}

func (im impl) Ext(path string) string {
	for i := len(path) - 1; i >= 0 && !im.isPathSeparator(path[i]); i-- {
		if path[i] == '.' {
			return path[i:]
		}
	}
	return ""
}

func (im impl) IsAbs(path string) bool {
	switch im {
	case windowsImpl:
		if isWindowsReservedName(path) {
			return true
		}
		l := im.volumeNameLen(path)
		if l == 0 {
			return false
		}
		path = path[l:]
		if path == "" {
			return false
		}
		return isSlash(path[0])
	default:
		return strings.HasPrefix(path, "/")
	}
}

func (im impl) Join(elems ...string) string {
	switch im {
	case windowsImpl:
		return im.windowsJoin(elems)
	default:
		return im.unixJoin(elems)
	}
}

func (im impl) Rel(basepath, targpath string) (string, error) {
	baseVol := im.VolumeName(basepath)
	targVol := im.VolumeName(targpath)
	base := im.Clean(basepath)
	targ := im.Clean(targpath)
	if im.sameWord(targ, base) {
		return ".", nil
	}
	base = base[len(baseVol):]
	targ = targ[len(targVol):]
	if base == "." {
		base = ""
	}
	// Can't use IsAbs - `\a` and `a` are both relative in Windows.
	baseSlashed := len(base) > 0 && base[0] == im.separator()
	targSlashed := len(targ) > 0 && targ[0] == im.separator()
	if baseSlashed != targSlashed || !im.sameWord(baseVol, targVol) {
		return "", errors.New("can't make " + targpath + " relative to " + basepath)
	}
	// Position base[b0:bi] and targ[t0:ti] at the first differing elements.
	bl := len(base)
	tl := len(targ)
	var b0, bi, t0, ti int
	for {
		for bi < bl && base[bi] != im.separator() {
			bi++
		}
		for ti < tl && targ[ti] != im.separator() {
			ti++
		}
		if !im.sameWord(targ[t0:ti], base[b0:bi]) {
			break
		}
		if bi < bl {
			bi++
		}
		if ti < tl {
			ti++
		}
		b0 = bi
		t0 = ti
	}
	if base[b0:bi] == ".." {
		return "", errors.New("can't make " + targpath + " relative to " + basepath)
	}
	if b0 != bl {
		// Base elements left. Must go up before going down.
		seps := strings.Count(base[b0:bl], string(im.separator()))
		size := 2 + seps*3
		if tl != t0 {
			size += 1 + tl - t0
		}
		buf := make([]byte, size)
		n := copy(buf, "..")
		for i := 0; i < seps; i++ {
			buf[n] = im.separator()
			copy(buf[n+1:], "..")
			n += 3
		}
		if t0 != tl {
			buf[n] = im.separator()
			copy(buf[n+1:], targ[t0:])
		}
		return string(buf), nil
	}
	return targ[t0:], nil
}

func (im impl) Split(path string) (string, string) {
	vol := im.VolumeName(path)
	i := len(path) - 1
	for i >= len(vol) && !im.isPathSeparator(path[i]) {
		i--
	}
	return path[:i+1], path[i+1:]
}

func (im impl) VolumeName(path string) string {
	return path[:im.volumeNameLen(path)]
}

func (im impl) ToURL(path string) *url.URL {
	switch im {
	case windowsImpl:
		return im.windowsToURL(path)
	default:
		return im.unixToURL(path)
	}
}

func (im impl) FromURL(u *url.URL) (string, error) {
	switch im {
	case windowsImpl:
		return im.windowsFromURL(u)
	default:
		return im.unixFromURL(u)
	}
}

func (im impl) volumeNameLen(path string) int {
	switch im {
	case windowsImpl:
		if len(path) < 2 {
			return 0
		}
		// with drive letter
		c := path[0]
		if path[1] == ':' && ('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z') {
			return 2
		}
		// is it UNC? https://msdn.microsoft.com/en-us/library/windows/desktop/aa365247(v=vs.85).aspx
		if l := len(path); l >= 5 && isSlash(path[0]) && isSlash(path[1]) &&
			!isSlash(path[2]) && path[2] != '.' {
			// first, leading `\\` and next shouldn't be `\`. its server name.
			for n := 3; n < l-1; n++ {
				// second, next '\' shouldn't be repeated.
				if isSlash(path[n]) {
					n++
					// third, following something characters. its share name.
					if !isSlash(path[n]) {
						if path[n] == '.' {
							break
						}
						for ; n < l; n++ {
							if isSlash(path[n]) {
								break
							}
						}
						return n
					}
					break
				}
			}
		}
		return 0
	default:
		return 0
	}
}

func (im impl) separator() uint8 {
	return uint8(im)
}

func (im impl) isPathSeparator(s uint8) bool {
	switch {
	case s == uint8(im):
		return true
	case im == windowsImpl && s == '/':
		return true
	default:
		return false
	}
}

func (im impl) fromSlash(path string) string {
	if im == unixImpl {
		return path
	}
	return strings.Replace(path, "/", string(im), -1)
}

func (im impl) toSlash(path string) string {
	if im == unixImpl {
		return path
	}
	return strings.Replace(path, string(im), "/", -1)
}

func (im impl) sameWord(a, b string) bool {
	switch im {
	case windowsImpl:
		return strings.EqualFold(a, b)
	default:
		return a == b
	}
}

// isSlash returns true if the given character is either a slash or a backslash,
// regardless of platform. For a platform-specific answer, use im.isPathSeparator.
func isSlash(c uint8) bool {
	return c == '\\' || c == '/'
}
