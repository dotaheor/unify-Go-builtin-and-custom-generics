
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
1. support const generic parameters (the draft only supports types now).
1. consistent looking of builtin and custom generics.
1. the body of a generic declaration is totally Go 1 compatible.
1. using generics is much like calling functions, so it is easy to understand.
  * supporting multiple outputs as a mini package.
  * supporting generic closure.

## Overview of this solution

Now, there are 5 kinds of code element declarations (except labels) in Go: `var`, `const`, `func`, `type`, and `import`.
This solution adds a new one `gen`, which means a generic declaration.

In the following examples, the generic input constraints are ignored.

A generic declartion looks like

```
gen GenName[in0 InputElemKind0, in1 InputElemKind1, ...] [out OutputElemKind] {
	...
}
```

where each `ElemKind` can be any of `var`, `const`, `func`, `type`, `import`, and `gen`.
(However, `var` inputs and outputs are almost never used for it is not much useful.
`func` and `import` inputs are also not very useful.)
The number of the outputs of a `gen` decalration can be zero or one.
If a `gen` has no outputs, then it is viewed as a pure contract.

Note: [the old revision](https://github.com/dotaheor/unify-Go-builtin-and-custom-generics/blob/51b200e5d0f959f8a0ae2110d52d528b9ad393a4/README.md) of this proposal permits multiple
outputs, which is prohibited by the current revision.

From the declaration, we can see that the `gen` declaration syntax form is very like a function declaration.
The difference is the parameters and results of a generic declaration are all code element kinds,
instead of value types.

## Some simple custom generic examples

### Exampe 1 (a `func` output):
```
// declaration
gen ConvertSlice[OldElement, NewElement type] [func] {
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

func strings2Interfaces = ConvertSlice[string, interfacce{}]

func main() {
	words := []string{"hello", "bye"}
	fmt.Println(strings2Interfaces(words)...)
	
	nums := []int{1, 2, 3}
	fmt.Println(ConvertSlice[int, interfacce{}](nums)...)
}
```

Note: by using [the simplifed form](#some-simple-single-output-gens-can-be-simplified) mentioned below,
the above `gen` can also be declared as:
```
gen ConvertSlice[OldElement, NewElement type] func (x []OldElement) []NewElement {
	if x == nil {
		return nil
	}
	y := make([]NewElement, 0, len(x))
	for i := range x {
		y = append(y, NewElement(x[i]))
	}
	return y
}
```

### Example 2 (a `type` output):
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

func main() {
	var intList List[int]
	intList = intList.Push(123)
	intList = intList.Push(456)
	intList = intList.Push(789)
	intList.Dump()
	
	var strList List[string]
	strList = intList.Push("abc")
	strList = intList.Push("mno")
	strList = intList.Push("xyz")
	strList.Dump()
}
```

### Example 3 (an `import` output):

If a `gen` outputs an `import`, we can think the `gen` outputs a mini-package.

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

BTW, here are two comparisons between this proposal and the latest official draft:
1. Map: [the draft](https://go-review.googlesource.com/c/go/+/187317/3/src/go/parser/testdata/map.go2) vs. [this proposal](https://github.com/dotaheor/unify-Go-builtin-and-custom-generics/blob/master/examples/src/go/parser/testdata/map.go2).
1. Set: [the draft](https://go-review.googlesource.com/c/go/+/187317/3/src/go/parser/testdata/set.go2) vs. [this proposal](https://github.com/dotaheor/unify-Go-builtin-and-custom-generics/blob/master/examples/src/go/parser/testdata/set.go2).

### Example 4 (a `gen` output):
```
// declaration
gen TreeMap[Key type] [gen] {
	// The only exported gen is used as the output of the generic.
	// NOTE: the name of the declared gen is not important,
	//       as long as it is exported.
	//       It is recommended to use the same name as the enclosing gen, if possible.
	gen TreeMap[Element type] type {
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
}

// use it

type stringIntTreeMap = TreeMap[string][int]

func main() {
	var tm stringIntTreeMap
	tm.Put("Go", 2009)
	...
}
```

We can call the inner `TreeMap` generic use case as a generic closure.
The use in the above example can also be called as a generic call chain.

Note: by using [the simplifed form](#some-simple-single-output-gens-can-be-simplified) mentioned below,
the above `gen` can also be declared as:
```
gen TreeMap[Key type] gen [Element type] type {
	type Tree struct {...}
	func (t *Tree) Put(k Key, e Element) {...}
	func (t *Tree) Get(k Key) Element {...}
	func (t *Tree) Has(k Key) bool {...}
	func (t *Tree) Delete(k Key)(Element, bool) {...}
}
```

## If the last generic in a generic call chain has only one input, then the `[]` surrounding the argument can be omitted.

If we observe builtin generic syntax carefully, we will find that the last generic arguments are not enclosed in `[]`.
For example: `array[5]int`, `slice[]int`, `map[string]int`, `chan int`.
(Surely, the `array` and `slice` identifier must be ommited in uses, so below for details.)

We can apply this same rule for custom generics.
For example, in the last example above, the generic use can be
```
type stringIntTreeMap = TreeMap[string]int
```

which is like the builtin `map` generic.

## How builtin generics are declared

Please note, in this solution, builtin generics still have some privileges.
The names of builtin generics can contain non-identifier letters,
and the represenations of builtin generic uses have more freedom.

The following shown builtin generic declarations are all "look-like", not "exactly-is".

Builtin array and slice declaration:
```
gen array[N const] gen {
	gen Array[T type] type {
		... // export an array type
	}
}

gen slice[] gen {
	gen Slice[T type] type {
		... // export an array type
	}
}
```

In their uses, the generic identifier `array` and `slice` must be absent. (This is a builtin generic privilege).

Builtin map declaration:
```
gen map[Tkey type] gen {
	gen Map[Tvalue type] type {
		... // export a map type
	}
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

(BTW, can we make `map` and `chan` become non-keywords?)

## `gen`s are also contracts

For example, the following no-outputs `gen` acts as a pure contract.

```
gen viaStrings(To, From type) {
	func _() {
		var t To
		var f From
		var x string = f.String()
		t.Set(string("")) // could also use t.Set(x)
	}
}
```

Generics with outputs can also be viewed as (non-pure) contracts.

The following `gen` implies the above contract.
```
gen SetViaStrings[To, From type] func {
	func Exported(s []From) []To {
		r := make([]To, len(s))
		for i, v := range s {
			r[i].Set(string(v.String()))
		}
		return r
	}
}
```

Another example: the above `TreeMap` generic can be delcared as
```
gen TreeMap[Tkey type] gen {
	comparable[Tkey] // call another contract to tighten the requirements for Tkey.
	                 // "comparable" is a builtin gen/contract.
	
	gen TreeMap[Tvalue type] type {
		... // export a tree map type
	}
}
```

Personally, I think we should use a hybrid way of the official contract draft
version 1 and version 2 to define contracts.

It would be great if it is possible that compilers can construct a contract from each `gen` block automatically.
A contrat is a rule set after all.
For many `gen` blocks, their contracts (the rule sets) are verbose to clearly write down manually,
for each of the sets might contain many rules and some of the rules might be subtle.
However, this is not hard job for computers.

## What is the meaningfullness of calling a contract generic in another generic?

For example, in the last example, the `TreeMap` calls the `comparable` generic.
However, its only exported `gen` implementation might not require the `Tkey`
type is comparable, which means, the TreeMap can support slice/func/map types
as key types, however, this is temporarily prohibited for the `comparable` generic
is called. Yes, this is exactly the meaningfullness of calling extra generics
as contracts to add more constraints in a generic to accept less valid inputs
than a generic implementation can actually support. This is because some supported
types might not be tested fully or other reasons. In other words, callig some
looks-irrelevant contracts in a `gen` tightens the conditions of the `gen`.

## Some simple single output `gen`s can be simplified

If the single output is a type or a function,
then we can simplify the `gen` declaration.
For example,

```
gen identity[T type] func {
	func Identity(x T) T {
		return x
	}
}

gen set[T type] type {
	type Set map[T]struct{}
}
```

can be simplifed as

```
gen identity[T type] func (x T) T {
	return x
}

gen set[T type] type map[T]struct{}
```

## Remaining problems

The above efforts don't unify the `new` and `make` builtin generic functions well.

A new generic must be declared (by using the just mentioned simplifed form above) and used as
```
gen new[T type] func() *T {
	var x T
	return &x
}

// use it:

var x = new[string]() // different from Go 1
```

and make:

```
gen make[T type] func {
	sameKind[T, []T] || sameKind[T, map[int]T] // sameKind is a builtin contract
	
	func Make(params ...int) T {
		// ....
	}
}

// use it:

var m = make[map[int]string]() // different from Go 1
var s = make[[]int](100)       // different from Go 1
```

To make the unification complete, we can add a rule that allows
merging the generic argument list into the genernal value argument list in
a generic function call, if the generic argument list contains only types.
By this rule, the above `new` and `make` calls can be written as:
```
var x = new(string)
var m = make(map[intstring)
var s = make([]int, 100)
```

