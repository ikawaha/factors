package factors

import (
	"sort"
	"strings"
)

const theta = "θ"

type stringSet map[string]struct{}

func newStringSet(items ...string) stringSet {
	ret := stringSet{}
	for _, v := range items {
		ret[v] = struct{}{}
	}
	return ret
}

// Set represents a set of necessary factors.
type Set struct {
	items      stringSet
	minimumLen int
	infinite   bool
}

// NewSet creates a set initialized given items.
func NewSet(items ...string) Set {
	var ret Set
	for _, v := range items {
		ret.Add(v)
	}
	return ret
}

// Items returns (sorted) items of the set.
func (s Set) Items() []string {
	if s.infinite || len(s.items) == 0 {
		return nil
	}
	ret := make([]string, 0, len(s.items))
	for k := range s.items {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

// Add adds a item to the set.
func (s *Set) Add(item string) {
	if s.infinite {
		return
	}
	if s.items == nil {
		s.items = stringSet{}
	}
	s.items[item] = struct{}{}
	if s.minimumLen == 0 || len(item) < s.minimumLen {
		s.minimumLen = len(item)
	}
}

// longest common substring
func longestCommonSubstring(x, y string) string {
	if len(x) == 0 || len(y) == 0 {
		return ""
	}
	var max int
	matrix := make([][]int, len(x))
	for i := range matrix {
		matrix[i] = make([]int, len(y))
	}
	var p int
	for i := 0; i < len(x); i++ {
		for j := 0; j < len(y); j++ {
			if x[i] == y[j] {
				if i == 0 || j == 0 {
					matrix[i][j] = 1
					if max < matrix[i][j] {
						max = matrix[i][j]
						p = i
					}
					continue
				}
				matrix[i][j] = matrix[i-1][j-1] + 1
				if max < matrix[i][j] {
					max = matrix[i][j]
					p = i
				}
			}
		}
	}
	return x[(p+1)-max : p+1]
}

// LongestCommon returns the longest common substring of items in the set.
func (s Set) LongestCommon() string {
	if len(s.items) == 0 || s.infinite {
		return ""
	}
	items := s.Items()
	ret := items[0]
	for i := 1; i < len(items); i++ {
		ret = longestCommonSubstring(ret, items[i])
		if ret == "" {
			return ""
		}
	}
	return ret
}

// Clear clears the set.
func (s *Set) Clear() {
	s.infinite = false
	s.minimumLen = 0
	s.items = nil
}

// SetInfinite sets this set infinite.
func (s *Set) SetInfinite() {
	s.infinite = true
	s.minimumLen = 0
	s.items = nil
}

// Infinite returns true if this set is infinite.
func (s Set) Infinite() bool {
	return s.infinite
}

// DropRedundantPrefix drops items which has prefix of other item in this set.
func (s *Set) DropRedundantPrefix() {
	if s.infinite || len(s.items) == 0 {
		return
	}
	ps := make([]string, 0, len(s.items))
	for k := range s.items {
		ps = append(ps, k)
	}
	sort.Strings(ps)
	items := make([]string, 0, len(ps))
	items = append(items, ps[0])
	for i := 1; i < len(ps); i++ {
		if strings.HasPrefix(ps[i], items[len(items)-1]) {
			continue
		}
		items = append(items, ps[i])
	}
	s.items = newStringSet(items...)
}

func sortByRevertedString(s []string) {
	sort.Slice(s, func(i, j int) bool {
		a := s[i]
		b := s[j]
		for len(a) > 0 && len(b) > 0 {
			if a[len(a)-1] == b[len(b)-1] {
				a = a[:len(a)-1]
				b = b[:len(b)-1]
				continue
			}
			return a[len(a)-1] < b[len(b)-1]

		}
		return len(a) < len(b)
	})
}

// DropRedundantSuffix drops items which has suffix of other item in this set.
func (s *Set) DropRedundantSuffix() {
	if s.infinite || len(s.items) == 0 {
		return
	}
	ss := make([]string, 0, len(s.items))
	for k := range s.items {
		ss = append(ss, k)
	}
	sortByRevertedString(ss)
	items := make([]string, 0, len(ss))
	items = append(items, ss[0])
	for i := 1; i < len(ss); i++ {
		if strings.HasSuffix(ss[i], items[len(items)-1]) {
			continue
		}
		items = append(items, ss[i])
	}
	s.items = newStringSet(items...)
}

// DropRedundantFragment drops items which contains of other item in this set.
func (s *Set) DropRedundantFragment() {
	fs := s.Items()
loop:
	for i := 0; i < len(fs); i++ {
		for j := 0; j < len(fs); j++ {
			if i == j || fs[i] == "" {
				continue
			}
			if strings.Contains(fs[j], fs[i]) {
				fs[j] = ""
				continue loop
			}
		}
	}
	items := make([]string, 0, len(fs))
	for _, v := range fs {
		if v != "" {
			items = append(items, v)
		}
	}
	s.items = newStringSet(items...)
}

// Len returns a size of this set.
func (s Set) Len() int {
	if s.infinite {
		return -1
	}
	return len(s.items)
}

// String returns string represents of this set.
func (s Set) String() string {
	if s.infinite {
		return theta
	}
	return "{" + strings.Join(s.Items(), ", ") + "}"
}

// UnionSet returns a union set of x and y.
func UnionSet(x, y Set) Set {
	var ret Set
	ret.infinite = x.infinite || y.infinite
	if ret.infinite {
		return ret
	}
	ret.items = make(stringSet, len(x.items)+len(y.items))
	for k := range x.items {
		ret.Add(k)
	}
	for k := range y.items {
		ret.Add(k)
	}
	ret.minimumLen = x.minimumLen
	if x.minimumLen > y.minimumLen {
		ret.minimumLen = y.minimumLen
	}
	return ret
}

// CrossSet returns a cross set of x and y.
// e.g. CrossSet {a, bc} and {xy, z} -> {axy, az, bcxy, bcz}
func CrossSet(x, y Set) Set {
	var ret Set
	ret.infinite = x.infinite || y.infinite
	if ret.infinite {
		return ret
	}
	ret.items = make(stringSet, len(x.items)*len(y.items))
	for k0 := range x.items {
		for k1 := range y.items {
			ret.items[k0+k1] = struct{}{}
		}
	}
	ret.minimumLen = x.minimumLen + y.minimumLen
	return ret
}

// BestSet chooses the best set from the given sets.
func BestSet(arg Set, args ...Set) Set {
	best := arg
	for _, v := range args {
		if best.minimumLen > v.minimumLen {
			continue
		}
		if best.minimumLen == v.minimumLen {
			if len(best.items) < len(v.items) {
				continue
			}
		}
		best = v
	}
	return best
}
