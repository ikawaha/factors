package factors

import (
	"reflect"
	"testing"
)

func TestCrossSet(t *testing.T) {
	type args struct {
		a Set
		b Set
	}
	tests := []struct {
		name string
		args args
		want Set
	}{
		{
			name: "cross set",
			args: args{
				a: NewSet("aa1", "bb1", "cc1"),
				b: NewSet("a2", "b2", "cccc2"),
			},
			want: NewSet(
				"aa1a2", "aa1b2", "aa1cccc2",
				"bb1a2", "bb1b2", "bb1cccc2",
				"cc1a2", "cc1b2", "cc1cccc2",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CrossSet(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CrossSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSet(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		args args
		want Set
	}{
		{
			name: "new without any item",
			args: args{
				items: nil,
			},
			want: Set{},
		},
		{
			name: "new with some items",
			args: args{
				items: []string{
					"hello",
					"goodbye",
				},
			},
			want: Set{
				infinite:   false,
				minimumLen: 5,
				items: map[string]struct{}{
					"hello":   {},
					"goodbye": {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSet(tt.args.items...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Add(t *testing.T) {
	type fields struct {
		undef      bool
		minimumLen int
		items      stringSet
	}
	type args struct {
		item string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Set
	}{
		{
			name: "add an item to empty set",
			fields: fields{
				undef:      false,
				minimumLen: 0,
				items:      nil,
			},
			args: args{
				item: "hello",
			},
			want: Set{
				infinite:   false,
				minimumLen: 5,
				items: stringSet{
					"hello": {},
				},
			},
		},
		{
			name: "add an item to a set",
			fields: fields{
				undef:      false,
				minimumLen: 7,
				items: stringSet{
					"goodbye": {},
				},
			},
			args: args{
				item: "hello",
			},
			want: Set{
				infinite:   false,
				minimumLen: 5,
				items: stringSet{
					"hello":   {},
					"goodbye": {},
				},
			},
		},
		{
			name: "add a duplicate item to a set",
			fields: fields{
				undef:      false,
				minimumLen: 5,
				items: stringSet{
					"goodbye": {},
					"hello":   {},
				},
			},
			args: args{
				item: "hello",
			},
			want: Set{
				infinite:   false,
				minimumLen: 5,
				items: stringSet{
					"hello":   {},
					"goodbye": {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Set{
				infinite:   tt.fields.undef,
				minimumLen: tt.fields.minimumLen,
				items:      tt.fields.items,
			}
			s.Add(tt.args.item)
			if !reflect.DeepEqual(s, tt.want) {
				t.Errorf("Add() = %#+v, want %#+v", s, tt.want)
			}
		})
	}
}

func TestSet_Clear(t *testing.T) {
	type fields struct {
		undef      bool
		minimumLen int
		items      stringSet
	}
	tests := []struct {
		name   string
		fields fields
		want   Set
	}{
		{
			name: "clear a set",
			fields: fields{
				undef:      true,
				minimumLen: 5,
				items: stringSet{
					"hello":   {},
					"goodbye": {},
					"aloha":   {},
				},
			},
			want: Set{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Set{
				infinite:   tt.fields.undef,
				minimumLen: tt.fields.minimumLen,
				items:      tt.fields.items,
			}
			s.Clear()
			if !reflect.DeepEqual(s, tt.want) {
				t.Errorf("Clear() = %v, want %v", s, tt.want)
			}
		})
	}
}

func TestSet_Clone(t *testing.T) {
	type fields struct {
		undef      bool
		minimumLen int
		items      stringSet
	}
	tests := []struct {
		name   string
		fields fields
		want   Set
	}{
		{
			name: "clone a set",
			fields: fields{
				undef:      false,
				minimumLen: 5,
				items: stringSet{
					"hello": {},
				},
			},
			want: Set{
				infinite:   false,
				minimumLen: 5,
				items: stringSet{
					"hello": {},
				},
			},
		},
		{
			name: "clone an infinite set",
			fields: fields{
				undef:      true,
				minimumLen: 5,
				items: stringSet{
					"hello": {},
				},
			},
			want: Set{
				infinite:   true,
				minimumLen: 5,
				items: stringSet{
					"hello": {},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Set{
				infinite:   tt.fields.undef,
				minimumLen: tt.fields.minimumLen,
				items:      tt.fields.items,
			}
			if got := s.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Items(t *testing.T) {
	type fields struct {
		undef      bool
		minimumLen int
		items      stringSet
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name:   "items of an empty set",
			fields: fields{},
			want:   nil,
		},
		{
			name: "items of a not empty set",
			fields: fields{
				undef:      false,
				minimumLen: 5,
				items: stringSet{
					"hello":   {},
					"goodbye": {},
				},
			},
			want: []string{"goodbye", "hello"},
		},
		{
			name: "items of a not empty infinite set",
			fields: fields{
				undef:      true,
				minimumLen: 5,
				items: stringSet{
					"hello":   {},
					"goodbye": {},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Set{
				infinite:   tt.fields.undef,
				minimumLen: tt.fields.minimumLen,
				items:      tt.fields.items,
			}
			if got := s.Items(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Items() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_Size(t *testing.T) {
	type fields struct {
		undef      bool
		minimumLen int
		items      stringSet
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "empty set",
			fields: fields{},
			want:   0,
		},
		{
			name: "size of a set",
			fields: fields{
				undef:      false,
				minimumLen: 0,
				items: stringSet{
					"hello":   {},
					"goodbye": {},
				},
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Set{
				infinite:   tt.fields.undef,
				minimumLen: tt.fields.minimumLen,
				items:      tt.fields.items,
			}
			if got := s.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_String(t *testing.T) {
	type fields struct {
		undef      bool
		minimumLen int
		items      stringSet
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "empty set",
			fields: fields{},
			want:   "{}",
		},
		{
			name: "infinite set",
			fields: fields{
				undef: true,
			},
			want: theta,
		},
		{
			name: "string representation of a set",
			fields: fields{
				items: stringSet{
					"hello":   {},
					"goodbye": {},
				},
			},
			want: "{goodbye, hello}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Set{
				infinite:   tt.fields.undef,
				minimumLen: tt.fields.minimumLen,
				items:      tt.fields.items,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_longestCommonSubstring(t *testing.T) {
	type args struct {
		x string
		y string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "apple, pineapple ->apple",
			args: args{
				x: "apple",
				y: "pineapple",
			},
			want: "apple",
		},
		{
			name: "pineapple, apple ->apple",
			args: args{
				x: "apple",
				y: "pineapple",
			},
			want: "apple",
		},
		{
			name: "ABRACADABRA, ECADADABRBCRDARA ->ADABR",
			args: args{
				x: "ABRACADABRA",
				y: "ECADADABRBCRDARA",
			},
			want: "ADABR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longestCommonSubstring(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("longestCommonSubstring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet_LongestCommon(t *testing.T) {
	type fields struct {
		items      stringSet
		minimumLen int
		infinite   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty set",
			fields: fields{
				items: newStringSet(),
			},
			want: "",
		},
		{
			name: "{AT, TT}",
			fields: fields{
				items:      newStringSet("AT", "TT"),
				minimumLen: 2,
			},
			want: "T",
		},
		{
			name: "{apple, pineapple, iOS_application}",
			fields: fields{
				items:      newStringSet("apple", "pineapple", "iOS_application"),
				minimumLen: 2,
			},
			want: "appl",
		},
		{
			name: "no common substring {hello, aloha, 123, 456}",
			fields: fields{
				items:      newStringSet("hello", "aloha", "123", "456"),
				minimumLen: 2,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Set{
				items:      tt.fields.items,
				minimumLen: tt.fields.minimumLen,
				infinite:   tt.fields.infinite,
			}
			if got := s.LongestCommon(); got != tt.want {
				t.Errorf("LongestCommon() = %v, want %v", got, tt.want)
			}
		})
	}
}
