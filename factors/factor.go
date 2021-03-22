package factors

import (
	"fmt"
)

// Factor represents a tuple of necessary factors for a regexp.
type Factor struct {
	Exact    Set
	Prefix   Set
	Suffix   Set
	Fragment Set
}

// NewFactor creates a factor tuple.
func NewFactor() Factor {
	return Factor{
		Exact:    Set{},
		Prefix:   Set{},
		Suffix:   Set{},
		Fragment: Set{},
	}
}

// NewFactorLiteral creates a factor tuple initialized with a given literal.
func NewFactorLiteral(literal string) Factor {
	return Factor{
		Exact:    NewSet(literal),
		Prefix:   NewSet(literal),
		Suffix:   NewSet(literal),
		Fragment: NewSet(literal),
	}
}

// NewFactorInfinite creates a factor tuple each factor set is infinite.
func NewFactorInfinite() Factor {
	return Factor{
		Exact:    Set{infinite: true},
		Prefix:   Set{infinite: true},
		Suffix:   Set{infinite: true},
		Fragment: Set{infinite: true},
	}
}

// NewFactorAnyChar creates a factor tuple initialized with regexp "any char".
func NewFactorAnyChar() Factor {
	ret := NewFactorLiteral("")
	ret.Exact.SetInfinite() // →∞
	return ret
}

// Add adds a literal to each factor set.
func (f *Factor) Add(literal string) {
	f.Exact.Add(literal)
	f.Prefix.Add(literal)
	f.Suffix.Add(literal)
	f.Fragment.Add(literal)
}

// Infinite returns true if there is a infinite set in the tuple.
func (f Factor) Infinite() bool {
	return f.Exact.infinite && f.Prefix.infinite && f.Suffix.infinite && f.Fragment.infinite
}

// String returns string representation of a tuple.
func (f Factor) String() string {
	return fmt.Sprintf("<exact:%s, prefix:%s, suffix:%s, fragment:%s>", f.Exact, f.Prefix, f.Suffix, f.Fragment)
}

// Alternate represents `a|b`
func Alternate(a, b Factor) Factor {
	var ret Factor
	ret.Exact = UnionSet(a.Exact, b.Exact)
	ret.Prefix = UnionSet(a.Prefix, b.Prefix)
	ret.Suffix = UnionSet(a.Suffix, b.Suffix)
	ret.Fragment = UnionSet(a.Fragment, b.Fragment)
	return ret
}

// Concatenate represents `a・b`
func Concatenate(a, b Factor) Factor {
	var ret Factor
	ret.Exact = CrossSet(a.Exact, b.Exact)

	ep := CrossSet(a.Exact, b.Prefix)
	ep.DropRedundantPrefix()
	ret.Prefix = BestSet(a.Prefix, ep)

	se := CrossSet(a.Suffix, b.Exact)
	se.DropRedundantSuffix()
	ret.Suffix = BestSet(b.Suffix, se)

	sp := CrossSet(a.Suffix, b.Prefix)
	sp.DropRedundantFragment()
	ret.Fragment = BestSet(a.Fragment, b.Fragment, sp)
	return ret
}
