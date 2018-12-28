package paths_test

import (
	"fmt"

	"github.com/apparentlymart/go-paths/paths"
)

func Example() {
	dir := paths.Windows.Dir(`c:\windows\system32\shell32.dll`)
	fmt.Println(paths.Windows.Join(dir, "../explorer.exe"))

	// Output:
	// c:\windows\explorer.exe
}

func ExampleUnix() {
	start := "/home/fred/.config/bar/baz"
	dir, fn := paths.Unix.Split(start)
	fmt.Println(dir, fn)
	fmt.Println(paths.Unix.Join(dir, fn))
	fmt.Println(paths.Unix.Join(dir, "boz"))

	// Output:
	// /home/fred/.config/bar/ baz
	// /home/fred/.config/bar/baz
	// /home/fred/.config/bar/boz
}

func ExampleWindows() {
	start := `c:\windows\system32\shell32.dll`
	dir, fn := paths.Windows.Split(start)
	fmt.Println(dir, fn)
	fmt.Println(paths.Windows.Join(dir, fn))
	fmt.Println(paths.Windows.Join(dir, "moricons.dll"))

	// Output:
	// c:\windows\system32\ shell32.dll
	// c:\windows\system32\shell32.dll
	// c:\windows\system32\moricons.dll
}

func ExampleSlash() {
	start := "/articles/go-patterns.html"
	dir, fn := paths.Slash.Split(start)
	fmt.Println(dir, fn)
	fmt.Println(paths.Slash.Join(dir, fn))
	fmt.Println(paths.Slash.Join(dir, "go-antipatterns.html"))

	// Output:
	// /articles/ go-patterns.html
	// /articles/go-patterns.html
	// /articles/go-antipatterns.html
}

func ExampleP_Split() {
	fmt.Println(paths.Unix.Split("/foo/bar"))
	fmt.Println(paths.Windows.Split("c:/foo/bar"))

	// Output:
	// /foo/ bar
	// c:/foo/ bar
}

func ExampleP_Join() {
	fmt.Println(paths.Unix.Join("/foo", "bar", "baz"))
	fmt.Println(paths.Windows.Join("c:/", "foo", "bar", "baz"))

	// Output:
	// /foo/bar/baz
	// c:\foo\bar\baz
}

func ExampleP_Clean() {
	fmt.Println(paths.Unix.Clean("a/b/../c"))
	fmt.Println(paths.Windows.Clean("c:/d/../e"))

	// Output:
	// a/c
	// c:\e
}
