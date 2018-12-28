package paths

import "testing"

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
