
# Generic is gen: mini-package - a solution to unify Go builtin and custom generics

This (immature) solution is extended from
[my many](https://gist.github.com/dotaheor/4b7496dba1939ce3f91a0f8cfccef927)
old [(immature) ideas](https://gist.github.com/dotaheor/c805d221ed86265d6e8bb4f16a714060).

Although the ideas in the solution are still immature for generic implementations,
I think they are really good to unify the appearances and explanations of generics.

In my opinion, the solution has much better readibilities than the generic design in C++, Rust, Java, etc.
The biggest advantage of this proposal is the new introduced `gen` elements are much like our familiar `func` element,
and using a `gen` is much like calling a function, which makes the proposal very easy to understand.

Comparing to the current official generic/contract draft,
personally, I think this proposal has the following advantages:
1. supports const generic parameters (the draft only supports types now).
1. consistent looking of builtin and custom generics.
1. the main part of the declaration syntax of generic types and functions is totally Go 1 compatible.
1. uses generics is much like calling functions, so it is easy to understand.

## The generic declaration syntax

Now, there are 5 kinds of code element declarations (except labels) in Go: `var`, `const`, `func`, `type`, and `import`.
This solution adds a new one `gen`, which means a generic declaration.

In the following examples, the generic input constraints are ignored.

The form for exported generic declartions looks like:
```
gen GenName[in0 InputElemKind0][in1 InputElemKind1] ... [inN InputElemKindN] [out OutputElemKind] {
	...
}
```

In the form,
* each `ElemKind` can be any of `var`, `const`, `func`, `type`, `import`, and `gen`.
  However, `var` inputs and outputs are almost never used for it is not much useful.
  `func` and `import` inputs are also not very useful. `const` outputs are also not very useful.
  `gen` outputs are ever rcommended to be used in older versions, but also not much useful in the current version.
  To summarize, now:
	* `InputElemKind` may only be `type` and `const`.
	* `OutputElemKind` may only be `type`, `func`, and `import`.
* the last portion `[out OutputElemKind]` is the output, it may be blank, `[]` but may not omitted (it is not optional).
  If it is not blank, then its surrounding `[]` can be omitted.
* the `[inN InputElemKindN]` portions are the inputs.
  Except the first one, the others are optional but they are may not be blank `[]`.
  The first one `[in0 InputElemKind0]` may be blank `[]`, but it is not optional.
  

_(Note: the syntax used in the current version is a little different from [older versions](https://github.com/dotaheor/unify-Go-builtin-and-custom-generics/tree/890bb969383a8c11a7f17308de8a4020488aeb0f).)_

Please don't be frighten by the new syntax. Their uses are much simpler and intuiative than it looks.
At this point, we just need to know that a `gen` declaration is like a function declaration and we can call it to use it.
The differences are
* the parameters and results of a generic declaration are all code element kinds, instead of value types.
* generic parameters and results are enclosed in `[]` instead of `()`.
* each generic input parameter is enclosed in a separated `[]`.

Let's view some simple examples to get how to declare and use generics.

## Some simple custom generic examples

### Exampe 1 (a single `func` output):

The following `gen` has two inputs (both are `type`) and one output (a `func`).
```
// declaration
gen ConvertSlice[OldElement type][NewElement type] [func] {
	// The only exported function is used as the output of the generic.
	// NOTE: the name of the declared function is not important,
	//       as long as it is exported.
	//       It is recommended to use the same name as the gen, if possible.
	func ConvertSlice(x []OldElement) []NewElement {
		if x == nil {
			return nil
		}
		y := make([]NewElement, 0, len(x))
		for i := range x {
			y = append(y, NewElement(x[i]))
		}
		return y
	}
	
	// There can be more functions declared, but they must be all
	// unexported, for this gen only allows one exported function.
	func anotherUnexportedFunction() {}
}

// use it

func strings2Interfaces = ConvertSlice[string][interfacce{}] // call the gen and bind the output to a function.
var _ = ConvertSlice[string][interfacce{}]                   // this is also allowed.

func main() {
	words := []string{"hello", "bye"}
	fmt.Println(strings2Interfaces(words)...)
	
	nums := []int{1, 2, 3}
	fmt.Println(ConvertSlice[int][interfacce{}](nums)...) // call the gen then call the output function
}
```

### Example 2 (a single `type` output):

The following `gen` has one input (a `type`) and one output (a `type`).
```
// declaration
gen List[T type] type {
	// The only exported type is used as the output of the generic.
	// NOTE: the name of the declared type is not important,
	//       as long as it is exported.
	//       It is recommended to use the same name as the gen, if possible.
	type List struct {
		Element T
		Next    *List
	}
	
	func (n *List) Push(e T) *List {...}
	func (n *List) Dump() {...}
	
	// Some other unexport variables/constants/types/functions
	// can be declared here.
	// ...
	var x = 1
	const n = 128
	func f() {}
	type t struct{}
}

// use it

type BoolList = List[bool] // call the gen then bind the output type to a type alias

func main() {
	var intList List[int] // call the gen then use the output type
	intList = intList.Push(123)
	intList = intList.Push(456)
	intList = intList.Push(789)
	intList.Dump()
	
	var strList List[string] // call the gen then use the output type
	strList = intList.Push("abc")
	strList = intList.Push("mno")
	strList = intList.Push("xyz")
	strList.Dump()
}
```

### Example 3 (another example with only one single `type` output):

The following `gen` has two inputs (both are `type`) and one output (a `type`).
```
// declaration
gen TreeMap[Key type][Element type] [type] {
	// Apply some constraints.
	assure Key.kind == int.kind || Key.kind == string.kind
	
	// The only exported type is used as the output of the generic.
	// NOTE: the name of the declared type is not important,
	//       as long as it is exported.
	//       It is recommended to use the same name as the enclosing gen, if possible.
	type TreeMap struct {...}
	func (t *TreeMap) Put(k Key, e Element) {...}
	func (t *TreeMap) Get(k Key) Element {...}
	func (t *TreeMap) Has(k Key) bool {...}
	func (t *TreeMap) Delete(k Key)(Element, bool) {...}
}

// use it

type stringIntTreeMap = TreeMap[string][int]

func main() {
	var tm stringIntTreeMap
	tm.Put("Go", 2009)
	...
	
	var tm2 TreeMap[int][bool]
	tm2.Put(1, true)
	...
}
```

### Example 4 (an `import` output):

The following `gen` has no inputs but has one output (an `import`). We can think the `gen` outputs a mini-package.
```
// declaration
gen Example[] [import] {

	// For a gen which ouputs an import, all the exported types
	// and functions declared in the gen body will be outputted,
	// their exported names are just their declaration names.
	//
	// For this specified gen, one type and one function will
	// be outputted together in a mini-package.
	
	type Bar struct{}
	func Foo(Bar) {}
}

// use it

import alib = Example[] // we can use alib as an imported package

func main() {
	var v alib.Bar
	alib.Foo(v)
}
```

## How builtin generics are declared

Please note, in this proposal, builtin generics still have some privileges over custom generics.
The names of builtin generics can contain non-identifier letters,
and the represenations of builtin generic uses have more freedom.

The following shown builtin generic declarations are all "look-like", not "exactly-is".

Builtin array and slice declaration:
```
gen array[N const][T type] type {
	... // export an array type
}


gen slice[][T type] type {
	... // export a slice type
}
```

In their uses, the generic identifier `array` and `slice` must be absent. (This is a builtin generic privilege).

Builtin map declaration:
```
gen map[Tkey type][Tvalue type] type {
	assure Tkey.comparable

	... // export a map type
}
```

Builtin channel declaration:
```
gen chan[T type] type {
	type C struct {
		...
	}
	
	// An operator function
	func (c C) <- (v T) {
		// ... send a value v to channel c
	}
	
	// Another operator function
	func <- (c C) (v T) {
		// ... receive a value from channel c
	}
}

// "<-chan" is the name of the gen.
gen <-chan[T type] type {
	type C struct {
		...
	}
	
	func <- (c C) (v T) {
		// ... receive a value from channel c
	}
}

// "chan<-" is the name of the gen.
gen chan<-[T type] type {
	type C struct {
		...
	}
	
	func (c C) <- (v T) {
		// ... send a value v to channel c
	}
}
```

The literal representations of directional channel types are also builtin generic privileges.

Operator function generics are also builtin generic privileges.

(BTW, can we make `map` and `chan` become non-keywords for better consistency?)

Built-in `new` and `make` generic function declarations:
```
gen new[T type] func {
	// must be exported
	func New() {
		var x T
		return &x
	}
}

gen make[T type] func {
	// apply some constraints
	assure T.kind == ([]T).kind || T.kind == (map[int]T).kind

	// must be exported
	func Make(params ...int) T {
		// ...
	}
}
```

## More about generic calls

### Cyclic calls between `gen`s declared in the same package are allowed

But `gen`s declared from different packages may not depend on cyclicly.

(Note, this part of the design is still very not mature.)

### For a `gen` with single `type` output, in its calls, the `[]` surrounding the last generic arguments may be omitted

If we observe builtin generic syntax carefully, we will find that the last generic arguments are not enclosed in `[]`.
For example: `array[5]int`, `slice[]int`, `map[string]int`, `chan int`.
(Surely, the `array` and `slice` identifier must be ommited in uses, so below for details.)

We can apply this same rule for custom generics.
For example, for the generic type declared in the above example 2 and 3, their calls may be
```
type BoolList = List bool
type stringIntTreeMap = TreeMap[string]int
```

(Is it good to make this rule mandatory?)

### For a `gen` with single `func` output, in its calls, the generic arguments may be inserted (at the beginning) into general argument list

For example, the built-in `new` and `make` generic, they may be called with two forms:
```
var x = new[string]() // different from Go 1
var y = new(string)   // same as Go 1

var m1 = make[map[int]string]() // different from Go 1
var s1 = make[[]int](100)       // different from Go 1
var m2 = make(map[int]string)   // same as Go 1
var s2 = make([]int, 100)       // same as Go 1
```

Call the generic function shown in example 1 shown above:
```
func main() {
	words := []string{"hello", "bye"}
	fmt.Println(strings2Interfaces(string, interfacce{}, words)...)
	
	nums := []int{1, 2, 3}
	fmt.Println(ConvertSlice(int, interfacce{}, nums)...)
}
```

Generic arguments can be inferred from general arguments' types.
For example, the generic function in example 1 may also be called as:
```
func main() {
	words := []string{"hello", "bye"}
	fmt.Println(strings2Interfaces(words)...)
	
	nums := []int{1, 2, 3}
	fmt.Println(ConvertSlice(nums)...)
}
```

Compilers can infer the first generic argument as the element type of `words` or `nums`,
and infer the second generic argument as the element type of the only parameter (`[]interface{}`) of `fmt.Println`.

## Contracts

Please read [the constraint expressions part](contracts.md) of this proposal for how to write basic generic parameter constraints.

Constraint expressions must be assured in the beginning of a `gen` declartion body to be used as the contract of the `gen`.

All the constraints assured (directly or indirectly) in a `gen` compose of a contraint set, or a contract.
We can use the name of the `gen` to represent the contract.
A `gen` with a blank output acts as a pure contract.

When using a `gen` as a contract, prefix the `assure` keyword to its call.
For example, the built-in generic `make` is called as a contract.
``` 
gen Convert [T1 type][T2 type] func {
	// Constraint T1 must be a slice or map.
	assure make[T2]

	// Constraint T1 must be a slice type
	// and element values of T2 may be
	// converted to element type of T1.
	assure T1.kind == []int.kind && T1.element(T2.element.value)

	func Convert(kvs T2) T1 {
		vs := make(T1, 0, len(kvs))
		for _, v := range kvs {
			vs = append(vs, v)
		}
		return v
	}
}

```

## Comparisons between this proposal and the latest official draft

* Map: [the draft](https://go-review.googlesource.com/c/go/+/187317/3/src/go/parser/testdata/map.go2) vs. [this proposal](https://github.com/dotaheor/unify-Go-builtin-and-custom-generics/blob/master/examples/src/go/parser/testdata/map.go2) (and [another multi-output version](https://github.com/dotaheor/unify-Go-builtin-and-custom-generics/blob/master/examples/src/go/parser/testdata/map2.go2)).
* Set: [the draft](https://go-review.googlesource.com/c/go/+/187317/3/src/go/parser/testdata/set.go2) vs. [this proposal](https://github.com/dotaheor/unify-Go-builtin-and-custom-generics/blob/master/examples/src/go/parser/testdata/set.go2).
* [More](https://github.com/dotaheor/unify-Go-builtin-and-custom-generics/tree/master/examples/src/go).

## More examples

(todo)

