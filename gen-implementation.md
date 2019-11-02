

# `gen` means both generic and generator

We can write a custom code generator to implement the generic proposal.
However, making the generator built into Go compiler will improve user experience much.

### An example to show how the code generation works

For a `gen` declared in a package `foo`, all the generated types and functions are put in a geneated sub-package of `foo`.

For example, for the following example:
```
// foo/gen-a.go
package foo

gen GenA[T1 type][T2 type] [import] {
	type s = struct {
		x T1
	}
	type Bar struct{
		s
		y T2
	}
	func (b Bar) S() s {
		return b.s
	}
	func Foo(Bar) {}
}

// foo/gen-b.go
package foo

gen GenB[T type] [type] {
	type Gen2 = []T
}
```

```
// main.go
package main

import "foo"

import gen_a_1 = foo.Gen1[int][string]
import gen_a_2 = foo.Gen1[int][bool]

type Bar1 = gen_a_1.Bar
type Bar2 = gen_a_2.Bar

type ListRune = foo.Gen2[rune]
type ListByte = foo.Gen2[byte]

func main() {
	var b1 Bar1
	var b2 Bar2
	_ = b1.S() == b2.S() // legal
}
```

is equivalent to the following generated result:

```
// foo/GenA
package GenA // It is better to put this generated package in an
             // internal package to avoid being imported in user code.

type s_int = struct {
	x int
}

type Bar_int_string struct{
	s_int
	y string
}

func (b Bar_int_string) S() s_int {
	return b.s
}

func Foo__Bar_int_string(Bar_int_string) {}

type Bar_int_bool struct{
	s_int
	y bool
}

func (b Bar_int_string) S() s_int {
	return b.s
}

func Foo__Bar_int_bool(Bar_int_bool) {}
```

```
// foo/GenA/_1
package _1 // Gen1[int][string]

import "foo/GenA"

type Bar = GenA.Bar_int_string
func Foo = GenA.Foo__Bar_int_string
```

```
// foo/GenA/_2
package _2 // Gen1[int][bool]

import "foo/GenA"

type Bar = GenA.Bar_int_bool
func Foo = GenA.Foo__Bar_int_bool
```

```
// foo/GenB
package GenB // It is better to put this generated package in an
             // internal package to avoid being imported in user code.

type Type_rune = []rune

type Type_byte = []byte
```

```
// foo/GenB/_1
package _1 // Gen2[rune]

type Type = Type_rune
```

```
// foo/GenB/_2
package _2 // Gen2[byte]

type Type = Type_byte
```

```
// main.go
package main

import gen_a_1 "foo/GenA/_1"
import gen_a_2 "foo/GenA/_2"
import gen_b_1 "foo/GenB/_1"
import gen_b_2 "foo/GenB/_2"

type Bar1 = gen_a_1.Bar
func Foo1 = gen_a_1.Foo
type Bar2 = gen_a_2.Bar
func Foo2 = gen_a_2.Foo

type ListRune = gen_b_1.Type
type ListByte = gen_b_2.Type

func main() {
	var b1 Bar1
	var b2 Bar2
	_ = b1.S() == b2.S() // legal
}
```

In generating, to avoid too long type and function generated names, we can use the `Foo_aHaSh` form as names instead.

----

Several days after I wrote down the above content, I suddenly realized that the real intention of
[the alias declarations propsoal](https://github.com/golang/go/issues/16339) might be for generics,
instead of the claimed large-scale refactoring in that proposal. Maybe, maybe, I will stop here.

