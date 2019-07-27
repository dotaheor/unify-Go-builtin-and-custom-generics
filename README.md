
# Generic is gen: super function - a solution to unify Go builtin and custom generics

This (immature) solution is extended from
[my many](https://gist.github.com/dotaheor/4b7496dba1939ce3f91a0f8cfccef927)
old [(immature) ideas](https://gist.github.com/dotaheor/c805d221ed86265d6e8bb4f16a714060).

Although the ideas in the solution are still immature for generic implementations,
I think they are really good to unify the appearances and explanations of generics.

In my opinion, the solution has much better readibilities than the generic design in C++, Rust, Java, etc.
The biggest advantage of this proposal is the new introduced `gen` elements are much like our familiar `func` element, which makes the proposal very easy to understand.

### Overview of this solution

Now, there are 5 kinds of code element declarations (except labels) in Go: `var`, `const`, `func`, `type`, and `import`.
This solution adds a new one `gen`, which means a generic declaration.

In the following examples, the generic input constraints are ignored.

A generic declartion looks like


```
gen GenName[in0 InputEleKind0, in1 InputEleKind1, ...] [out0 OutputEleKind0, out1 OutputEleKind1, ...] {
	...
}
```

The form is very like a function declaration.
The difference is the parameters and results of a generic declaration are all code element kinds.
In other words, the parameters and results of a generic declaration can be `var`, `const`, `func`, `type`, `import`, and `gen`.

### Some simple custom generic examples

Exampe 1 (single `func` output):
```
// declaration
gen ConvertSlice[SliceElement, NewElement type] [func] {
	func Convert(x []SliceElement) []NewElement {
		if x == nil {
			return nil
		}
		y := make([]NewElement, 0, len(x))
		for i := range x {
			y = append(y, NewElement(x[i]))
		}
		return y
	}
	
	export Convert
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

Example 2 (single `type` output):
```
// declaration
gen List[T type] type {
	type node struct {
		Element T
		Next    *node
	}
	
	func (n *node) Push(e T) *node {...}
	func (n *node) Dump() {...}
	
	export node
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

Example 3 (a single `import` output):
```
// declaration
gen Example[] [import] {
	type Bar struct{}
	func Foo(Bar) {}
	
	export {
		Bar,
		Foo,
	}
}

// use it

import alib = Example[]

func main() {
	var v alib.Bar
	alib.Foo(v)
}
```

Example 4 (a single `gen` output):
```
// declaration
gen TreeMap[Key type] [gen] {
	export gen[Element type] type {
		type Tree struct {...}
		func (t *Tree) Put(k Key, e Element) {...}
		func (t *Tree) Get(k Key) Element {...}
		func (t *Tree) Has(k Key) bool {...}
		func (t *Tree) Delete(k Key)(Element, bool) {...}
		
		export Tree
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

We can call the `TreeMap` generic use case as a generic chain with two generics.
The uses in the above three other examples can also be called as generic chain,
but each of them only uses one generic.

### If the last generic in a generic chain use has only one input, then the `[]` surrounding the argument can be omitted.

For example, in the last example above, the generic use can be
```
type stringIntTreeMap = TreeMap[string]int
```

which is like the builtin `map` generic.

### A generic input can be delcared as optional if the input is the only input of a geneirc.

For example, for the generic declared as
```
gen Something[N? const] [outT type] {
	// a possible implementaion
	on absent(N) {
		...
	} else {
		...
	}
	...
}
```

we can use it as
```
type X = Something[16]
type Y = Something[]
```

### How builtin generics are declared

Please note, in this solution, builtin generics still have some privileges.
The names of builtin generics can contain non-identifier letters,
and the represenations of builtin generic uses have more freedom.

The following shown builtin generic declarations are all "look-like", not "exactly-is".

Builtin array and slice declaration:
```
gen array[N? const] gen {
	on absent(N) {
		export gen[T type] type {
			... // export a slice type
		}
	} 
	
	export gen[T type] type {
		... // export an array type
	}
}
```

In it uses, the generic identifier `array` must be absent. (This is a builtin generic privilege).

Builtin map declaration:
```
gen map[Tkey type] gen {
	export gen[T type] type {
		... // export a map type
	}
}
```

Builtin channel declaration:
```
gen chan[dir? const] gen {
	on absent(dir) {
		export gen[T type] type {
			... // export a bi-directional channel type
		
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
			
			export C
		}
	}
	
	on dir == true {
		export gen[T type] type {
			... // export a receive-only channel type
		
			type C struct {
				...
			}
			
			func <- (c C) (v T) {
				// ... receive a value from channel c
			}
			
			export C
		}
	}
	
	export gen[T type] type {
		... // export a send-only channel type
		
		type C struct {
			...
		}
		
		func (c C) <- (v T) {
			// ... send a value v to channel c
		}
		
		export C
	}
}
```

The literal representations of directional channel types are also builtin generic privileges.

Operator generic delcations (use `+` as example):
```
gen +[Ta?, Tb type] func {
	... // internal implementation
	on kind(Tb) == string {
		...
	}
	on kind(Tb) == int {
		...
	}
	...
}
```

Operator generics are also builtin generic privileges.

### The above shown operator generic and optional generic inputs might be not good ideas

It might be better to split slice and array as two different generics by not using optional inputs.

It might also be better to split the channel generic described above as three different generics: `chan`, `chan<-` and `<-chan`.

And it might be best not to support operator generics.

### Work with contracts

The contract idea proposed in the [Go 2 draft](https://go.googlesource.com/proposal/+/master/design/go2draft-contracts.md)
is a great idea. However, I think it can be improved.
For example, there is a contract defined in the draft as
```
contract viaStrings(t To, f From) {
	var x string = f.String()
	t.Set(string("")) // could also use t.Set(x)
}
```

I think the values `t` and `f` shouldn't appear in the contract prototype.
It would be better to define the contract as the following one:
```
contract viaStrings(To, From) {
	var t To
	var f From
	var x string = f.String()
	t.Set(string("")) // could also use t.Set(x)
}
```

Further more, I think the `contract` keyword is unecessary.
In fact, the above contract can be defined as no-outputs `gen` instead.
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

Yes, no-outputs generics will be viewed as pure contracts,
and generics with outputs can also be viewed as (non-pure) contracts.

The following is a re-written of the `SetViaStrings` generic function shown in the Go 2 draft.
```
gen SetViaStrings[To, From type] func {
	viaStrings[To, From] // call the contract (another generic)
	
	export func(s []From) []To {
		r := make([]To, len(s))
		for i, v := range s {
			r[i].Set(v.String())
		}
		return r
	}
}
```

Another example: the builtin map generic can be delcared as
```
gen TreeMap[Tkey type] gen {
	comparable[Tkey]
	
	export gen[T type] type {
		... // export a tree map type
	}
}
```

where `comparable` a builtin contract (a builtin `gen`).

### What is the meaningfullness of calling a contract generic in another generic?

For example, in the last example, the `TreeMap` calls the `comparable` generic.
However, its only exported `gen` implementation might not require the `Tkey`
type is comparable, which means, the TreeMap can support slice/func/map types
as key types, however, this is temporarily prohibited for the `comparable` generic
is called. Yes, this is exactly the meaningfullness of calling extra generics
as contracts to add more constraints in a generic to accept less valid inputs
than a generic implementation can actually support. This is because some supported
types might not be tested fully or other reasons. In other words, callig some
looks-irrelevant contracts in a `gen` tightens the conditions of the `gen`.

### The `export` keyword can be removed from this proposal

In fact, the `export` keyword is also not very essential. We can comply with the current Go conventions. If a `gen` only exports one `type` or `func` element, then there can only be exactly one type or function which is declared as exported (first letter is upper case) in the `gen` body. If a `gen` exports an `import`, then there can be multiple elements declared as exported. A generic declartion will look like

```
gen GenName[in0 InputEleKind0, in1 InputEleKind1, ...] [out OutputEleKind] {
	...
}
```

### Some somple single output `gen` can be simplified

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

gen set[T type] map[T]struct{}
```

or (anonymous generics)

```
gen [T type] func identity (x T) T {
	return x
}

gen [T type] type set map[T]struct{}
```

### Remaining problems

The above efforts don't unify the `new` and `make` builtin generic functions well.

