package factors

import (
	"reflect"
	"testing"
)

func TestAlternate(t *testing.T) {
	type args struct {
		a Factor
		b Factor
	}
	tests := []struct {
		name string
		args args
		want Factor
	}{
		{
			name: "θ|<{a}, {a}, {a}, {a}>",
			args: args{
				a: NewFactorInfinite(),
				b: NewFactorLiteral("a"),
			},
			want: Factor{
				Exact:    Set{infinite: true},
				Prefix:   Set{infinite: true},
				Suffix:   Set{infinite: true},
				Fragment: Set{infinite: true},
			},
		},
		{
			name: "<{a}, {a}, {a}, {a}>|<{b}, {b}, {b}, {b}>",
			args: args{
				a: NewFactorLiteral("a"),
				b: NewFactorLiteral("b"),
			},
			want: Factor{
				Exact:    NewSet("a", "b"),
				Prefix:   NewSet("a", "b"),
				Suffix:   NewSet("a", "b"),
				Fragment: NewSet("a", "b"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Alternate(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Alternate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConcatenate(t *testing.T) {
	type args struct {
		a Factor
		b Factor
	}
	tests := []struct {
		name string
		args args
		want Factor
	}{
		{
			name: "θ|<{a}, {a}, {a}, {a}>",
			args: args{
				a: NewFactorInfinite(),
				b: NewFactorLiteral("a"),
			},
			want: Factor{
				Exact:    Set{infinite: true},
				Prefix:   Set{infinite: true},
				Suffix:   Set{infinite: true},
				Fragment: Set{infinite: true},
			},
		},
		{
			name: "<{a}, {a}, {a}, {a}>・<{b}, {b}, {b}, {b}>",
			args: args{
				a: NewFactorLiteral("a"),
				b: NewFactorLiteral("b"),
			},
			want: Factor{
				Exact:    NewSet("ab"),
				Prefix:   NewSet("a"),
				Suffix:   NewSet("b"),
				Fragment: NewSet("a", "b"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Concatenate(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Concatenate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactor_String(t *testing.T) {
	type fields struct {
		Exact    Set
		Prefix   Set
		Suffix   Set
		Fragment Set
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Factor{
				Exact:    tt.fields.Exact,
				Prefix:   tt.fields.Prefix,
				Suffix:   tt.fields.Suffix,
				Fragment: tt.fields.Fragment,
			}
			if got := f.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactor_Undef(t *testing.T) {
	type fields struct {
		Exact    Set
		Prefix   Set
		Suffix   Set
		Fragment Set
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Factor{
				Exact:    tt.fields.Exact,
				Prefix:   tt.fields.Prefix,
				Suffix:   tt.fields.Suffix,
				Fragment: tt.fields.Fragment,
			}
			if got := f.Infinite(); got != tt.want {
				t.Errorf("Infinite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFactor(t *testing.T) {
	tests := []struct {
		name string
		want Factor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFactor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFactor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFactorAnyChar(t *testing.T) {
	tests := []struct {
		name string
		want Factor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFactorAnyChar(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFactorAnyChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFactorEmptyLiteral(t *testing.T) {
	tests := []struct {
		name string
		want Factor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFactorLiteral(""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFactorEmptyLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFactorLiteral(t *testing.T) {
	type args struct {
		literal string
	}
	tests := []struct {
		name string
		args args
		want Factor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFactorLiteral(tt.args.literal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFactorLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFactorUndef(t *testing.T) {
	tests := []struct {
		name string
		want Factor
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFactorInfinite(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFactorInfinite() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestBestPrefix(t *testing.T) {
//	type args struct {
//		x Set
//		y Set
//	}
//	tests := []struct {
//		name string
//		args args
//		want Set
//	}{
//		{
//			name: "infinite",
//			args: args{
//				x: Set{infinite: true},
//				y: Set{},
//			},
//			want: Set{infinite: true},
//		},
//		{
//			name: "overlap",
//			args: args{
//				x: Set{
//					items: stringSet{
//						"a":     {},
//						"b":     {}, // overlap
//						"c":     {},
//						"y":     {}, // overlap
//						"hello": {},
//					},
//					minimumLen: 1,
//				},
//				y: Set{
//					items: stringSet{
//						"b":       {}, // overlap
//						"x":       {},
//						"y":       {}, // overlap
//						"z":       {},
//						"goodbye": {},
//					},
//					minimumLen: 1,
//				},
//			},
//			want: Set{
//				items: stringSet{
//					"a":       {},
//					"b":       {},
//					"c":       {},
//					"goodbye": {},
//					"hello":   {},
//					"x":       {},
//					"y":       {},
//					"z":       {},
//				},
//				minimumLen: 1,
//			},
//		},
//		{
//			name: "prefix",
//			args: args{
//				x: Set{
//					items: stringSet{
//						"a":     {}, // prefix of "aa"
//						"aa":    {},
//						"b":     {}, // prefix of "bcd"
//						"y":     {}, // prefix of "y"
//						"hello": {},
//					},
//					minimumLen: 1,
//				},
//				y: Set{
//					items: stringSet{
//						"bcd":     {},
//						"c":       {},
//						"x":       {},
//						"y":       {},
//						"z":       {},
//						"goodbye": {},
//						"hell":    {}, // prefix of "hello"
//					},
//					minimumLen: 1,
//				},
//			},
//			want: Set{
//				items: stringSet{
//					"a":       {},
//					"b":       {},
//					"c":       {},
//					"goodbye": {},
//					"hell":    {},
//					"x":       {},
//					"y":       {},
//					"z":       {},
//				},
//				minimumLen: 1,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := BestPrefix(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("BestPrefix() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func TestBestSuffix(t *testing.T) {
//	type args struct {
//		x Set
//		y Set
//	}
//	tests := []struct {
//		name string
//		args args
//		want Set
//	}{
//		{
//			name: "infinite",
//			args: args{
//				x: Set{infinite: true},
//				y: Set{},
//			},
//			want: Set{infinite: true},
//		},
//		{
//			name: "overlap",
//			args: args{
//				x: Set{
//					items: stringSet{
//						"a":     {},
//						"b":     {}, // overlap
//						"c":     {},
//						"y":     {}, // overlap
//						"hello": {},
//					},
//					minimumLen: 1,
//				},
//				y: Set{
//					items: stringSet{
//						"b":       {}, // overlap
//						"x":       {},
//						"y":       {}, // overlap
//						"z":       {},
//						"goodbye": {},
//					},
//					minimumLen: 1,
//				},
//			},
//			want: Set{
//				items: stringSet{
//					"a":       {},
//					"b":       {},
//					"c":       {},
//					"goodbye": {},
//					"hello":   {},
//					"x":       {},
//					"y":       {},
//					"z":       {},
//				},
//				minimumLen: 1,
//			},
//		},
//		{
//			name: "suffix",
//			args: args{
//				x: Set{
//					items: stringSet{
//						"a":       {},
//						"b":       {}, // suffix of "dcb"
//						"y":       {}, // suffix of "y"
//						"hello":   {},
//						"goodbye": {},
//					},
//					minimumLen: 1,
//				},
//				y: Set{
//					items: stringSet{
//						"dcb":    {},
//						"c":      {},
//						"x":      {},
//						"y":      {},
//						"z":      {},
//						"byebye": {},
//						"bye":    {}, // suffix of "goodbye" and "byebye"
//					},
//					minimumLen: 1,
//				},
//			},
//			want: Set{
//				items: stringSet{
//					"a":     {},
//					"b":     {},
//					"c":     {},
//					"bye":   {},
//					"hello": {},
//					"x":     {},
//					"y":     {},
//					"z":     {},
//				},
//				minimumLen: 1,
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := BestSuffix(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("BestSuffix() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

/*
func TestBestFragment(t *testing.T) {
	type args struct {
		x Set
		y Set
	}
	tests := []struct {
		name string
		args args
		want Set
	}{
		{
			name: "infinite",
			args: args{
				x: Set{infinite: true},
				y: Set{},
			},
			want: Set{infinite: true},
		},
		{
			name: "same sets",
			args: args{
				x: Set{
					items: stringSet{
						"aaa": {},
						"bbb": {},
						"ccc": {},
					},
					minimumLen: 3,
				},
				y: Set{
					items: stringSet{
						"aaa": {},
						"bbb": {},
						"ccc": {},
					},
					minimumLen: 3,
				},
			},
			want: Set{
				items: stringSet{
					"aaa": {},
					"bbb": {},
					"ccc": {},
				},
				minimumLen: 3,
			},
		},
		{
			name: "{abc}, {}",
			args: args{
				x: Set{
					items: stringSet{
						"abc": {},
					},
					minimumLen: 3,
				},
				y: Set{
					items:      stringSet{},
					minimumLen: 0,
				},
			},
			want: Set{
				items:      stringSet{},
				minimumLen: 0,
			},
		},
		{
			name: "{abc}, {def}",
			args: args{
				x: Set{
					items: stringSet{
						"abc": {},
					},
					minimumLen: 3,
				},
				y: Set{
					items: stringSet{
						"def": {},
					},
					minimumLen: 3,
				},
			},
			want: Set{
				items:      stringSet{},
				minimumLen: 0,
			},
		},
		{
			name: "{abc}, {bcd}",
			args: args{
				x: Set{
					items: stringSet{
						"abc": {},
					},
					minimumLen: 3,
				},
				y: Set{
					items: stringSet{
						"bcd": {},
					},
					minimumLen: 3,
				},
			},
			want: Set{
				items: stringSet{
					"bc": {},
				},
				minimumLen: 2,
			},
		},
		{
			name: "{abc, abd}, {bcd}",
			args: args{
				x: Set{
					items: stringSet{
						"abc": {},
						"abd": {},
					},
					minimumLen: 3,
				},
				y: Set{
					items: stringSet{
						"bcd": {},
					},
					minimumLen: 3,
				},
			},
			want: Set{
				items: stringSet{
					"bc": {},
					"d":  {},
				},
				minimumLen: 1,
			},
		},
		{
			name: "{abc, cbd}, {abd}",
			args: args{
				x: Set{
					items: stringSet{
						"abc": {},
						"cbd": {},
					},
					minimumLen: 3,
				},
				y: Set{
					items: stringSet{
						"abd": {},
					},
					minimumLen: 3,
				},
			},
			want: Set{
				items: stringSet{
					"ab": {},
					"bd": {},
				},
				minimumLen: 2,
			},
		},
		{
			name: "{abc, cbd}, {abd, cbc}",
			args: args{
				x: Set{
					items: stringSet{
						"abc": {},
						"cbd": {},
					},
					minimumLen: 3,
				},
				y: Set{
					items: stringSet{
						"abd": {},
						"cbc": {},
					},
					minimumLen: 3,
				},
			},
			want: Set{
				items: stringSet{
					"ab": {},
					"bd": {},
					"bc": {},
					"cb": {},
				},
				minimumLen: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BestFragment(tt.args.x, tt.args.y); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BestFragment() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
