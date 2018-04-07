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
	switch l.(type) { // type assertion to check if its list or db
	case *list.List:
		e.ll = l.(*list.List)
		// case *models.GhostList:
		// 	e.ll = l.(*models.GhostList)
	}
	e.el = e.ll.PushBack(e)
}

func (e *entry) setMRU(l interface{}) {
	e.detach()
	switch l.(type) { // type assertion to check if its list or db
	case *list.List:
		e.ll = l.(*list.List)
		// case *models.GhostList:
		// 	e.ll = l.(*models.GhostList)
	}

	e.el = e.ll.PushFront(e)
}

func (e *entry) detach() {
	if e.ll != nil {
		e.ll.Remove(e.el)
	}
}
