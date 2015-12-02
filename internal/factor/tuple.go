package factor

import (
	"fmt"
)

type FactorTuple struct {
	Exact, Pref, Suff, Frag Factor
}

func (f FactorTuple) Undef() bool {
	return f.Exact.Undef && f.Pref.Undef &&
		f.Suff.Undef && f.Frag.Undef
}

func (f FactorTuple) String() string {
	if f.Undef() {
		return "--"
	} else if !f.Exact.Undef {
		return fmt.Sprintf("%v", f.Exact)
	}
	return fmt.Sprintf("pref:%v, suff:%v, frag:%v", f.Pref, f.Suff, f.Frag)
}

func NewFactorTuple() FactorTuple {
	var f FactorTuple
	f.Exact = NewSet()
	f.Pref = f.Exact
	f.Suff = f.Exact
	f.Frag = f.Exact
	return f
}

func NewFactorTupleLiteral(str string) FactorTuple {
	var f FactorTuple
	f.Exact.Add(str)
	f.Pref = f.Exact
	f.Suff = f.Exact
	f.Frag = f.Exact
	return f
}

func NewFactorTupleUndef() FactorTuple {
	return FactorTuple{
		Exact: Factor{Undef: true},
		Pref:  Factor{Undef: true},
		Suff:  Factor{Undef: true},
		Frag:  Factor{Undef: true},
	}
}

func NewFactorTupleAnyChar() FactorTuple {
	f := NewFactorTupleUndef()
	f.Frag.Undef = false
	f.Frag.Add("")
	return f
}

func NewFactorTupleEmptyString() FactorTuple {
	var f FactorTuple
	f.Exact.Add("")
	f.Pref = f.Exact
	f.Suff = f.Exact
	f.Frag = f.Exact
	return f
}

func Alternate(x, y FactorTuple) FactorTuple {
	//fmt.Printf("alt: %+v, %+v\n", x, y)
	//defer func() { fmt.Printf("  ->%+v\n", x) }()

	undef := x.Exact.Undef
	x.Exact.Union(y.Exact)
	if !x.Exact.Undef {
		return x
	}
	if !undef {
		x.Suff = x.Pref.Clone()
		x.Frag = x.Pref.Clone()
	}
	x.Pref.Union(y.Pref)
	x.Suff.Union(y.Suff)
	x.Frag.Union(y.Frag)
	return x
}

func Concatinate(x, y FactorTuple) (xy FactorTuple) {
	//fmt.Printf("cat: %+v, %+v\n", x, y)
	//defer func() { fmt.Printf("  ->%+v\n", xy) }()

	xy.Exact = CrossSet(x.Exact, y.Exact)

	if !xy.Exact.Undef {
		xy.Pref = xy.Exact
		xy.Suff = xy.Exact
		xy.Frag = xy.Exact
		return
	}

	if !x.Exact.Undef && !y.Pref.Undef {
		xy.Pref = CrossSet(x.Exact, y.Pref)
	} else {
		xy.Pref = x.Pref
	}

	if !x.Suff.Undef && !y.Exact.Undef {
		xy.Suff = CrossSet(x.Suff, y.Exact)
	} else {
		xy.Suff = y.Suff
	}
	cs := CrossSet(x.Suff, y.Pref)
	if b := BestSet(&cs, &x.Frag, &y.Frag); b != nil {
		xy.Frag = *b
	} else {
		xy.Frag = Factor{Undef: true}
	}
	return
}
