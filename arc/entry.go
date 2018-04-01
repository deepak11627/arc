package arc

import (
	"container/list"
)

type entry struct {
	key   interface{}
	value interface{}
	ll    ListService
	el    *list.Element
	ghost bool
}

func (e *entry) setLRU(l interface{}) {
	e.detach()
	e.ll = l.(*list.List)
	e.el = e.ll.PushBack(e)
}

func (e *entry) setMRU(l interface{}) {
	e.detach()
	e.ll = l.(*list.List)
	e.el = e.ll.PushFront(e)
}

func (e *entry) detach() {
	if e.ll != nil {
		e.ll.Remove(e.el)
	}
}
