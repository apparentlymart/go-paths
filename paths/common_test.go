package paths

import (
	"testing"
)

func TestClean(t *testing.T) {
	tests := []struct {
		path, want string
	}{
		// Already clean
		{"abc", "abc"},
		{"abc/def", "abc/def"},
		{"a/b/c", "a/b/c"},
		{".", "."},
		{"..", ".."},
		{"../..", "../.."},
		{"../../abc", "../../abc"},
		{"/abc", "/abc"},
		{"/", "/"},

		// Empty is current dir
		{"", "."},

		// Remove trailing slash
		{"abc/", "abc"},
		{"abc/def/", "abc/def"},
		{"a/b/c/", "a/b/c"},
		{"./", "."},
		{"../", ".."},
		{"../../", "../.."},
		{"/abc/", "/abc"},

		// Remove doubled slash
		{"abc//def//ghi", "abc/def/ghi"},
		{"//abc", "/abc"},
		{"///abc", "/abc"},
		{"//abc//", "/abc"},
		{"abc//", "abc"},

		// Remove . elements
		{"abc/./def", "abc/def"},
		{"/./abc/def", "/abc/def"},
		{"abc/.", "abc"},

		// Remove .. elements
		{"abc/def/ghi/../jkl", "abc/def/jkl"},
		{"abc/def/../ghi/../jkl", "abc/jkl"},
		{"abc/def/..", "abc"},
		{"abc/def/../..", "."},
		{"/abc/def/../..", "/"},
		{"abc/def/../../..", ".."},
		{"/abc/def/../../..", "/"},
		{"abc/def/../../../ghi/jkl/../../../mno", "../../mno"},
		{"/../abc", "/abc"},

		// Combinations
		{"abc/./../def", "def"},
		{"abc//./../def", "def"},
		{"abc/../../././../def", "../../def"},
	}
	windowsTests := []struct {
		path, want string
	}{
		{`c:`, `c:.`},
		{`c:\`, `c:\`},
		{`c:\abc`, `c:\abc`},
		{`c:abc\..\..\.\.\..\def`, `c:..\..\def`},
		{`c:\abc\def\..\..`, `c:\`},
		{`c:\..\abc`, `c:\abc`},
		{`c:..\abc`, `c:..\abc`},
		{`\`, `\`},
		{`/`, `\`},
		{`\\i\..\c$`, `\c$`},
		{`\\i\..\i\c$`, `\i\c$`},
		{`\\i\..\I\c$`, `\I\c$`},
		{`\\host\share\foo\..\bar`, `\\host\share\bar`},
		{`//host/share/foo/../baz`, `\\host\share\baz`},
		{`\\a\b\..\c`, `\\a\b\c`},
		{`\\a\b`, `\\a\b`},
		{`\\a\b\`, `\\a\b`},
		{`\\folder\share\foo`, `\\folder\share\foo`},
		{`\\folder\share\foo\`, `\\folder\share\foo`},
	}

	impls := map[string]P{
		"Unix":    Unix,
		"Windows": Windows,
	}

	for name, p := range impls {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				t.Run(test.path, func(t *testing.T) {
					got := p.Clean(test.path)
					want := test.want
					if name == "Windows" {
						want = Windows.(impl).fromSlash(want)
					}
					if got != want {
						t.Errorf("wrong result for %s.Clean(%q)\ngot:  %s\nwant: %s", name, test.path, got, want)
					}
				})
			}
			switch name {
			case "Windows":
				for _, test := range windowsTests {
					t.Run(test.path, func(t *testing.T) {
						got := p.Clean(test.path)
						want := test.want
						if got != want {
							t.Errorf("wrong result for %s.Clean(%q)\ngot:  %s\nwant: %s", name, test.path, got, want)
						}
					})
				}
			}
		})
	}
}

func TestSplit(t *testing.T) {
	type Test struct {
		path, dir, file string
	}

	implTests := map[string][]Test{
		"Unix": {
			{"a/b", "a/", "b"},
			{"a/b/", "a/b/", ""},
			{"a/", "a/", ""},
			{"a", "", "a"},
			{"/", "/", ""},
		},
		"Windows": {
			{`a\b`, `a\`, `b`},
			{`a\b\`, `a\b\`, ""},
			{`a\`, `a\`, ``},
			{"a", "", "a"},
			{`\`, `\`, ""},
			{`c:`, `c:`, ``},
			{`c:/`, `c:/`, ``},
			{`c:/foo`, `c:/`, `foo`},
			{`c:/foo/bar`, `c:/foo/`, `bar`},
			{`//host/share`, `//host/share`, ``},
			{`//host/share/`, `//host/share/`, ``},
			{`//host/share/foo`, `//host/share/`, `foo`},
			{`\\host\share`, `\\host\share`, ``},
			{`\\host\share\`, `\\host\share\`, ``},
			{`\\host\share\foo`, `\\host\share\`, `foo`},
		},
	}

	impls := map[string]P{
		"Unix":    Unix,
		"Windows": Windows,
	}

	for implName, tests := range implTests {
		t.Run(implName, func(t *testing.T) {
			impl := impls[implName]
			for _, test := range tests {
				t.Run(test.path, func(t *testing.T) {
					if d, f := impl.Split(test.path); d != test.dir || f != test.file {
						t.Errorf("Split(%q) = %q, %q, want %q, %q", test.path, d, f, test.dir, test.file)
					}
				})
			}
		})
	}
}

func TestJoin(t *testing.T) {
	type Test struct {
		elem []string
		path string
	}

	implTests := map[string][]Test{
		"Unix": {
			// zero parameters
			{[]string{}, ""},

			// one parameter
			{[]string{""}, ""},
			{[]string{"/"}, "/"},
			{[]string{"a"}, "a"},

			// two parameters
			{[]string{"a", "b"}, "a/b"},
			{[]string{"a", ""}, "a"},
			{[]string{"", "b"}, "b"},
			{[]string{"/", "a"}, "/a"},
			{[]string{"/", "a/b"}, "/a/b"},
			{[]string{"/", ""}, "/"},
			{[]string{"//", "a"}, "/a"},
			{[]string{"/a", "b"}, "/a/b"},
			{[]string{"a/", "b"}, "a/b"},
			{[]string{"a/", ""}, "a"},
			{[]string{"", ""}, ""},

			// three parameters
			{[]string{"/", "a", "b"}, "/a/b"},
		},
		"Windows": {
			{[]string{`directory`, `file`}, `directory\file`},
			{[]string{`C:\Windows\`, `System32`}, `C:\Windows\System32`},
			{[]string{`C:\Windows\`, ``}, `C:\Windows`},
			{[]string{`C:\`, `Windows`}, `C:\Windows`},
			{[]string{`C:`, `a`}, `C:a`},
			{[]string{`C:`, `a\b`}, `C:a\b`},
			{[]string{`C:`, `a`, `b`}, `C:a\b`},
			{[]string{`C:`, ``, `b`}, `C:b`},
			{[]string{`C:`, ``, ``, `b`}, `C:b`},
			{[]string{`C:`, ``}, `C:.`},
			{[]string{`C:`, ``, ``}, `C:.`},
			{[]string{`C:.`, `a`}, `C:a`},
			{[]string{`C:a`, `b`}, `C:a\b`},
			{[]string{`C:a`, `b`, `d`}, `C:a\b\d`},
			{[]string{`\\host\share`, `foo`}, `\\host\share\foo`},
			{[]string{`\\host\share\foo`}, `\\host\share\foo`},
			{[]string{`//host/share`, `foo/bar`}, `\\host\share\foo\bar`},
			{[]string{`\`}, `\`},
			{[]string{`\`, ``}, `\`},
			{[]string{`\`, `a`}, `\a`},
			{[]string{`\\`, `a`}, `\a`},
			{[]string{`\`, `a`, `b`}, `\a\b`},
			{[]string{`\\`, `a`, `b`}, `\a\b`},
			{[]string{`\`, `\\a\b`, `c`}, `\a\b\c`},
			{[]string{`\\a`, `b`, `c`}, `\a\b\c`},
			{[]string{`\\a\`, `b`, `c`}, `\a\b\c`},
		},
	}

	implTests["Windows"] = append(implTests["Windows"], implTests["Unix"]...)

	impls := map[string]P{
		"Unix":    Unix,
		"Windows": Windows,
	}

	for implName, tests := range implTests {
		t.Run(implName, func(t *testing.T) {
			im := impls[implName]
			for _, test := range tests {
				t.Run(test.path, func(t *testing.T) {
					expected := im.(impl).fromSlash(test.path)
					if p := im.Join(test.elem...); p != expected {
						t.Errorf("Join(%q) = %q, want %q", test.elem, p, expected)
					}
				})
			}
		})
	}
}
