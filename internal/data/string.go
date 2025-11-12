package data

import (
	"fmt"
)


// StringValue represents a string value.
type StringValue struct {
	Val string
}

// NewStringValue creates a new StringValue.
func NewStringValue(val string) *StringValue {
	return &StringValue {Val: val}
}

// String returns the string representation of the StringValue.
func (s *StringValue) String() string {
	return fmt.Sprintf("StringValue(%s)", s.Val)
}	



