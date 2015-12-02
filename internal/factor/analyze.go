package factor

import (
	"bytes"
	"fmt"
	"io"
	"regexp/syntax"
	"strings"
	"unicode"
)

const (
	TolerateCharClassCount = 100
)

var (
	RegexpOpCodeTable = []string{
		"", "NoMatch", "EmptyMatch", "Literal", "CC",
		"AnyCharNL", "AnyChar", "^", "$", "\\A", "$",
		"\\b", "\\B", "()", "*", "+", "?", "RT", "・", "|",
	}
)

type Op syntax.Op

func (o Op) String() string {
	op := int(o)
	if op >= 0 && op < len(RegexpOpCodeTable) {
		return RegexpOpCodeTable[op]
	}
	return "Unknown"
}

type Node struct {
	Regexp   *syntax.Regexp
	Child    []*Node
	Factor   FactorTuple
	Internal bool
}

//for debug
func (n Node) String() string {
	return fmt.Sprintf("op:%v, re:%v, factor:%v", Op(n.Regexp.Op), n.Regexp, n.Factor)
}

//for debug
func PrintTree(n *Node, depth int) {
	if n == nil {
		return
	}
	fmt.Printf("%v%+v\n", strings.Repeat("-", depth), n)
	for _, c := range n.Child {
		PrintTree(c, depth+1)
	}
}

type Edge struct {
	from *Node
	to   *Node
}

func walk(n *Node) ([]*Node, []Edge) {
	if n == nil {
		return nil, nil
	}

	ns := []*Node{n}
	es := []Edge{}
	switch n.Regexp.Op {
	case syntax.OpBeginLine, syntax.OpBeginText,
		syntax.OpEndLine, syntax.OpEndText:
		n.Internal = true
	}

	if len(n.Child) != 0 {
		tmp := &Node{
			Regexp:   n.Regexp,
			Child:    n.Child,
			Internal: true,
		}
		ns = append(ns, tmp)
		es = append(es, Edge{from: n, to: tmp})
		n = tmp
	}
	for _, c := range n.Child {
		if c.Regexp.Op == syntax.OpCharClass {
			tmp := &Node{
				Regexp:   c.Regexp,
				Factor:   c.Factor,
				Internal: true,
			}
			ns = append(ns, tmp)
			es = append(es, Edge{from: n, to: tmp})
			es = append(es, Edge{from: tmp, to: c})
		} else {
			es = append(es, Edge{from: n, to: c})
		}
		cn, ce := walk(c)
		ns = append(ns, cn...)
		es = append(es, ce...)
	}
	return ns, es
}

