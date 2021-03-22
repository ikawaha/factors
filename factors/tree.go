package factors

import (
	"bytes"
	"fmt"
	"io"
	"regexp/syntax"
)

const maxPrintLen = 30

// Node represents a node of a parse tree.
type Node struct {
	Factor   Factor
	Regexp   *syntax.Regexp
	Child    []*Node
	Internal bool
}

// String returns string representation of a node.
func (n Node) String() string {
	return fmt.Sprintf("op:%v, re:%v, factor:%v", Op(n.Regexp.Op), n.Regexp, n.Factor)
}

// Edge represents a edge.
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

// Dot writes a node in dot format.
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
		p := escape(abbr(ni.Factor.Prefix.String()))
		s := escape(abbr(ni.Factor.Suffix.String()))
		f := escape(abbr(ni.Factor.Fragment.String()))
		fmt.Fprintf(w, "\t\"%p\" [label=\"{ %s |{ %s | %s | %s | %s }}\"];\n", ni, l, e, p, s, f)
	}
	for _, e := range edges {
		fmt.Fprintf(w, "\t\"%p\" -- \"%p\"\n", e.from, e.to)
	}
	fmt.Fprintln(w, dotFooter)
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
