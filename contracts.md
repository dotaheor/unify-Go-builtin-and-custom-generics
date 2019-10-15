
There are two kinds of ways to describe a constraint:
1. use constraint expressions.
1. through contract calls.

## The `assure` keyword

Each `assure` line describe a constraint (see following section for examples).

_(Other candidates to replace the `assure` keyword: `require`, `must`, etc.)_

## Properties

`type` properties:
* `T.kind`, means the kind of the type represented by `T`.
* `T.value`, means an unspecified value of the type represented by `T`.
* `T.name`, `""` for `T` represents neither a defined type or a type alias.
* `T.defined`, whether or not the type represetned by `T` represents a defined type.
   Note: `T.name != "" && !T.defined` is `true` means `T` is a type alias.
* `T.comparable`, whether or not the type represetned by `T` represents a comparable type.
* `T.signed`: whether or not the type represetned by `T` is a signed numeric type.
   (There must be an aforementioned contract constraints `T` to represent
   an integer or floating-point type to use this property.)
* `T.underlying`: the underlying type of the type represetned by `T`.
* `T.base`: the base type of the pointer type represetned by `T`.
   (There must be an aforementioned contract constraints `T` to represent
   a pointer type).
* `T.key`: the key type of the map type represetned by `T`.
   (There must be an aforementioned contract constraints `T` to represent
   a map type).
* `T.element`: the element type of the type represetned by `T`.
   (There must be an aforementioned contract constraints `T` to represent
   an array, slice, map, or channel type to use this property.)
* `T.receivable`: whether or not the type represetned by `T` represents a receivable channel type.
   (There must be an aforementioned contract constraints `T` to represent
   a channel type).
* `T.sendable`: whether or not the type represetned by `T` represents a sendable channel type.
   (There must be an aforementioned contract constraints `T` to represent
   a channel type).
* `T.methsets`: the method set of the type represetned by `T`.
* `T.fields`: the field set of the type represetned by `T`.
   (There must be an aforementioned contract constraints `T` to represent
   a struct type).
* `T.selectors`: the selector set (both methods and fields) of the type represetned by `T`.
* `T.variadic`: whether or not the type represetned by `T` represents a variadic function type.
   (There must be an aforementioned contract constraints `T` to represent
   a function type).
* `T.inputs.count`: the number of parameters of the function type represetned by `T`.
   (There must be an aforementioned contract constraints `T` to represent
   a function type).
* `T.inputs.0`: the first parameter type of the function type represetned by `T`.
   (There must be an aforementioned contract constraints `T` to represent
   a function type).
* `T.outputs.count`: the number of results of the function type represetned by `T`.
   (There must be an aforementioned contract constraints `T` to represent
   a function type).
* `T.outputs.0`: the first result type of the function type represetned by `T`.
   (There must be an aforementioned contract constraints `T` to represent
   a function type).
* `T.length`: the length of the array type represetned by `T`.
   (There must be an aforementioned contract constraints `T` to represent
   an array type).

`const` properties:
* `C.name`: the name of the constant represented by `C` is signed.
* `C.typed`: whether or not the constant represented by `C` is typed.
   (Maybe, it is good to require all generic constants must be typed.)
* `C.type`: the type or default type of a constant represented by `C`.

(`var`, `func`, `import` and `gen` can also be used as contract parameters/arguments,
but doing this will bring much complexity. So this is not supported temporarily.)

## Built-in contract expressions

All contract expressions are built-in. Custom contract expressions are not supported.

Simple ones:
```
assure T.defined
assure T.comparable
assure N > M // N and M must be two generic consts.
```

Specify the type represetned by `T` muse have a specified method:
```
assure T.methods.M func(string) int
```

Specify the struct type represetned by `T` muse have a specified field:
```
assure T.fields.X int
```

Specify the type represetned by `T` muse have a specified selector:
```
assure T.selectors.X int
assure T.selectors.F func(string) int // F can be either a method or a field of a function type
```

Some more complex ones:
```
assure T.methods (
		.M1 func(string) int
		.M2 func(..int) (string, error)
		Ty.methods // embed a method set
	)
assure T.fields (
		.X int
		.Y string
		Tx.fields // embed a field set
	)
assure T.selectors (
		.X int
		.F func(string) int
		Tz.selectors // embed a selector set
	)
```

## Built-in contract calls

#### `comparable[Tx, Ty type]`

Whether the values of the input types are comparable with each other.

Examples:
```
assure comparable[Ta, Tb]
assure comparable[Tx.base, map[Ty]int]
```

#### `assignable[Td, Ts type]`

Whether the values of `Ts` can be assigned to type `Td`.

Examples:
```
assure assignable[[]int, Ta]
assure assignable[interface{M()}, Tx]
```

#### `convertible[Td, Ts type]`

Whether the values of `Ts` can be converted to type `Td`.

Examples:
```
assure convertible[[]int, Ta]
assure convertible[interface{M()}, Tx]
```

#### identical[Tx, Ty, T ...type]

Whether or not the input types are identical types.

For example,
```
assure identical[Tx.base, Ty.element, Tz.key]
```

#### distinct[Tx, Ty, T ...type]

Whether or not the input types are distinct types.

For example,
```
assure distinct[Tx.element, Ty.inputs.0, Tz.outputs.1]
```

#### sameKind[Tx, Ty, T ...type]

Whether or not the input types belong to the same kind.

For example,
```
assure sameKind[map[any]any, Tx, Ty]
```

Here, let's suppose [`any` is an alias of `interface{}`](https://github.com/golang/go/issues/33232).


#### anyKind[Tx type, Ts ...type]

Whether or not first input type is any kind of the kinds of the following input types.

For example,
```
assure anyKind[Ta]                     // Ta can be any kind
assure anyKind[Tx, string]             // Tx must be of string kind
assure anyKind[Ty, string, int, int64] // Ty can be any of string, int or int64 kind
```

#### `impelements[Tx, Ty type]`

Whether the type represented by `Tx` implements the interface type represented by `Ty`.

Examples:
```
assure impelements[Ta, interface{M1()}]
assure impelements[Ta.element, Tm.key]
```

_(This contract is a little overlapping with the `T.methods` expression mentioned above.)_

## Custom contract calls

Custom `gen` declarations can be used as custom contracts.

An Example:
```
gen Min(T type) func {
	assure T.comparable // or: assure comparable[T, T]
	
	func Min(x, t T) T {
		if x < y {
			return x
		} else {
			return y
		}
	}
}

gne Max(T type) func {
	assure Min[T] // Apply the contract of the Min gen to the Max gen.
	              // The effect is the same as appending all the assure
		      // lines in the Min gen to the Max gen.
	
	func Max(x, t T) T {
		if x > y {
			return x
		} else {
			return y
		}
	}
}
```

In fact, a no-outputs `gen` is totally for constraints definition purpose.
```
gen ConditionSet[Tx, Ty, Tz type] {
	type S = []Tx
	assure assignable[S, Ty]
	assure sameKind(Tz, func()]
	type M struct {Y Ty}
	assure convertible[Tz.inputs.0, M]
	// ... more conditions
}
```

The `ConditionSet` `gen` can be called in other `gen` declartions to constraint some types.
```
gen MyGen[Ta, Tb, Tc, Td] [type] {
	assure ConditionSet[Tb, Tc, Td]
	
	type MyGen struct {
		A Ta
		B Tb
		C Tc
		D Td
	}
}
```


