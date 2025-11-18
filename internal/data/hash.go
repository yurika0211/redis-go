package data

import (
	"fmt"
)

// hash represents a hash table
type Hash struct {
	fields map[string]string
}

/**
 * NewHash creates a new Hash value.
 */
func NewHash() *Hash {
	return &Hash{fields: make(map[string]string)}
}

/**
 * NewElement creates a new element in the Hash.
 */
func NewElement(key string, value string, h *Hash) map[string]string {
	h.fields[key] = value
	return h.fields
}

/**
 * String returns the string representation of the Hash.
 */
func (h *Hash) String() string {
	var result string
	for k, v := range h.fields {
		result += fmt.Sprintf("Hash(%s:%s)", k, v)
	}
	return result
}

// SetField sets a field in the hash.
func (h *Hash) SetField(field string, val string) {
	h.fields[field] = val
}

// GetField retrieves a field value and a boolean indicating existence.
func (h *Hash) GetField(field string) (string, bool) {
	v, ok := h.fields[field]
	return v, ok
}
