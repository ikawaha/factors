package factor

import (
	"io"
	"regexp/syntax"

	"github.com/ikawaha/regexp/internal/factor"
)

type Factor struct {
	Exact, Prefix, Suffix, Fact []string
}

func Analyze(re *syntax.Regexp) Factor {
	tuple := factor.Analyze(re)
	return Factor{
		Exact:  tuple.Exact.Items(),
		Prefix: tuple.Pref.Items(),
		Suffix: tuple.Suff.Items(),
		Fact:   tuple.Frag.Items(),
	}
}

func DebugParse(w io.Writer, re *syntax.Regexp) Factor {
	root := factor.Parse(re)
	root.Dot(w)
	tuple := root.Factor
	return Factor{
		Exact:  tuple.Exact.Items(),
		Prefix: tuple.Pref.Items(),
		Suffix: tuple.Suff.Items(),
		Fact:   tuple.Frag.Items(),
	}
}
