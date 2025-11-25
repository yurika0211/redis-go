package data 

import (
	"fmt"
)

type SortedSet struct {
	Val map[string]int
}

/**
 * NewSortedSet creates a new SortedSet.
 */
func NewSortedSet(val map[string]int) *SortedSet {
	return &SortedSet{Val: val}
}

/**
 * String returns the string representation of the SortedSet.
 */
func (s *SortedSet) String() string {
	return fmt.Sprintf("SortedSet(%v)", s.Val)
}

func (s *SortedSet) Add(score string , member string) {
	// For simplicity, we convert score to an integer.
	// In a real implementation, you'd handle errors and different types.
	var intScore int
	fmt.Sscanf(score, "%d", &intScore)
	s.Val[member] = intScore
}