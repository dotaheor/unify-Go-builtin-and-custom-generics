
# A solution to unify Go builtin and custom generics

This (immature) solution is extended from
[my many](https://gist.github.com/dotaheor/4b7496dba1939ce3f91a0f8cfccef927)
old [(immature) ideas](https://gist.github.com/dotaheor/c805d221ed86265d6e8bb4f16a714060).

Although the ideas in the solution are still immature for generic implementations,
I think they are really good to unify the appearances and explanations of generics.

In my opinion, the solution has much better readibilities than the generic design in C++, Rust, Java, etc.

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
	var intList *List[int]
	intList = intList.Push(123)
	intList = intList.Push(456)
	intList = intList.Push(789)
	intList.Dump()
	
	var strList *List[string]
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
	func Foo(T) {}
	
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
		}
	}
	
	on dir == true {
		export gen[T type] type {
			... // export a receive-only channel type
		}
	}
	
	export gen[T type] type {
		.../ export a send-only channel type
	}
}
```

The literal representations of directional channel types are also builtin generic privileges.

Operator generic delcations (use `+` as example):
```
gen +[a?, b var] var {
	... // internal implementation
}
```

Operation generics are also builtin generic privileges.
