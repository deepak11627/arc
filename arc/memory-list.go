package arc

import (
	"container/list"
)

type MemoryList struct {
	list.List
}

func NewMemoryList() ListService {
	return &MemoryList{}
}
