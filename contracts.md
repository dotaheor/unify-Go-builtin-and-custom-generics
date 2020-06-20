
This contract proposal is a part of [the generic design](README.md).
It can also be used with other generic type and function declartion proposals.

This proposal suggests using some familar expressions used in daily Go programming
to constrain generic (type and const, etc) parameters.

Personally, I think this design combines the advantages of the contract draft v1 and v2.
* v1: prone to change constraints accidentally, non-precise, ambiguities existing, but **almost possibility complete**.
* v2: **formal, precise**, but possibility limited.

## Properties

`type` properties:
* `T.kind`, means the kind of the type represented by `T`.
* `T.value`, means an unspecified (addressable) value of the type represented by `T`.
* `T.name`, means the name of the type represented by `T`. `""` for unnamed types.
* `T.signed`: whether or not the type represetned by `T` is a signed numeric type.
   (Not a very elementary property.
   There must be an aforementioned contract constrainting `T` to represent
   an integer or floating-point type to use this property.)
* `T.orderable`, whether or not the values of the type represented by `T` can be compared with `<` and `>`, etc.
   (Not a very elementary property).
* `T.comparable`, whether or not the type represented by `T` represents a comparable type.
* `T.embeddable`, whether or not ~~the type represented by~~ `T` is embeddable.
* `T.base`: the base type of the pointer type represented by `T`.
   (There must be an aforementioned contract constrainting `T` to represent
   a pointer type).
* `T.key`: the key type of the map type represented by `T`.
   (There must be an aforementioned contract constrainting `T` to represent
   a map, slice, array or string type. If `T` is a slice, array or string type,
   then `T.key` is `int`).
* `T.element`: the element type of the type represented by `T`.
   (There must be an aforementioned contract constrainting `T` to represent
   an array, slice, map, or channel type to use this property.)
* `T.length`: the length of the array type represented by `T`.
   (There must be an aforementioned contract constrainting `T` to represent
   an array type).
* `T.receivable`: whether or not the type represented by `T` represents a receivable channel type.
   (There must be an aforementioned contract constrainting `T` to represent
   a channel type).
* `T.sendable`: whether or not the type represented by `T` represents a sendable channel type.
   (There must be an aforementioned contract constrainting `T` to represent
   a channel type).
* `T.methods`: the method set of the type represented by `T`.
* `T.fields`: the field set of the type represented by `T`.
   (There must be an aforementioned contract constrainting `T` to represent
   a struct type).
* `T.selectors`: the selector set (both methods and fields) of the type represented by `T`.
* `T.variadic`: whether or not the type represented by `T` represents a variadic function type.
   (There must be an aforementioned contract constrainting `T` to represent
   a function type).
* `T.inputs.count`: the number of parameters of the function type represented by `T`.
   (There must be an aforementioned contract constrainting `T` to represent
   a function type).
* `T.inputs.0`: the first parameter type of the function type represented by `T`.
   (There must be an aforementioned contract constrainting `T` to represent
   a function type).
* `T.outputs.count`: the number of results of the function type represented by `T`.
   (There must be an aforementioned contract constrainting `T` to represent
   a function type).
* `T.outputs.0`: the first result type of the function type represented by `T`.
   (There must be an aforementioned contract constrainting `T` to represent
   a function type).

`const` properties (the offical contract draft 2 doesn't support `const` generic parameters now,
but the properties are shown here anyway):
* ~~`C.name`: the name of the constant represented by `C` is signed.~~
   ~~The value might be `""` for intermediate values and literal constants, such as `T.length.name`.~~
* `C.typed`: whether or not the constant represented by `C` is typed.
   (Maybe, it is good to require all generic constants must be typed.)
* `C.type`: the type or default type of a constant represented by `C`.

Some of the type properties can be used as constants.

_(`var`, `func`, `import` and `gen` can also be used as contract parameters/arguments,
but doing this will bring much complexity. So this is not supported temporarily.)_

Note, the following properties were removed from this propsoal:
* `T.alias`: in fact, every generic parameter `T` is an alias,
   and the type represented by `T` is always not an alias,
   so this property is not essential.
* `T.underlying`: the property is too fundamental to be much useful and used directly.
* `T.defined`, same as the above one.
* `T.receiveonly`: makes the uses of generic arguments verbose.
* `T.sendonly`: same as the above one.

## The `assure` keyword

Each `assure` line describe a constraint, or a mini contract. (See following sections for examples).

_(Other candidates to replace the `assure` keyword: `ensure`, `require`, `must`, `assert`, etc.)_

Syntax
```
assure expression
```

* if `expression` is a boolean, it must be true to pass the checking.
* if `expression` is an integer, it must be non-zero to pass the checking.
* if `expression` is a supposed fact, the fact must exist to pass the checking.

The proposed syntax is mainly to describe the following example purpose.
It can be another better form.

Not all expessions used the following examples are valid expressions used in daily Go programming. Please read the comments to get their meanings

## Constrain examples

Simple ones:
```
// T must represent a comparable type.
assure T.comparable

// Same as the above one.
assure T.value == T.value

// N and M must be two consts (either generic parameters or not).
assure N > M

// T must represent an array type which length is not smaller than 8.
assure T.kind == [0]int.kind && T.length >= 8

// T1, T2 and T3 must represent the same type.
assure T1 == T2 == T3
```

