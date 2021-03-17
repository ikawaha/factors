package factors

import (
	"fmt"
)

type Factor struct {
	Exact    Set
	Prefix   Set
	Suffix   Set
	Fragment Set
}

func NewFactor() Factor {
	return Factor{
		Exact:    Set{},
		Prefix:   Set{},
		Suffix:   Set{},
		Fragment: Set{},
	}
}

func NewFactorLiteral(literal string) Factor {
	return Factor{
		Exact:    NewSet(literal),
		Prefix:   NewSet(literal),
		Suffix:   NewSet(literal),
		Fragment: NewSet(literal),
	}
}

func NewFactorInfinite() Factor {
	return Factor{
		Exact:    Set{infinite: true},
		Prefix:   Set{infinite: true},
		Suffix:   Set{infinite: true},
		Fragment: Set{infinite: true},
	}
}

func NewFactorAnyChar() Factor {
	ret := NewFactorLiteral("")
	ret.Exact.Diverge() // →∞
	return ret
}

func (f *Factor) Add(literal string) {
	f.Exact.Add(literal)
	f.Prefix.Add(literal)
	f.Suffix.Add(literal)
	f.Fragment.Add(literal)
}

func (f Factor) Infinite() bool {
	return f.Exact.infinite && f.Prefix.infinite && f.Suffix.infinite && f.Fragment.infinite
}

func (f Factor) String() string {
	return fmt.Sprintf("<exact:%s, prefix:%s, suffix:%s, fragment:%s>", f.Exact, f.Prefix, f.Suffix, f.Fragment)
}

// Alternate represents `a|b`
func Alternate(a, b Factor) Factor {
	//fmt.Printf("alt: %+v, %+v\n", a, b)
	//defer func() { fmt.Printf("  ->%+v\n", a) }()
	var ret Factor
	ret.Exact = UnionSet(a.Exact, b.Exact)
	ret.Prefix = UnionSet(a.Prefix, b.Prefix)
	ret.Suffix = UnionSet(a.Suffix, b.Suffix)
	ret.Fragment = UnionSet(a.Fragment, b.Fragment)
	return ret
}

// Concatenate represents `a・b`
func Concatenate(a, b Factor) Factor {
	//fmt.Printf("cat: %+v, %+v\n", a, b)
	//defer func() { fmt.Printf("  ->%+v\n", ret) }()
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
