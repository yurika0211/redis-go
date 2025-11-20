package data

import (
	"container/list"
)

// List represents a list value.
type List struct {
	List *list.List
}

// NewList creates a new List and appends the provided values in order.
func NewList(vals []string) *List {
	l := list.New()
	for _, v := range vals {
		l.PushBack(v)
	}
	return &List{List: l}
}

/**
 * Push a new value to the list.
 */
func (this *List) Push(val string) *List {
	this.List.PushBack(val)
	return this
}

// Values returns all elements in the list as a slice of strings (from front to back).
func (this *List) Values() []string {
	out := make([]string, 0, this.List.Len())
	for e := this.List.Front(); e != nil; e = e.Next() {
		if s, ok := e.Value.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
