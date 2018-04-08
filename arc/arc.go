// Package arc implements an Adaptive replacement cache
package arc

import (
	"fmt"
	"sync"

	"github.com/deepak11627/arc/utils"
)

// ARC struct to implement ARC cache
// In Cache:
// - T1: Pages that have been accessed at least once
// - T2: Pages that have been accessed at least twice
// Ghost:
// - B1: Evicted from T1
// - B2: Evicted from T2
// Adapt:
// - Hit in B1 should increase size of T1, drop entry from T2 to B2
// - Hit in B2 should increase size of T2, drop entry from T1 to B1
type ARC struct {
	p      int
	c      int
	t1     ListService
	t2     ListService
	b1     ListService
	b2     ListService
	mutex  sync.RWMutex
	len    int
	cache  map[interface{}]*entry
	logger Logger
	db     DBService
}

// Option type setting params dynamically
type Option func(*ARC)

// SetLogger function to set logger dynamically
func SetLogger(l Logger) func(*ARC) {
	return func(arc *ARC) {
		arc.logger = l
	}
}

// SetDatabaseListService function to set logger dynamically
func SetDatabaseListService(db DBService) func(*ARC) {
	return func(arc *ARC) {
		arc.db = db
	}
}

// NewARC returns a new Adaptive Replacement Cache (ARC).
func NewARC(c int, t1, t2, b1, b2 ListService, opts ...Option) CacheService {
	arc := &ARC{
		p:     0,
		c:     c,
		t1:    t1,
		t2:    t2,
		b1:    b1,
		b2:    b2,
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

		a.logger.Debug("Adding a new entry item to cache.", "item", fmt.Sprintf("%+v", ent))

		a.req(ent)
		a.cache[key] = ent
	} else {
		a.logger.Debug("Item found in cache, will adjust its position", "item_key", fmt.Sprintf("%s", key))
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
		a.logger.Debug("Reading a value from cache, will adjust its position", "item_keu", fmt.Sprintf("%s", key))
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
		a.logger.Debug("Case 1", "item", fmt.Sprintf("%+v", ent))
		// repetitive entry so should go into MRU
		// Case I
		// x ∈ T1 ∪ T2 (a hit in ARC(c) and DBL(2c)): Move x to the top of T2
		ent.setMRU(a.t2)
	} else if ent.ll == a.b1 {

		a.logger.Debug("Case 2", "item", fmt.Sprintf("%+v", ent))
		// Case II
		// Cache Miss in t1 and t2
		// x ∈ B1 (a miss in ARC(c), a hit in DBL(2c)):
		// Adapt p = min{ c, p + max{ |B2| / |B1|, 1} }. REPLACE(p).
		// Move x to the top of T2 and place it in the cache.
		// Adaptation
		var d int
		if a.b1.Len() >= a.b2.Len() {
			d = 1
		} else {
			d = a.b2.Len() / a.b1.Len()
		}
		a.p = utils.Min(a.p+d, a.c)

		a.replace(ent)
		ent.setMRU(a.t2)
	} else if ent.ll == a.b2 {
		a.logger.Debug("Case 3", "item", fmt.Sprintf("%+v", ent))
		// Case III
		// Cache Miss in t1 and t2
		// x ∈ B2 (a miss in ARC(c), a hit in DBL(2c)):
		// Adapt p = max{ 0, p – max{ |B1| / |B2|, 1} } . REPLACE(p).
		// Move x to the top of T2 and place it in the cache.
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
		a.logger.Debug("Case 4", "item", fmt.Sprintf("%+v", ent))
		// Case IV
		// x ∈ L1 ∪ L2 (a miss in DBL(2c) and ARC(c)):
		// case (i) |L1| = c:
		//   If |T1| < c then delete the LRU page of B1 and REPLACE(p).
		//   else delete LRU page of T1 and remove it from the cache.
		// case (ii) |L1| < c and |L1| + |L2|≥ c:
		//   if |L1| + |L2|= 2c then delete the LRU page of B2.
		//   REPLACE(p) .
		// Put x at the top of T1 and place it in the cache.
		if a.t1.Len()+a.b1.Len() == a.c {
			// Case A
			if a.t1.Len() < a.c {
				a.delLRU(a.b1)
				if a.db != nil {
					a.db.Remove("B1")
				}
				a.replace(ent)
			} else {
				a.delLRU(a.t1)
			}
		} else if a.t1.Len()+a.b1.Len() < a.c {
			// Case B
			if a.t1.Len()+a.t2.Len()+a.b1.Len()+a.b2.Len() >= a.c {
				if a.t1.Len()+a.t2.Len()+a.b1.Len()+a.b2.Len() == 2*a.c {
					a.delLRU(a.b2)
					if a.db != nil {
						a.db.Remove("B2")
					}
				}
				a.replace(ent)
			}
		}
		ent.setMRU(a.t1)
	}
	a.logger.Debug("Adaptation value was", "p", a.p)
}

func (a *ARC) delLRU(l ListService) {
	lru := l.Back()
	a.logger.Debug("Removing item from list", "item", fmt.Sprintf("%+v", lru))
	l.Remove(lru)
	a.len--
	delete(a.cache, lru.Value.(*entry).key)
}

func (a *ARC) replace(ent *entry) {
	// if (|T1| ≥ 1) and ((x ∈ B2 and |T1| = p) or (|T1| > p))
	//   then move the LRU page of T1 to the top of B1 and remove it from the cache.
	// else move the LRU page in T2 to the top of B2 and remove it from the cache.
	if a.t1.Len() > 0 && ((a.t1.Len() > a.p) || (ent.ll == a.b2 && a.t1.Len() == a.p)) {
		lru := a.t1.Back().Value.(*entry)
		a.logger.Debug("Moving item from T1 to B1", "item", fmt.Sprintf("%+v", lru))
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b1)
		// Archieve  Evicted items to database
		if a.db != nil {
			a.db.PushFront("B1", lru.key, lru.value)
		}

	} else {
		lru := a.t2.Back().Value.(*entry)
		a.logger.Debug("Moving item from T2 to B2", "item", fmt.Sprintf("%+v", lru))
		lru.value = nil
		lru.ghost = true
		a.len--
		lru.setMRU(a.b2)
		// Archieve  Evicted items to database
		if a.db != nil {
			a.db.PushFront("B2", lru.key, lru.value)
		}
	}
}

// Traverse prints the items of a list
func (a *ARC) Traverse() {
	utils.RenderMessageHeading("Items are cached.")

	// Iterate through list and print its contents.
	for k, v := range a.cache {
		fmt.Printf("Value at %s is %s (isGhost - %v) \n", k, v.value, v.ghost)

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
