Most content of this variant is the same as (or similar to) [the main proposal](README.md).
This variant doesn't need the `gen` keyword presented in the the main proposal.

[The contract part](contracts.md) has no differences to the main proposal.

## The generic declaration syntax (of this proposal variant)

The form for a generic declartion looks like:
```
package GenName[inType0, inType1, ..., inTypeN] (
	... // assure lines to constraint type parameters
) {
	... // Go 1 code, such as type, function/method declarations
}
```

The `( ... // assure lines )` part is optional. It doesn't need to present if no constraints are required for type parameters.

In other words, a custom generic are declared as mini-package in a single source file. 

The declaration form differs to the main proposal in several points:
1. It uses the old `package` keyword (instead of inventing a `gen` keyword) to declared a generic.
1. All generic parameters are enclosed in one pair of [].
1. There is not the output part in the declaration.
1. Type paramemters are not followed by the `type` keyword. This means `const` parameters are not supported in this proposal.

Let's view some simple examples to get how to declare and use generics.
Please note that the syntax forms are different when using generic type arguments between functions, types and packages:
* `aGenericFunc(Type1, Type2, ..., TypeN, value1, value2, ...)`, just like `make(ContainerType, length)`
* `aGenericType[Type1][Type2]...[TypeN_1]TypeN`, just like `map[Key]Element`
* `aGenericPackage[Type1, Type2, ..., TypeN]`

## Some simple custom generic examples

### Exampe 1 (declare a generic function):


```
// declaration
package ConvertSlice[OldElement, NewElement] (
	// The contract this genric must satisfy, see the contract page for details.
	assure NewElement(OldElement.value) // convertiable from OldElement to NewElement
) {
	// The name of the function is the same as the generic,
	// so this generic can be called as a function.
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
}

// use it

// Function generic type arguments must be enclosed in a pair of ()
// along with value arguments but the type arguments must be put
// at the beginning of the list.
func strings2Interfaces = ConvertSlice(string, interfacce{})
// This is also allowed.
var _ = ConvertSlice(string, interfacce{})

func main() {
	words := []string{"hello", "bye"}
	fmt.Println(strings2Interfaces(words)...)
	
	// Mix type arguments up with value arguments.
	nums := []int{1, 2, 3}
	fmt.Println(ConvertSlice(int, interfacce{}, nums)...)
}
```

### Example 2 (declare a generic type):

```
// declaration
package List[T] {
	// The name of the type is the same as the generic,
	// so this generic can be called as a generic type.
	type List struct {
		Element T
		Next    *List
	}
	
	func (n *List) Push(e T) *List {...}
	func (n *List) Dump() {...}
}

// use it

// Use the generic type.
// Except the last generic type argument, each of the other ones
// must be enclosed in a pair of [].
type BoolList = List bool

func main() {
	// Use the generic type with argument "int".
	var intList List int 
	intList = intList.Push(123)
	intList = intList.Push(456)
	intList = intList.Push(789)
	intList.Dump()
	
	// Use the generic type with argument "string".
	var strList List string
	strList = intList.Push("abc")
	strList = intList.Push("mno")
	strList = intList.Push("xyz")
	strList.Dump()
}
```

### Example 3 (declare another generic type):

```
// declaration
package TreeMap[Key, Element] (
	// Apply some constraints.
	assure Key.kind == int.kind || Key.kind == string.kind
) {
	// The name of the type is the same as the generic,
	// so this generic can be called as a generic type.
	type TreeMap struct {...}
	func (t *TreeMap) Put(k Key, e Element) {...}
	func (t *TreeMap) Get(k Key) Element {...}
	func (t *TreeMap) Has(k Key) bool {...}
	func (t *TreeMap) Delete(k Key)(Element, bool) {...}
}

// use it

// Use the generic type.
// Except the last generic type argument, each of the other ones
// must be enclosed in a pair of [].
type stringIntTreeMap = TreeMap[string]int

func main() {
	var tm stringIntTreeMap
	tm.Put("Go", 2009)
	...
	
	var tm2 TreeMap[int]bool
	tm2.Put(1, true)
	...
}
```

### Example 4 (declare a generic package):

```
// declaration
package Example[T1, T2] {

	// 
	
	type Bar struct{T1; x *T2}
	func Foo(b *Bar) *T2 {
		return b.x
	}
}

// use it

// All type arguments of a generic package instance
// must be enclosed in the same pair of [].
import alib Example[string, int] // we can use alib as an imported package

func main() {
	var v alib.Bar
	alib.Foo(&v)
}
```

## More examples

```
package Merge[T] {
	func Merge(ss ...[]T) []T {
		n := 0
		for _, s := range ss {
			n += len(s)
		}
		
		r := make([]T, 0, n)
		for _, s := range ss {
			r = append(r, s...)
		}
		
		return r
	}
}

// use it:

var a, b = []string{"hello"}, []string{"world", "!"}
var x, y, z = []{1, 2, 3}, []{5, 5}, []{9}

// By making use of type inference, generic type arguments
// don't need to present here.
var _ = Merge(x, y, z) // [1 2 3 5 5 9]
var _ = Merge(a, b)    // ["hello" "world" "!"]
```

```
package Keys[M] (
	assure M.kind == Map
) {
	type K = M.key

	func Keys(m M) []K {
		if m == nil {
			return nil
		}
		keys := make([]K, 0, len(m))
		for k := range m {
			keys = append(keys, k)
		}
		return keys
	}
}

// use it:

var m = map[string]int{"foo", 1, "bar": 2}

// By making use of type inference, generic type arguments
// don't need to present here.
var _ = Keys(m) // ["foo" "bar"]
```

```
package IncreaseStat[T] (
	assure T.fields.N int
) {
	func IncreaseStat(t *T) {
		t.N++
	}
}

// use it:

type Foo struct {
	x bool
	N int
}

type Bar struct {
	N int
	y string
}

var f Foo
var b Bar

// By making use of type inference, generic type arguments
// don't need to present here.
IncreaseStat(&f) // f.N = 1
IncreaseStat(&b) // b.N = 1
```

```
package Vector[T] {
	type Vector []T
	
	func (v *Vector) Push(x T) {
		*v = append(*v, x)
	}
}
```

```
package Smallest[T] (
	assure T.orderable == true
) {
	func Smallest(s []T) T {
		r := s[0]
		for _, v := range s[1:] {
			if v < r {
				r = v
			}
		}
		return r
	}
}
```

```
package Map[F, T] { 
	func Map(s []F, f func(F) T) []T {
		r := make([]T, len(s))
		for i, v := range s {
			r[i] = f(v)
		}
		return r
	}
}
```

```
package graph

package Graph[Node, Edge] (
	assure Node.methods.Edges() []Edge       // Node must has the specified method
	assure Edge.methods.Nodes() (Node, Node) // Edge must has the specified method
) {
	type Graph struct {
		nodes []Node
		edges []Edge
	}
	
	func (g *Graph) ShortestPath(from, to Node) []Edge { ... }
	
	func (g *Graph) SetNodes(nodes ...Node) { ... }
}

// use it in another package

import graph

type node struct { ... }
func (n *node) Edges() []*edge { ... }
type edge struct { ... }
func (e *edge) Nodes (from, to *node) { ... }

var n1 = &node{ ... }
var n2 = &node{ ... }
var g graph.Graph[*node]*edge
g.SetNodes(n1, n2)
edges := g.ShortestPath(n1, n2)
```



