package arc

import (
	"container/list"
	"fmt"
)

type MemoryList struct {
	*list.List
}

func NewMemoryList() ListService {
	return MemoryList{list.New()}
}

func (ml *MemoryList) Value() {
	fmt.Printf("%d", ml.Len())
}
