package easy

import (
	"sort"
	"testing"
	"testing/quick"
)

func TestSortWithIndex(t *testing.T) {
	// Empty input.
	if n2o := SortWithIndex(sort.IntSlice([]int{})); len(n2o) != 0 {
		t.Error("sorting empty slice yields ", n2o)
	}
	// Non-empty input.
	if err := quick.Check(checkSortWithIndexNonEmpty, nil); err != nil {
		t.Error(err)
	}
}

func checkSortWithIndexNonEmpty(s []int) bool {
	s = append(s, 0)
	t := make([]int, len(s))
	copy(t, s)
	n2o := SortWithIndex(sort.IntSlice(s))
	for i, v := range s[1:] {
		if v < s[i] {
			return false
		}
	}
	for n, o := range n2o {
		if s[n] != t[o] {
			return false
		}
	}
	return true
}

func TestReverse(t *testing.T) {
	if err := quick.Check(checkReverse, nil); err != nil {
		t.Error(err)
	}
}

func checkReverse(s []int) bool {
	sort.Sort(Reverse{sort.IntSlice(s)})
	if len(s) > 0 {
		for i, v := range s[1:] {
			if s[i] < v {
				return false
			}
		}
	}
	return true
}