func escape(s string) string {
	var b bytes.Buffer
	for _, r := range s {
		switch r {
		case '\\', '|', '"', '{', '}', '[', ']':
			b.WriteRune('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

const maxPrintLen = 30

func abbr(s string) string {
	if len(s) <= maxPrintLen {
		return s
	}
	for i := range s {
		if i > maxPrintLen {
			s = s[:i] + "...}"
			break
		}
	}
	return s
}

func (n *Node) Dot(w io.Writer) {
	const (
		dotHeader = "graph regexptree {\tdpi=48\tgraph [style=filed];\tnode [shape=record];"
		dotFooter = "}"
	)
	nodes, edges := walk(n)
	fmt.Fprintln(w, dotHeader)
	for _, ni := range nodes {
		if ni.Internal {
			op := escape(Op(ni.Regexp.Op).String())
			fmt.Fprintf(w, "\t\"%p\" [shape=doublecircle, label=\"%s\"];\n", ni, op)
			continue
		}
		l := escape(ni.Regexp.String())
		e := escape(abbr(ni.Factor.Exact.String()))
		p := escape(abbr(ni.Factor.Pref.String()))
		s := escape(abbr(ni.Factor.Suff.String()))
		f := escape(abbr(ni.Factor.Frag.String()))
		fmt.Fprintf(w, "\t\"%p\" [label=\"{ %s |{ %s | %s | %s | %s }}\"];\n", ni, l, e, p, s, f)
	}
	for _, e := range edges {
		fmt.Fprintf(w, "\t\"%p\" -- \"%p\"\n", e.from, e.to)
	}
	fmt.Fprintln(w, dotFooter)
}

type analyzeMode uint8

const (
	simple analyzeMode = iota + 1
	tree
)

func analyze(re *syntax.Regexp, mode analyzeMode) (n *Node) {
	//println("analyze", re.String())
	//defer func() { fmt.Printf("  ->%+v\n", n) }()

	if re == nil {
		return nil
	}
	switch re.Op {
	case syntax.OpNoMatch:
		return &Node{Factor: NewFactorTupleUndef(), Regexp: re}

	case syntax.OpEmptyMatch,
		syntax.OpBeginLine, syntax.OpEndLine,
		syntax.OpBeginText, syntax.OpEndText,
		syntax.OpWordBoundary, syntax.OpNoWordBoundary:
		return &Node{Factor: NewFactorTupleEmptyString(), Regexp: re}

	case syntax.OpLiteral:
		if re.Flags&syntax.FoldCase != 0 {
			switch len(re.Rune) {
			case 0:
				return &Node{Factor: NewFactorTupleUndef(), Regexp: re}
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
				n = analyze(re1, simple)
				n.Regexp = re
				return
			}
			re1 := &syntax.Regexp{
				Op:    syntax.OpLiteral,
				Flags: syntax.FoldCase,
			}
			fact := NewFactorTupleLiteral("")
			for i := range re.Rune {
				re1.Rune = re.Rune[i : i+1]
				fact = Concatinate(fact, Analyze(re1))
			}
			return &Node{Factor: fact, Regexp: re}
		}
		return &Node{Factor: NewFactorTupleLiteral(string(re.Rune)), Regexp: re}

	case syntax.OpAnyCharNotNL, syntax.OpAnyChar:
		return &Node{Factor: NewFactorTupleAnyChar(), Regexp: re}

	case syntax.OpCapture:
		n0 := analyze(re.Sub[0], mode)
		n = &Node{Factor: n0.Factor, Regexp: re}
		if mode == tree {
			n.Child = append(n.Child, n0)
		}
		return

	case syntax.OpConcat:
		if len(re.Sub) == 0 {
			return &Node{Factor: NewFactorTupleUndef(), Regexp: re} //XXX
		}
		if len(re.Sub) == 1 {
			n0 := analyze(re.Sub[0], mode)
			n = &Node{Factor: n0.Factor, Regexp: re}
			if mode == tree {
				n.Child = append(n.Child, n0)
			}
			return
		}
		n0, n1 := analyze(re.Sub[0], mode), analyze(re.Sub[1], mode)
		n = &Node{Factor: Concatinate(n0.Factor, n1.Factor), Regexp: re}
		if mode == tree {
			n.Child = append(n.Child, n0, n1)
		}
		for i := 2; i < len(re.Sub); i++ {
			ni := analyze(re.Sub[i], mode)
			n.Factor = Concatinate(n.Factor, ni.Factor)
			if mode == tree {
				n.Child = append(n.Child, ni)
			}
		}
		return

	case syntax.OpAlternate:
		if len(re.Sub) == 0 {
			return &Node{Factor: NewFactorTupleUndef(), Regexp: re}
		}
		if len(re.Sub) == 1 {
			n0 := analyze(re.Sub[0], mode)
			return &Node{Factor: n0.Factor, Regexp: re}
		}
		if mode != tree {
			for i := 0; i < len(re.Sub); i++ {
				if re.Sub[i].Op == syntax.OpStar {
					return &Node{Factor: NewFactorTupleUndef(), Regexp: re}
				}
			}
		}
		n0 := analyze(re.Sub[0], mode)
		if mode != tree && n0.Factor.Undef() {
			return &Node{Factor: NewFactorTupleUndef(), Regexp: re}
		}
		n1 := analyze(re.Sub[1], mode)
		if mode != tree && n1.Factor.Undef() {
			return &Node{Factor: NewFactorTupleUndef(), Regexp: re}
		}
		f := n0.Factor
		n0.Factor.Exact = n0.Factor.Exact.Clone()
		n0.Factor.Pref = n0.Factor.Pref.Clone()
		n0.Factor.Suff = n0.Factor.Suff.Clone()
		n0.Factor.Frag = n0.Factor.Frag.Clone()
		n = &Node{Factor: Alternate(f, n1.Factor), Regexp: re}
		if mode == tree {
			n.Child = append(n.Child, n0, n1)
		}
		for i := 2; i < len(re.Sub); i++ {
			ni := analyze(re.Sub[i], mode)
			if mode != tree && ni.Factor.Undef() {
				return &Node{Factor: NewFactorTupleUndef(), Regexp: re}
			}
			n.Factor = Alternate(n.Factor, ni.Factor)
			if mode == tree {
				n.Child = append(n.Child, ni)
			}
		}
		return

	case syntax.OpQuest:
		n = &Node{Factor: NewFactorTupleUndef(), Regexp: re}
		if mode == tree {
			n0 := Parse(re.Sub[0])
			n.Child = append(n.Child, n0)
		}
		return

	case syntax.OpStar:
		n = &Node{Factor: NewFactorTupleUndef(), Regexp: re}
		if mode == tree {
			n.Child = append(n.Child, Parse(re.Sub[0]))
		}
		return

	case syntax.OpRepeat:
		if re.Min == 0 {
			n0 := analyze(re.Sub[0], mode)
			n = &Node{Factor: NewFactorTupleUndef(), Regexp: re}
			if mode == tree {
				n.Child = append(n.Child, n0)
			}
			return
		}
		fallthrough
	case syntax.OpPlus:
		n0 := analyze(re.Sub[0], mode)
		if !n0.Factor.Exact.Undef {
			n0.Factor.Pref = n0.Factor.Exact
			n0.Factor.Suff = n0.Factor.Exact.Clone()
			n0.Factor.Frag = n0.Factor.Exact.Clone()
			n0.Factor.Exact.Clear()
			n0.Factor.Exact.Undef = true
		}
		n = &Node{Factor: n0.Factor, Regexp: re}
		if mode == tree {
			n.Child = append(n.Child, n0)
		}
		return n

	case syntax.OpCharClass:
		if len(re.Rune) == 0 {
			return &Node{Factor: NewFactorTupleEmptyString(), Regexp: re}
		}
		if len(re.Rune) == 1 {
			return &Node{Factor: NewFactorTupleLiteral(string(re.Rune[0])), Regexp: re}
		}
		n := 0
		for i := 0; i < len(re.Rune); i += 2 {
			n += int(re.Rune[i+1] - re.Rune[i])
		}
		if n > TolerateCharClassCount {
			return &Node{Factor: NewFactorTupleAnyChar(), Regexp: re}
		}
		f := NewFactorTuple()
		for i := 0; i < len(re.Rune); i += 2 {
			lo, hi := re.Rune[i], re.Rune[i+1]
			for rr := lo; rr <= hi; rr++ {
				f.Exact.Add(string(rr)) //f = Alternate(fact, NewFactorTupleLiteral(string(rr)))
			}
		}
		return &Node{Factor: f, Regexp: re}
	}
	return &Node{Factor: NewFactorTupleUndef(), Regexp: re}
}

func Analyze(re *syntax.Regexp) FactorTuple {
	n := analyze(re, simple)
	return n.Factor
}

func Parse(re *syntax.Regexp) (n *Node) {
	return analyze(re, tree)
}
