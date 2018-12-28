package paths_test

import (
	"fmt"

	"github.com/apparentlymart/go-paths/paths"
)

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
