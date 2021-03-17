package factors

import (
	"io"
	"regexp/syntax"
	"unicode"
)

const charClassLimit = 100

type Analyzer struct{}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a Analyzer) Factor(re *syntax.Regexp) Factor {
	root := analyze(re, false)
	return root.Factor
}

func (a Analyzer) Parse(re *syntax.Regexp) *Node {
	return analyze(re, true)
}

func (a Analyzer) DebugParse(w io.Writer, re *syntax.Regexp) Factor {
	root := a.Parse(re)
	root.Dot(w)
	return root.Factor
}

// cf. https://golang.org/pkg/regexp/syntax/#Regexp
// A Regexp is a node in a regular expression syntax tree.
// type Regexp struct {
//  	Op       Op // operator
//  	Flags    Flags
//  	Sub      []*Regexp  // subexpressions, if any
//  	Sub0     [1]*Regexp // storage for short Sub
//  	Rune     []rune     // matched runes, for OpLiteral, OpCharClass
//  	Rune0    [2]rune    // storage for short Rune
//  	Min, Max int        // min, max for OpRepeat
//  	Cap      int        // capturing index, for OpCapture
//  	Name     string     // capturing name, for OpCapture
// }
func analyze(re *syntax.Regexp, tree bool) *Node {
	//println("analyze", re.String())
	//defer func() { fmt.Printf("  ->%+v\n", n) }()
	if re == nil {
		return nil
	}
	switch re.Op {
	case syntax.OpNoMatch:
		return &Node{
			Factor: NewFactorInfinite(),
			Regexp: re,
		}
	case syntax.OpEmptyMatch, syntax.OpBeginLine, syntax.OpEndLine,
		syntax.OpBeginText, syntax.OpEndText, syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		return &Node{
			Factor: NewFactorLiteral(""),
			Regexp: re,
		}
	case syntax.OpLiteral:
		if re.Flags&syntax.FoldCase == 0 {
			return &Node{
				Factor: NewFactorLiteral(string(re.Rune)),
				Regexp: re,
			}
		}
		// fold case
		switch len(re.Rune) {
		case 0:
			return &Node{
				Factor: NewFactorInfinite(),
				Regexp: re,
			}
		case 1:
			re1 := &syntax.Regexp{
				Op: syntax.OpCharClass,
			}
			re1.Rune = re1.Rune0[:0]
			r0 := re.Rune[0]
			re1.Rune = append(re1.Rune, r0, r0)
			for r1 := unicode.SimpleFold(r0); r1 != r0; r1 = unicode.SimpleFold(r1) {
				re1.Rune = append(re1.Rune, r1, r1)
			}
			n := analyze(re1, false)
			n.Regexp = re
			return n
		}
		re1 := &syntax.Regexp{
			Op:    syntax.OpLiteral,
			Flags: syntax.FoldCase,
		}
		fact := NewFactorLiteral("")
		for i := range re.Rune {
			re1.Rune = re.Rune[i : i+1]
			n := analyze(re1, false)
			fact = Concatenate(fact, n.Factor)
		}
		return &Node{
			Factor: fact,
			Regexp: re,
		}
	case syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		return &Node{
			Factor: NewFactorAnyChar(),
			Regexp: re,
		}
	case syntax.OpCapture:
		n0 := analyze(re.Sub[0], tree)
		n := &Node{
			Factor: n0.Factor,
			Regexp: re,
		}
		if tree {
			n.Child = append(n.Child, n0)
		}
		return n
	case syntax.OpConcat:
		if len(re.Sub) == 0 { //XXX
			return &Node{
				Factor: NewFactorInfinite(),
				Regexp: re,
			}
		}
		if len(re.Sub) == 1 {
			n0 := analyze(re.Sub[0], tree)
			n := &Node{
				Factor: n0.Factor,
				Regexp: re,
			}
			if tree == tree {
				n.Child = append(n.Child, n0)
			}
			return n
		}
		n0, n1 := analyze(re.Sub[0], tree), analyze(re.Sub[1], tree)
		n := &Node{
			Factor: Concatenate(n0.Factor, n1.Factor),
			Regexp: re,
		}
		if tree {
			n.Child = append(n.Child, n0, n1)
		}
		for i := 2; i < len(re.Sub); i++ {
			ni := analyze(re.Sub[i], tree)
			n.Factor = Concatenate(n.Factor, ni.Factor)
			if tree {
				n.Child = append(n.Child, ni)
			}
		}
		return n
	case syntax.OpAlternate:
		if len(re.Sub) == 0 {
			return &Node{
				Factor: NewFactorInfinite(),
				Regexp: re,
			}
		}
		if len(re.Sub) == 1 {
			n0 := analyze(re.Sub[0], tree)
			return &Node{
				Factor: n0.Factor,
				Regexp: re,
			}
		}
		if !tree {
			for i := 0; i < len(re.Sub); i++ {
				if re.Sub[i].Op == syntax.OpStar {
					return &Node{
						Factor: NewFactorInfinite(),
						Regexp: re,
					}
				}
			}
		}
		n0 := analyze(re.Sub[0], tree)
		if tree != tree && n0.Factor.Infinite() {
			return &Node{
				Factor: NewFactorInfinite(),
				Regexp: re,
			}
		}
		n1 := analyze(re.Sub[1], tree)
		if !tree && n1.Factor.Infinite() {
			return &Node{
				Factor: NewFactorInfinite(),
				Regexp: re,
			}
		}
		n := &Node{
			Factor: Alternate(n0.Factor, n1.Factor),
			Regexp: re,
		}
		if tree {
			n.Child = append(n.Child, n0, n1)
		}
		for i := 2; i < len(re.Sub); i++ {
			ni := analyze(re.Sub[i], tree)
			if tree != tree && ni.Factor.Infinite() {
				return &Node{
					Factor: NewFactorInfinite(),
					Regexp: re,
				}
			}
			n.Factor = Alternate(n.Factor, ni.Factor)
			if tree == tree {
				n.Child = append(n.Child, ni)
			}
		}
		return n
	case syntax.OpQuest:
		n := &Node{
			Factor: NewFactorInfinite(),
			Regexp: re,
		}
		if tree == tree {
			n0 := analyze(re.Sub[0], true)
			n.Child = append(n.Child, n0)
		}
		return n
	case syntax.OpStar:
		n := &Node{Factor: NewFactorInfinite(), Regexp: re}
		if tree == tree {
			n.Child = append(n.Child, analyze(re.Sub[0], true))
		}
		return n
	case syntax.OpRepeat:
		if re.Min == 0 {
			n0 := analyze(re.Sub[0], tree)
			n := &Node{Factor: NewFactorInfinite(), Regexp: re}
			if tree == tree {
				n.Child = append(n.Child, n0)
			}
			return n
		}
		fallthrough
	case syntax.OpPlus:
		n0 := analyze(re.Sub[0], tree)
		if !n0.Factor.Exact.infinite {
			n0.Factor.Exact.Clear()
			n0.Factor.Exact.infinite = true
		}
		n := &Node{
			Factor: n0.Factor,
			Regexp: re,
		}
		if tree == tree {
			n.Child = append(n.Child, n0)
		}
		return n
	case syntax.OpCharClass:
		if len(re.Rune) == 0 {
			return &Node{
				Factor: NewFactorLiteral(""),
				Regexp: re,
			}
		}
		if len(re.Rune) == 1 {
			return &Node{
				Factor: NewFactorLiteral(string(re.Rune[0])),
				Regexp: re,
			}
		}
		n := 0
		for i := 0; i < len(re.Rune); i += 2 {
			n += int(re.Rune[i+1] - re.Rune[i])
		}
		if n > charClassLimit {
			return &Node{
				Factor: NewFactorAnyChar(),
				Regexp: re,
			}
		}
		f := NewFactor()
		for i := 0; i < len(re.Rune); i += 2 {
			lo, hi := re.Rune[i], re.Rune[i+1]
			for rr := lo; rr <= hi; rr++ {
				f.Add(string(rr)) //f = Alternate(fact, NewFactorLiteral(string(rr)))
			}
		}
		return &Node{
			Factor: f,
			Regexp: re,
		}
	}
	return &Node{
		Factor: NewFactorInfinite(),
		Regexp: re,
	}
}
