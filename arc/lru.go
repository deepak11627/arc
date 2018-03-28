package arc

import (
	"container/list"
	"fmt"

	"github.com/deepak11627/arc/utils"
)

// LRU is Least recently Used cache
type LRU struct {
	*list.List
	Size int
	Heap map[string]string
}

// NewLRU create a new LRU cache
func NewLRU(size int) ListService {

	lru := &LRU{
		Size: size,
		List: list.New(),
		//Heap:  make(map[string]string, 0),
	}
	// for _, opt := range opts {
	// 	opt(lru)
	// }

	return lru
}

// Traverse prints the items of a list
func (c *LRU) Traverse() {
	utils.RenderMessageHeading("Current cached items are")
	// Iterate through list and print its contents.
	for e := c.Front(); e != nil; e = e.Next() {
		fmt.Print(fmt.Sprintf(" %s", e.Value))

	}

	utils.RenderMessageEnd()
}
