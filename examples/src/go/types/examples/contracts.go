package p

// A pure contracts can viewed as a `gen` without any output.

// A contract requires a type has a specified method.	
gen C[T type] {
	var x T
	var _ interface {
	    m()
	} = x
}
	
// Contracts may be empty.
gen Empty[] {}
	
	
// Like functions, contracts (gens) may not be grouped.
gen ( // ERROR 
	C1[T] {} // ERROR 
	C2[T] {} // ERROR 
	C3[T] {} // ERROR 
)
	
	
// A contract specifies methods and types for each of the
// type parameters it constrains.
gen Stringer[T type] {
    var x T
	var _ interface {
	    String() string
	} = x
}
	
gen Sequence[T type] {
	var x T
	var _ T = x[:]
}

// Contracts may constrain multiple type parameters
// in mutually recursive ways.
gen G[Node, Edge type] {
	var n Node
	var _ interface {
	    Edges() []Edge
	} = n
	
	var e Edge
	var _ interface {
	    Nodes() (from Node, to Node)
	} = e
}

gen Graph [Node, Edge type] import {
	G[Node, Edge] // apply an extra contract

	// For a gen which ouputs an import, all the exported types
	// and functions declared in the gen body will be outputted,
	// their exported names are just their declaration names.
	//
	// For this specified gen, one type and two functions will
	// be outputted together in a mini-pacakge.

    type Graph struct { /* ... */ }
    
    func New(nodes []Node) *Graph(Node, Edge) { panic("unimplemented") }
    
    func ShortestPath(from, to Node) []Edge { panic("unimplemented") }
}

