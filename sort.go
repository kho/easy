package easy

import (
	"sort"
)

type WithIndex struct {
	Sorter   sort.Interface
	NewToOld []int
}

func NewWithIndex(sorter sort.Interface) *WithIndex {
	n := sorter.Len()
	newToOld := make([]int, n)
	for i := range newToOld {
		newToOld[i] = i
	}
	return &WithIndex{sorter, newToOld}
}

func (s *WithIndex) Len() int           { return s.Sorter.Len() }
func (s *WithIndex) Less(i, j int) bool { return s.Sorter.Less(i, j) }
func (s *WithIndex) Swap(i, j int) {
	s.NewToOld[i], s.NewToOld[j] = s.NewToOld[j], s.NewToOld[i]
	s.Sorter.Swap(i, j)
}

func SortWithIndex(sorter sort.Interface) []int {
	wi := NewWithIndex(sorter)
	sort.Sort(wi)
	return wi.NewToOld
}