More complex ones:
```
// Tx must represent a pointer type, Ty must represent a map type,
// and the values of the base type of Tx are assignable and comparable
// to values of the element type of Ty.
assure Tx.kind == (*int).kind
assure Ty.kind == (map[int]int).kind
assure Ty.element.value = Tx.base.value
assure Ty.element.value == Tx.base.value

// Ta must represent a slice type, Tb must represent a receiveable
// channel type, and values of the element type of Tb are convertiable
// to values of the element type of Ta.
assure Ta.kind == ([]int).kind && Tb.kind == (chan int).kind
asusre Tb.receivable
assure Ta.element(Tb.element.value)
```

More:
```
// Specify the type represetned by T must have a specified method.
assure T.methods.M func(string) int

// Specify the struct type represetned by T must have a specified field.
assure T.fields.X int

// Specify the type represetned by T must have two specified selectors.
assure T.selectors.X int
assure T.selectors.F func(string) int // F can be either a method or a field of a function type

// Some more complex ones:
assure T.methods {
		.M1 func(string) int
		.M2 func(..int) (string, error)
		Ty.methods // embed a method set (a.k.a., T implements Ty)
	}
assure T.fields {
		.X int
		.Y string
		Tx.fields // embed a field set
	}
assure T.selectors {
		.X int
		.F func(string) int
		Tz.selectors // embed a selector set
	}
```

## Multiple constraints can be group in `()`

For exxample:
```
assure (
	Ty.kind == (map[int]int).kind
	Ty.element.value = Tx.base.value
	Ta.selectors {
		.X int
		.F func(string) int
		Tb.selectors
	}
)
```

## Intermediate type declarations can show up between `assure` lines

For convenience, some intermediate types are allowed to declared between assure lines.
For example:
```
assure T.kind == (map[int]int).kind
type K, E = T.key, T.element
type I interface {
	M1() []K
	M2(E) bool
}
type S struct {
	Name string
	Keys []K
}
assure T.methods { I.methods }
assure T.fields { S.field }
```

## More thinking

### Allow `assure` lines being used in non-generic code?

Use `assure` lines as assert statements, but only limited to constant expressions, in non-geneirc code. Good?

### Built-in non-elementary properties?

In fact, the above listed properties `T.orderable` and `T.signed` are not very elementary.
There are more such non-elementary properties: `T.addable`, `T.numeric`, `T.interger`, `T.floatingpoint`, `T.complex`, etc.

Good to add these ones? Or use the following introduced built-in contracts or kinds instead? 

### Built-in contracts?

Perhaps, it is good to predeclare some built-in named contracts to make some constraints
less verbose and more readable and standardized.

For example, `isNumeric[T]`, `isFloatingPoint[T]` and  `isInteger[T]` are more readable than
`T.numeric`, `T.floatingpoint` and `T.interger`, respectively.

More examples:
```
assure isArray[T] // isArray is a built-in contract
// is more readable and standardize than
assure T.kind == [0]int.kind

assure isInteger[T] && t.signed // isInteger is a built-in contract
// is more readable and less verbose than
assure T.kind == int.kind || T.kind == int8.kind || T.kind == int16.kind || T.kind == int32.kind || T.kind == int64.kind

assure sameKind[T1, T2, T3, T4] // sameKind is a built-in contract
// is less verbose than
assure T1.kind == T2.kind == T3.kind == T4.kind

assure anyKind[T, T1, T2, T3] // anyKind is a built-in contract
// is less verbose (but also less readable?) than
assure T.kind == T1.kind || T.kind == T2.kind || T.kind == T3.kind

// (BTW, are the following two lines readable?)
assure T.kind == (T1 || T2 || T3).kind
assure sameKind[T, T1 || T2 || T3]
```

### Built-in kinds?

The expression `T.kind == [0]int.kind` might be readable enough, but the `[0]int` in it
adds some irrelevant noises and might cuase some misleading.

Is it good to view `kind`s as integer values and predeclared all the kinds like
```
const (
	Bool = 1 << iota
	String
	Int
	Uint
	...
	Pointer
	Struct
	Array
	Slice
	Map
	Function
	Channel
	Interface
	...
	Signed =        Int | Int8 | Int16 | Int32 | Int64
	Unsgined =      Uint | Uint8 | Uint16 | Uint32 | Uint64 | Uintptr
	Integer =       Signed | Unsgined
	FloatingPoint = Float32 | Float64
	Complex =       Complex64 | Complex128
	Numeric =       Integer | FloatingPoint | Complex
	Orderable =     Integer | FloatingPoint | String
	Addable =       String | Numeric
	Basic =         Bool | Addable
	Ptr =           Pointer | UnsafePointer
	Container =     Array | Slcie | Map
	Any =           Basic | Ptr | Struct | Container | Function | Channel | Interface
)
```
?

Then the expression `T.kind == [0]int.kind` may be re-written as `T.kind == Array`,
which is cleaner and more readable.

More examples:
```
assure T.kind == int.kind || T.kind == int8.kind || T.kind == int16.kind || T.kind == int32.kind || T.kind == int64.kind
// may be re-written as
assure T.kind & Signed

assure T.kind == T1.kind || T.kind == T2.kind || T.kind == T3.kind
// may be re-written as
assure T.kind & (T1.kind | T2.kind | T3.kind)
```

## Tailored for the contract draft v2

In the above examples, the contract syntax used is defined in [this generic propsoal](README.md).
For [the offical contract draft 2](https://go.googlesource.com/proposal/+/master/design/go2draft-contracts.md),
some modifications are needed:
* the `[]` which encloses generic arguments of built-in contracts should be replaced with `()`.
* the `assure` lines can be put in `contract` declartion bodies, but it would be better to
  replace the `assure` keyword with another better way which is more suitable for the draft.
  For example, an `assure expression` line can be written as `{ expression }`.





