package factor

import (
	"reflect"
	"sort"
	"testing"
)

func TestAdd(t *testing.T) {
	var (
		f     Factor
		items = []string{"aaaaa", "bbb", "cc"}
	)
	for _, v := range items {
		f.Add(v)
		f.Add(v)
	}
	if f.MinFactorLen != 2 {
		t.Errorf("min factor len: got %+v, expected %+v, %+v", f.MinFactorLen, 2, f)
	}

	if len(f.data) != len(items) {
		t.Errorf("data size: got %+v, expected %+v", len(f.data), len(items))
	}
	for _, v := range items {
		if _, ok := f.data[v]; !ok {
			t.Errorf("got %+v, expected %+v", len(f.data), len(items))
			break
		}
	}
}

func TestSize(t *testing.T) {
	var (
		f     Factor
		items = []string{"a", "b", "c"}
	)
	if f.Size() != 0 {
		t.Errorf("got %+v, expected %+v", len(f.data), 0)
	}
	for _, v := range items {
		f.Add(v)
	}
	if f.Size() != len(items) {
		t.Errorf("got %+v, expected %+v", len(f.data), len(items))
	}
}

func TestItems(t *testing.T) {
	var (
		f     Factor
		items = []string{"b", "a", "c"}
	)
	for _, v := range items {
		f.Add(v)
	}

	list := f.Items()

	sort.Strings(items)
	for i := 0; i < len(items); i++ {
		if items[i] != list[i] {
			t.Errorf("got %+v, expected %+v", list, items)
			break
		}
	}
}

func TestUnion(t *testing.T) {
	var (
		x, y Factor
		s    = []string{"aa1", "bb1", "cc1", "a2", "b2", "cccc2"}
		s1   = s[:3]
		s2   = s[3:]
	)

	for _, v := range s1 {
		x.Add(v)
	}
	if x.MinFactorLen != 3 {
		t.Errorf("min factor len: got %+v, expected %+v\nx=%+v", x.MinFactorLen, 3, x)
	}

	for _, v := range s2 {
		y.Add(v)
	}
	if y.MinFactorLen != 2 {
		t.Errorf("min factor len: got %+v, expected %+v\ny=%+v", y.MinFactorLen, 2, y)
	}

	x.Union(y)
	sort.Strings(s)
	list := x.Items()
	if !reflect.DeepEqual(s, list) {
		t.Errorf("got %+v, expected %+v", list, s)
	}
	if x.MinFactorLen != 2 {
		t.Errorf("min factor len: got %+v, expected %+v\n%+v", x.MinFactorLen, 2, x)
	}

	y.Undef = true
	x.Union(y)
	if !x.Undef {
		t.Errorf("undef: expected true, %#v", x)
	}
	if x.data != nil {
		t.Errorf("data: expected nil, %#v", x)
	}
}

func TestCrossSet(t *testing.T) {
	var (
		x, y  Factor
		s     = []string{"aa1", "bb1", "cc1", "a2", "b2", "cccc2"}
		s1    = s[:3]
		s2    = s[3:]
		cross = []string{
			"aa1a2", "aa1b2", "aa1cccc2",
			"bb1a2", "bb1b2", "bb1cccc2",
			"cc1a2", "cc1b2", "cc1cccc2",
		}
	)

	for _, v := range s1 {
		x.Add(v)
	}
	if x.MinFactorLen != 3 {
		t.Errorf("min factor len: got %+v, expected %+v\nx=%+v", x.MinFactorLen, 3, x)
	}

	for _, v := range s2 {
		y.Add(v)
	}
	if y.MinFactorLen != 2 {
		t.Errorf("min factor len: got %+v, expected %+v\ny=%+v", y.MinFactorLen, 2, y)
	}

	xy := CrossSet(x, y)
	if xy.Undef {
		t.Errorf("undef: expected false, %#v", xy)
	}
	if xy.MinFactorLen != x.MinFactorLen+y.MinFactorLen {
		t.Errorf("min factor len: got %+v, expected %+v\nxy=%+v", xy.MinFactorLen, 2, xy)
	}
	if !reflect.DeepEqual(xy.Items(), cross) {
		t.Errorf("got %+v, expected %v", xy, cross)
	}

	y.Undef = true
	xy = CrossSet(x, y)
	if !xy.Undef {
		t.Errorf("undef: expected true, %#v", xy)
	}
	if xy.data != nil {
		t.Errorf("data: expected nil, %#v", xy)
	}
}

func TestClear(t *testing.T) {
	var f Factor
	f.Add("a")
	f.Undef = true
	f.Clear()
	if f.Undef || len(f.data) != 0 || f.MinFactorLen != 0 {
		t.Errorf("does not cleared undef, %+v", f)
	}
}

func TestClone(t *testing.T) {
	var x Factor
	data := []string{"aa", "bb", "cc", "a", "b", "cccc"}
	for _, v := range data {
		x.Add(v)
	}
	y := x.Clone()
	if !reflect.DeepEqual(x, y) {
		t.Errorf("got %#v, expected %#v", y, x)
	}
	x.Undef = true
	x.Add("hello")
	if reflect.DeepEqual(x, y) {
		t.Errorf("does not cloned, x=%#v, y=%#v", y, x)
	}
}

func TestBestSet(t *testing.T) {
	e := struct{}{}
	a := Factor{Undef: true}
	b := Factor{MinFactorLen: 1, data: StrageType{"a": e, "b": e, "ccc": e}}
	c := Factor{MinFactorLen: 2, data: StrageType{"aa": e, "bb": e, "cc": e}}
	d := Factor{MinFactorLen: 2, data: StrageType{"aa": e, "bbb": e}}
	best := BestSet(&a, &b, &c, &d)
	if best != &d {
		t.Errorf("got %+v, expected %+v", best, d)
	}
}
