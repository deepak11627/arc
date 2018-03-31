// Package arc implements an Adaptive replacement cache
package arc

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/deepak11627/arc/utils"
)

// ARC struct to implement ARC cache
type ARC struct {
	p      int
	c      int
	t1     *list.List
	t2     *list.List
	b1     *list.List
	b2     *list.List
	mutex  sync.RWMutex
	len    int
	cache  map[interface{}]*entry
	logger Logger
}

// Option type setting params dynamically
type Option func(*ARC)

// SetLogger function to set logger dynamically
func SetLogger(l Logger) func(*ARC) {
	return func(arc *ARC) {
		arc.logger = l
	}
}

// NewARC returns a new Adaptive Replacement Cache (ARC).
func NewARC(c int, opts ...Option) CacheService {
	arc := &ARC{
		p:     0,
		c:     c,
		t1:    list.New(),
		t2:    list.New(),
		b1:    list.New(),
		b2:    list.New(),
		len:   0,
		cache: make(map[interface{}]*entry, c),
	}
	for _, o := range opts {
		o(arc)
	}
	return arc
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
	if ent.ll == a.t1 || ent.ll == a.t2 {

		// repetitive entry so should go into MRU
		// Case I
		ent.setMRU(a.t2)
	} else if ent.ll == a.b1 {

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

		//ent.setLRU(a.t2)
		ent.setMRU(a.t2)
	} else if ent.ll == a.b2 {
		// Case III
		// Cache Miss in t1 and t2

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
		if a.t1.Len()+a.b1.Len() == a.c {
			// Case A
			if a.t1.Len() < a.c {
				a.delLRU(a.b1)
				a.replace(ent)
			} else {
				a.delLRU(a.t1)
			}
		} else if a.t1.Len()+a.b1.Len() < a.c {
			// Case B
			if a.t1.Len()+a.t2.Len()+a.b1.Len()+a.b2.Len() >= a.c {
				if a.t1.Len()+a.t2.Len()+a.b1.Len()+a.b2.Len() == 2*a.c {
					a.delLRU(a.b2)
				}
				a.replace(ent)
			}
		}
		ent.setMRU(a.t1)
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
		//	lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b1)
	} else {
		lru := a.t2.Back().Value.(*entry)
		//	lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b2)
	}
}

// Traverse prints the items of a list
func (a *ARC) Traverse() {
	utils.RenderMessageHeading(fmt.Sprintf("Cache Size is %d, %d # items are cached.\n", a.c, a.Len()))

	// Iterate through list and print its contents.
	for k, v := range a.cache {
		if !v.ghost {
			fmt.Printf("Value at %s is %s\n", k, v.value)
		}
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
