package data

import (
	"fmt"
	"strings"
)

type Set struct {
	members map[string]struct{}
}

/**
 * NewSet creates a new Set.
 */
func NewSet(val string) *Set {
	return &Set{members: map[string]struct{}{val: {}}}
}

/**
 * String returns the string representation of the Set.
 */
func (s *Set) String() string {
	var members []string
	for k := range s.members {
		members = append(members, k)
	}
	return fmt.Sprintf("Set(%s)", strings.Join(members, ", "))
}

func (s *Set) AddMember(member string) {
	s.members[member] = struct{}{}
}

func (s *Set) Members() []string {
	var members []string
	for k := range s.members {
		members = append(members, k)
	}
	return members
}
