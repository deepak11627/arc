package arc

import "container/list"

type ListService interface {
	PushFront(v interface{}) *list.Element
	Remove(e *list.Element) interface{}
	Traverse()
	Len() int
}

// Option are func to dynamically set LRU values
type Option func(*LRU)
