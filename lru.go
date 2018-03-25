package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//CacheSize the number of maximum values to cache.
var CacheSize int

// LRUCache is Least recently Used cache
type LRUCache struct {
	*list.List
	Count int
	Heap  map[string]string
}

// Option are func to dynamically set LRUCache values
type Option func(*LRUCache)

// NewLRUCache create a new LRUCache cache
func NewLRUCache(opts ...Option) *LRUCache {

	lru := &LRUCache{
		Count: 0,
		List:  list.New(),
		Heap:  make(map[string]string, 0),
	}
	for _, opt := range opts {
		opt(lru)
	}

	return lru
}

func (c *LRUCache) Traverse() {
	var items []interface{}
	// Iterate through list and print its contents.
	for e := c.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}

	RenderMessageHeading("Current cached items are")
	// Iterate through list and print its contents.
	for _, item := range items {
		switch item.(type) {
		case int:
			fmt.Printf("  | %d |  ", item.(int))
		case float32:
			fmt.Printf("  | %f |  ", item.(float32))
		default:
			fmt.Print(fmt.Sprintf("  | %s  |  ", item))
		}
	}

	RenderMessageEnd()
}

func main() {

	Message("Please enter maximum number of keys which caching system should store. ")
	for CacheSize == 0 {
		SetCacheSize()
	}
	fmt.Printf("Thank you %s", CacheSize)
}

//SetCacheSize takes input from user and sets value for CacheSize
func SetCacheSize() {
	reader := bufio.NewReader(os.Stdin)
	val, err := reader.ReadString('\n')
	if err != nil {
		Message("Problem reading the entered value.")
		os.Exit(1)
	}
	val = strings.Replace(val, "\n", "", -1)
	ok, err := strconv.Atoi(val)
	if err != nil {
		Message("Numbers are the only accepted value.")
	} else {
		CacheSize = ok
	}
}
func showOptions() {
	RenderMessageHeading("")
}
