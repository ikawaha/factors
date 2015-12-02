package factor

import (
	"sort"
	"strings"
)

type StrageType map[string]struct{}

type Factor struct {
	Undef        bool
	MinFactorLen int
	data         StrageType
}

func NewSet() Factor {
	return Factor{data: make(StrageType)}
}

func (f Factor) Items() []string {
	if f.Undef {
		return nil
	}
	ret := make([]string, 0, len(f.data))
	for k, _ := range f.data {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

func (f Factor) Size() int {
	return len(f.data)
}

func (f *Factor) Add(str string) {
	if f.data == nil {
		f.data = StrageType{}
	}
	if _, ok := f.data[str]; !ok {
		f.data[str] = struct{}{}
		if f.MinFactorLen == 0 || len(str) < f.MinFactorLen {
			f.MinFactorLen = len(str)
		}
	}
}

func (f *Factor) Union(ff Factor) {
	if f.Undef || ff.Undef {
		f.Undef, f.data = true, nil
		return
	}
	for k := range ff.data {
		f.Add(k)
	}
	if ff.MinFactorLen < f.MinFactorLen {
		f.MinFactorLen = ff.MinFactorLen
	}
}

func (f *Factor) Clear() {
	f.Undef = false
	f.MinFactorLen = 0
	f.data = nil
}

func (f Factor) Clone() Factor {
	var dst Factor
	dst.Undef = f.Undef
	dst.MinFactorLen = f.MinFactorLen
	dst.data = StrageType{}
	for k, v := range f.data {
		dst.data[k] = v
	}
	return dst
}

func (f Factor) String() string {
	if f.Undef {
		return "θ"
	}
	tmp := make([]string, 0, len(f.data))
	for k := range f.data {
		tmp = append(tmp, k)
	}
	sort.Strings(tmp)
	return "{" + strings.Join(tmp, ",") + "}"
}

func CrossSet(x, y Factor) (xy Factor) {
	xy.Undef = x.Undef || y.Undef
	if xy.Undef {
		return
	}
	for a := range x.data {
		for b := range y.data {
			xy.Add(a + b)
		}
	}
	xy.MinFactorLen = x.MinFactorLen + y.MinFactorLen
	return
}

func BestSet(sets ...*Factor) (best *Factor) {
	//fmt.Printf("best:%+v\n", sets)
	//defer func() { fmt.Printf("-->%+v\n", best) }()

	if len(sets) == 0 {
		return nil
	}
	best = sets[0]
	for i := 1; i < len(sets); i++ {
		if sets[i].Undef {
			continue
		}
		if best.Undef || (best.MinFactorLen <= sets[i].MinFactorLen &&
			best.Size() >= sets[i].Size()) {
			best = sets[i]
		}
	}
	return best
}
