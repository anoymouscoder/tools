// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package completion

import (
	"strings"
	"testing"

	. "golang.org/x/tools/gopls/internal/lsp/regtest"
)

func TestPostfixSnippetCompletion(t *testing.T) {
	const mod = `
-- go.mod --
module mod.com

go 1.12
`

	cases := []struct {
		name          string
		posRegExp     string
		before, after string
	}{
		{
			name:      "sort",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var foo []int
	foo.sort
}
`,
			after: `
package foo

import "sort"

func _() {
	var foo []int
	sort.Slice(foo, func(i, j int) bool {
	$0
})
}
`,
		},
		{
			name:      "sort_renamed_sort_package",
			posRegExp: "\n}",
			before: `
package foo

import blahsort "sort"

var j int

func _() {
	var foo []int
	foo.sort
}
`,
			after: `
package foo

import blahsort "sort"

var j int

func _() {
	var foo []int
	blahsort.Slice(foo, func(i, j2 int) bool {
	$0
})
}
`,
		},
		{
			name:      "last",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var s struct { i []int }
	s.i.last
}
`,
			after: `
package foo

func _() {
	var s struct { i []int }
	s.i[len(s.i)-1]
}
`,
		},
		{
			name:      "reverse",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var foo []int
	foo.reverse
}
`,
			after: `
package foo

func _() {
	var foo []int
	for i, j := 0, len(foo)-1; i < j; i, j = i+1, j-1 {
	foo[i], foo[j] = foo[j], foo[i]
}

}
`,
		},
		{
			name:      "slice_range",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	type myThing struct{}
	var foo []myThing
	foo.range
}
`,
			after: `
package foo

func _() {
	type myThing struct{}
	var foo []myThing
	for i, mt := range foo {
	$0
}
}
`,
		},
		{
			name:      "append_stmt",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var foo []int
	foo.append
}
`,
			after: `
package foo

func _() {
	var foo []int
	foo = append(foo, $0)
}
`,
		},
		{
			name:      "append_expr",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var foo []int
	var _ []int = foo.append
}
`,
			after: `
package foo

func _() {
	var foo []int
	var _ []int = append(foo, $0)
}
`,
		},
		{
			name:      "slice_copy",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var foo []int
	foo.copy
}
`,
			after: `
package foo

func _() {
	var foo []int
	fooCopy := make([]int, len(foo))
copy(fooCopy, foo)

}
`,
		},
		{
			name:      "map_range",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var foo map[string]int
	foo.range
}
`,
			after: `
package foo

func _() {
	var foo map[string]int
	for k, v := range foo {
	$0
}
}
`,
		},
		{
			name:      "map_clear",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var foo map[string]int
	foo.clear
}
`,
			after: `
package foo

func _() {
	var foo map[string]int
	for k := range foo {
	delete(foo, k)
}

}
`,
		},
		{
			name:      "map_keys",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var foo map[string]int
	foo.keys
}
`,
			after: `
package foo

func _() {
	var foo map[string]int
	keys := make([]string, 0, len(foo))
for k := range foo {
	keys = append(keys, k)
}

}
`,
		},
		{
			name:      "channel_range",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	foo := make(chan int)
	foo.range
}
`,
			after: `
package foo

func _() {
	foo := make(chan int)
	for e := range foo {
	$0
}
}
`,
		},
		{
			name:      "var",
			posRegExp: "\n}",
			before: `
package foo

func foo() (int, error) { return 0, nil }

func _() {
	foo().var
}
`,
			after: `
package foo

func foo() (int, error) { return 0, nil }

func _() {
	i, err := foo()
}
`,
		},
		{
			name:      "var_single_value",
			posRegExp: "\n}",
			before: `
package foo

func foo() error { return nil }

func _() {
	foo().var
}
`,
			after: `
package foo

func foo() error { return nil }

func _() {
	err := foo()
}
`,
		},
		{
			name:      "var_same_type",
			posRegExp: "\n}",
			before: `
package foo

func foo() (int, int) { return 0, 0 }

func _() {
	foo().var
}
`,
			after: `
package foo

func foo() (int, int) { return 0, 0 }

func _() {
	i, i2 := foo()
}
`,
		},
		{
			name:      "print_scalar",
			posRegExp: "\n}",
			before: `
package foo

func _() {
	var foo int
	foo.print
}
`,
			after: `
package foo

import "fmt"

func _() {
	var foo int
	fmt.Printf("foo: %v\n", foo)
}
`,
		},
		{
			name:      "print_multi",
			posRegExp: "\n}",
			before: `
package foo

func foo() (int, error) { return 0, nil }

func _() {
	foo().print
}
`,
			after: `
package foo

import "fmt"

func foo() (int, error) { return 0, nil }

func _() {
	fmt.Println(foo())
}
`,
		},
		{
			name:      "string split",
			posRegExp: "\n}",
			before: `
package foo

func foo() []string {
	x := "test"
	return x.split
}`,
			after: `
package foo

import "strings"

func foo() []string {
	x := "test"
	return strings.Split(x, "$0")
}`,
		},
		{
			name:      "string slice join",
			posRegExp: "\n}",
			before: `
package foo

func foo() string {
	x := []string{"a", "test"}
	return x.join
}`,
			after: `
package foo

import "strings"

func foo() string {
	x := []string{"a", "test"}
	return strings.Join(x, "$0")
}`,
		},
		{
			name:      "method for struct",
			posRegExp: "meth()",
			before: `
package foo

type Foo struct{}

Foo.meth
`,
			after: `
package foo

type Foo struct{}

func (foo Foo) ${1:}(${2:}) ${3:} {
	$0
}`,
		},
	}

	r := WithOptions(
		Settings{
			"experimentalPostfixCompletions": true,
		},
	)
	r.Run(t, mod, func(t *testing.T, env *Env) {
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				c.before = strings.Trim(c.before, "\n")
				c.after = strings.Trim(c.after, "\n")

				env.CreateBuffer("foo.go", c.before)

				pos := env.RegexpSearch("foo.go", c.posRegExp)
				completions := env.Completion("foo.go", pos)
				if len(completions.Items) != 1 {
					t.Fatalf("expected one completion, got %v", completions.Items)
				}

				env.AcceptCompletion("foo.go", pos, completions.Items[0])

				if buf := env.Editor.BufferText("foo.go"); buf != c.after {
					t.Errorf("\nGOT:\n%s\nEXPECTED:\n%s", buf, c.after)
				}
			})
		}
	})
}
