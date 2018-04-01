package arc

import (
	"container/list"
)

// CacheService is interface for ARC
type CacheService interface {
	Get(key interface{}) (value interface{}, ok bool)
	Put(key, value interface{}) bool
	Traverse()
	Len() int
}

// Logger is used for logging
type Logger interface {
	// Debug logging: an informative message that can aid in debugging.
	Debug(msg string, keyvals ...interface{})
	// Info logging: a high-level event occurred such as a business transaction completed.
	Info(msg string, keyvals ...interface{})
	// Warn logging: the component is not operating optimally but is not in immediate danger of disrupting service.
	Warn(msg string, keyvals ...interface{})
	// Error logging: The component is not operating correctly and needs urgent assistance.
	Error(msg string, keyvals ...interface{})
}

// DBService is a generic User operations interface
type DBService interface {
	Get(ListID int, key interface{}) (interface{}, error)
	Add(ListID int, key interface{}, value interface{}) error
}

// EntryService is to allow ghost entries to be stored into the database
type EntryService interface {
	setLRU(l interface{})
	setMRU(l interface{})
}

// ElementService is to maintain list items
type ElementService interface {
	Next() ElementService
	Prev() ElementService
}

// ListService for generic lists
type ListService interface {
	Len() int
	Back() *list.Element
	Front() *list.Element
	Remove(e *list.Element) interface{}
	PushBack(i interface{}) *list.Element
	PushFront(i interface{}) *list.Element
}
