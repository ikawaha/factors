package factors

import (
	"reflect"
	"regexp/syntax"
	"testing"
)

func syntaxRegexp(t *testing.T, re string) *syntax.Regexp {
	t.Helper()
	ret, err := syntax.Parse(re, syntax.Perl)
	if err != nil {
		t.Fatalf("syntex parse error: %v", err)
	}
	return ret
}

func Test_analyze(t *testing.T) {
	type args struct {
		re   *syntax.Regexp
		tree bool
	}
	tests := []struct {
		name string
		args args
		want Factor
	}{
		{
			name: "literal",
			args: args{
				re:   syntaxRegexp(t, `a`),
				tree: false,
			},
			want: NewFactorLiteral("a"),
		},
		{
			name: "a|b",
			args: args{
				re:   syntaxRegexp(t, `a|b`),
				tree: false,
			},
			want: Factor{
				Exact: Set{
					items: stringSet{
						"a": {},
						"b": {},
					},
					minimumLen: 1,
				},
				Prefix: Set{
					items: stringSet{
						"a": {},
						"b": {},
					},
					minimumLen: 1,
				},
				Suffix: Set{
					items: stringSet{
						"a": {},
						"b": {},
					},
					minimumLen: 1,
				},
				Fragment: Set{
					items: stringSet{
						"a": {},
						"b": {},
					},
					minimumLen: 1,
				},
			},
		},
		{
			name: "ab",
			args: args{
				re:   syntaxRegexp(t, `ab`),
				tree: false,
			},
			want: Factor{
				Exact: Set{
					items: stringSet{
						"ab": {},
					},
					minimumLen: 2,
				},
				Prefix: Set{
					items: stringSet{
						"ab": {},
					},
					minimumLen: 2,
				},
				Suffix: Set{
					items: stringSet{
						"ab": {},
					},
					minimumLen: 2,
				},
				Fragment: Set{
					items: stringSet{
						"ab": {},
					},
					minimumLen: 2,
				},
			},
		},
		{
			name: "a+",
			args: args{
				re:   syntaxRegexp(t, `a+`),
				tree: false,
			},
			want: Factor{
				Exact: Set{
					infinite: true,
				},
				Prefix: Set{
					items: stringSet{
						"a": {},
					},
					minimumLen: 1,
				},
				Suffix: Set{
					items: stringSet{
						"a": {},
					},
					minimumLen: 1,
				},
				Fragment: Set{
					items: stringSet{
						"a": {},
					},
					minimumLen: 1,
				},
			},
		},
		{
			name: "a.",
			args: args{
				re:   syntaxRegexp(t, `a.`),
				tree: false,
			},
			want: Factor{
				Exact: Set{
					infinite: true,
				},
				Prefix: Set{
					items: stringSet{
						"a": {},
					},
					minimumLen: 1,
				},
				Suffix: Set{
					infinite: true,
				},
				Fragment: Set{
					items: stringSet{
						"a": {},
					},
					minimumLen: 1,
				},
			},
		},
		{
			name: "X[abc]Y",
			args: args{
				re:   syntaxRegexp(t, `X[abc]Y`),
				tree: false,
			},
			want: Factor{
				Exact:    NewSet("XaY", "XbY", "XcY"),
				Prefix:   NewSet("XaY", "XbY", "XcY"),
				Suffix:   NewSet("XaY", "XbY", "XcY"),
				Fragment: NewSet("XaY", "XbY", "XcY"),
			},
		},
		{
			name: "(AG|GA)ATA((TT)*)",
			args: args{
				re:   syntaxRegexp(t, `(AG|GA)ATA((TT)*)`),
				tree: false,
			},
			want: Factor{
				Exact: Set{
					infinite: true,
				},
				Prefix: Set{
					items: stringSet{
						"AGATA": {},
						"GAATA": {},
					},
					minimumLen: 5,
				},
				Suffix: Set{
					infinite: true,
				},
				Fragment: Set{
					items: stringSet{
						"AGATA": {},
						"GAATA": {},
					},
					minimumLen: 5,
				},
			},
		},
		{
			name: "((GA|AAA)*)(TA|AG)",
			args: args{
				re:   syntaxRegexp(t, `((GA|AAA)*)(TA|AG)`),
				tree: false,
			},
			want: Factor{
				Exact: Set{
					infinite: true,
				},
				Prefix: Set{
					infinite: true,
				},
				Suffix: Set{
					items: stringSet{
						"TA": {},
						"AG": {},
					},
					minimumLen: 2,
				},
				Fragment: Set{
					items: stringSet{
						"TA": {},
						"AG": {},
					},
					minimumLen: 2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := analyze(tt.args.re, tt.args.tree); !reflect.DeepEqual(got.Factor, tt.want) {
				t.Errorf("analyze() = %v, want %v", got, tt.want)
			}
		})
	}
}
