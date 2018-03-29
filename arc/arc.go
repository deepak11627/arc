// Package arc implements an Adaptive replacement cache
package arc

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/deepak11627/arc/utils"
)

type ARC struct {
	p     int
	c     int
	t1    *list.List
	b1    *list.List
	t2    *list.List
	b2    *list.List
	mutex sync.RWMutex
	len   int
	cache map[interface{}]*entry
}

// New returns a new Adaptive Replacement Cache (ARC).
func NewARC(c int) *ARC {
	return &ARC{
		p:     0,
		c:     c,
		t1:    list.New(),
		b1:    list.New(),
		t2:    list.New(),
		b2:    list.New(),
		len:   0,
		cache: make(map[interface{}]*entry, c),
	}
}

// Put inserts a new key-value pair into the cache.
// This optimizes future access to this entry (side effect).
func (a *ARC) Put(key, value interface{}) bool {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	ent, ok := a.cache[key]
	if ok != true {
		a.len++

		ent = &entry{
			key:   key,
			value: value,
			ghost: false,
		}

		a.req(ent)
		a.cache[key] = ent
	} else {
		if ent.ghost {
			a.len++
		}
		ent.value = value
		ent.ghost = false
		a.req(ent)
	}
	return ok
}

// Get retrieves a previously via Set inserted entry.
// This optimizes future access to this entry (side effect).
func (a *ARC) Get(key interface{}) (value interface{}, ok bool) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	ent, ok := a.cache[key]
	if ok {
		a.req(ent)
		return ent.value, !ent.ghost
	}
	return nil, false
}

// Len determines the number of currently cached entries.
// This method is side-effect free in the sense that it does not attempt to optimize random cache access.
func (a *ARC) Len() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.len
}

func (a *ARC) req(ent *entry) {
	//fmt.Printf("Entry is %+v\n", ent)
	//fmt.Printf("ARC is %+v\n", a)
	if ent.ll == a.t1 || ent.ll == a.t2 {
		//	fmt.Printf("case 1")
		// Case I
		ent.setMRU(a.t2)
	} else if ent.ll == a.b1 {
		//	fmt.Printf("case 2")
		// Case II
		// Cache Miss in t1 and t2

		// Adaptation
		var d int
		if a.b1.Len() >= a.b2.Len() {
			d = 1
		} else {
			d = a.b2.Len() / a.b1.Len()
		}
		a.p = utils.Min(a.p+d, a.c)

		a.replace(ent)
		ent.setLRU(a.t2)
	} else if ent.ll == a.b2 {
		// Case III
		// Cache Miss in t1 and t2
		//	fmt.Printf("case 3")
		// Adaptation
		var d int
		if a.b2.Len() >= a.b1.Len() {
			d = 1
		} else {
			d = a.b1.Len() / a.b2.Len()
		}
		a.p = utils.Max(a.p-d, 0)

		a.replace(ent)
		ent.setMRU(a.t2)
	} else if ent.ll == nil {
		// Case IV
		//	fmt.Printf("case 4")
		if a.t1.Len()+a.b1.Len() == a.c {
			// Case A
			//		fmt.Printf("case 5")
			if a.t1.Len() < a.c {
				//			fmt.Printf("case 6")
				a.delLRU(a.b1)
				a.replace(ent)
			} else {
				//			fmt.Printf("case 7")
				a.delLRU(a.t1)
			}
		} else if a.t1.Len()+a.b1.Len() < a.c {
			// Case B
			//		fmt.Printf("case 8")
			if a.t1.Len()+a.t2.Len()+a.b1.Len()+a.b2.Len() >= a.c {
				//			fmt.Printf("case 10")
				if a.t1.Len()+a.t2.Len()+a.b1.Len()+a.b2.Len() == 2*a.c {
					//				fmt.Printf("case 11")
					a.delLRU(a.b2)
				}
				a.replace(ent)
			}
		}
		//	fmt.Printf("case 9")
		ent.setLRU(a.t1)
	}
}

func (a *ARC) delLRU(list *list.List) {
	lru := list.Back()
	list.Remove(lru)
	a.len--
	delete(a.cache, lru.Value.(*entry).key)
}

func (a *ARC) replace(ent *entry) {
	if a.t1.Len() > 0 && ((a.t1.Len() > a.p) || (ent.ll == a.b2 && a.t1.Len() == a.p)) {
		lru := a.t1.Back().Value.(*entry)
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setLRU(a.b1)
	} else {
		lru := a.t2.Back().Value.(*entry)
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b2)
	}
}

// Traverse prints the items of a list
func (a *ARC) Traverse() {
	utils.RenderMessageHeading(fmt.Sprintf("Cache Size is %d, %d items are", a.c, a.len))

	// Iterate through list and print its contents.
	for k, v := range a.cache {
		fmt.Printf("Value at %s is %s\n", k, v.value)
	}

	utils.RenderMessageEnd()
	fmt.Println("\nT1 items are")
	for e := a.t1.Front(); e != nil; e = e.Next() {
		fmt.Printf("%s -> %s\n", e.Value.(*entry).key, e.Value.(*entry).value)
	}
	fmt.Println("\nT2 items are")
	for e := a.t2.Front(); e != nil; e = e.Next() {
		fmt.Printf("%s -> %s\n", e.Value.(*entry).key, e.Value.(*entry).value)
	}
	fmt.Println("\nB1 items are")
	for e := a.b1.Front(); e != nil; e = e.Next() {
		fmt.Printf("%s -> %s\n", e.Value.(*entry).key, e.Value.(*entry).value)
	}
	fmt.Println("\nB2 items are")
	for e := a.b2.Front(); e != nil; e = e.Next() {
		fmt.Printf("%s -> %s\n", e.Value.(*entry).key, e.Value.(*entry).value)
	}
}
